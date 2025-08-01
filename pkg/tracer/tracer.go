package tracer

import (
	"context"
	"errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"time"

	"github.com/SkyAPM/go2sky"
	sgin "github.com/SkyAPM/go2sky-plugins/gin/v3"
	sreporter "github.com/SkyAPM/go2sky/reporter"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"github.com/opentracing-contrib/go-gin/ginhttp"
	"github.com/opentracing/opentracing-go"
	zipkinot "github.com/openzipkin-contrib/zipkin-go-opentracing"
	"github.com/openzipkin/zipkin-go"
	zreporter "github.com/openzipkin/zipkin-go/reporter"
	zipkinreporter "github.com/openzipkin/zipkin-go/reporter/http"
	jaegerconfig "github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Tracer interface {
	ServerInterceptor() grpc.UnaryServerInterceptor
	ClientInterceptor() grpc.UnaryClientInterceptor
	GinMiddleware(r *gin.Engine) gin.HandlerFunc
	Close() error
}

type ZipkinConfig struct {
	reporter zreporter.Reporter
	tracer   opentracing.Tracer
}

func NewZipkin(zipkinAddr, service, host string, random float64) (Tracer, error) {
	if random > 1 || random < 0 {
		return nil, errors.New("random number must be between 0 and 1")
	}
	// 初始化 zipkin tracer
	report := zipkinreporter.NewReporter(zipkinAddr)
	endpoint, err := zipkin.NewEndpoint(service, host)
	if err != nil {
		return nil, err
	}

	// 设置取样概率
	sampler, err := zipkin.NewBoundarySampler(random, time.Now().UnixNano())
	if err != nil {
		return nil, err
	}

	tracer, err := zipkin.NewTracer(
		report,
		zipkin.WithLocalEndpoint(endpoint),
		zipkin.WithSampler(sampler),
	)
	if err != nil {
		return nil, err
	}

	t := zipkinot.Wrap(tracer)

	return &ZipkinConfig{
		reporter: report,
		tracer:   t,
	}, nil
}

func (z *ZipkinConfig) ServerInterceptor() grpc.UnaryServerInterceptor {
	opentracing.SetGlobalTracer(z.tracer)
	return grpc_opentracing.UnaryServerInterceptor()
}

func (z *ZipkinConfig) ClientInterceptor() grpc.UnaryClientInterceptor {
	opentracing.SetGlobalTracer(z.tracer)
	return grpc_opentracing.UnaryClientInterceptor()
}

func (z *ZipkinConfig) GinMiddleware(_ *gin.Engine) gin.HandlerFunc {
	return ginhttp.Middleware(z.tracer)
}

func (z *ZipkinConfig) Close() error {
	return z.reporter.Close()
}

type JaegerConfig struct {
	closer io.Closer
	tracer opentracing.Tracer
}

func NewJaeger(jaegerAddr, service string, random float64) (Tracer, error) {
	var style string
	if random > 1 || random < 0 {
		return nil, errors.New("random number must be between 0 and 1")
	}
	if random == 1 || random == 0 {
		style = "const"
	} else {
		style = "probabilistic"
	}
	cfg := jaegerconfig.Configuration{
		ServiceName: service,
		Sampler: &jaegerconfig.SamplerConfig{
			Type:  style,
			Param: random,
		},
		Reporter: &jaegerconfig.ReporterConfig{
			CollectorEndpoint: jaegerAddr,
		},
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		return nil, err
	}

	return &JaegerConfig{
		closer: closer,
		tracer: tracer,
	}, nil
}

func (j *JaegerConfig) ServerInterceptor() grpc.UnaryServerInterceptor {
	opentracing.SetGlobalTracer(j.tracer)
	return grpc_opentracing.UnaryServerInterceptor()
}

func (j *JaegerConfig) ClientInterceptor() grpc.UnaryClientInterceptor {
	opentracing.SetGlobalTracer(j.tracer)
	return grpc_opentracing.UnaryClientInterceptor()
}

func (j *JaegerConfig) GinMiddleware(_ *gin.Engine) gin.HandlerFunc {
	return ginhttp.Middleware(j.tracer)
}

func (j *JaegerConfig) Close() error {
	return j.closer.Close()
}

type SkyWalkingConfig struct {
	tracer   *go2sky.Tracer
	reporter go2sky.Reporter
}

func NewSkyWalking(SkyWalkingAddr, service, instance string, random float64) (Tracer, error) {
	if random > 1 || random < 0 {
		return nil, errors.New("random number must be between 0 and 1")
	}
	rep, err := sreporter.NewGRPCReporter(SkyWalkingAddr)
	if err != nil {
		return nil, err
	}

	tracer, err := go2sky.NewTracer(service, go2sky.WithInstance(instance), go2sky.WithReporter(rep), go2sky.WithSampler(random))
	if err != nil {
		return nil, err
	}

	return &SkyWalkingConfig{
		tracer:   tracer,
		reporter: rep,
	}, nil
}

func (s *SkyWalkingConfig) ServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		extractor := func(headerKey string) (string, error) {
			if v := md.Get(headerKey); len(v) > 0 { // 使用动态传入的 headerKey
				return v[0], nil
			}
			return "", nil
		}

		span, ctx, err := s.tracer.CreateEntrySpan(ctx, info.FullMethod, extractor)

		if err != nil {
			return handler(ctx, req)
		}

		span.Tag(go2sky.TagHTTPMethod, "GRPC")
		span.Tag(go2sky.TagURL, info.FullMethod)

		defer func() {
			if r := recover(); r != nil {
				err = status.Errorf(codes.Internal, "panic: %v", r)
				span.Error(time.Now(), err.Error())
				span.Tag("error", "true")
			}
			if err != nil {
				span.Error(time.Now(), err.Error())
				span.Tag("error", "true")
			}
			span.End()
		}()

		resp, err = handler(ctx, req)
		return resp, err
	}
}

func (s *SkyWalkingConfig) ClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md, ok := metadata.FromOutgoingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		extractor := func(headerKey, headerValue string) error {
			md.Set(headerKey, headerValue)
			ctx = metadata.NewOutgoingContext(ctx, md)
			return nil
		}

		span, err := s.tracer.CreateExitSpan(ctx, method, cc.Target(), extractor)

		if err != nil {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		span.Tag(go2sky.TagHTTPMethod, "GRPC")
		span.Tag(go2sky.TagURL, method)

		defer func() {
			if r := recover(); r != nil {
				err = status.Errorf(codes.Internal, "panic: %v", r)
				span.Error(time.Now(), err.Error())
				span.Tag("error", "true")
			}
			if err != nil {
				span.Error(time.Now(), err.Error())
				span.Tag("error", "true")
			}
			span.End()
		}()

		err = invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}

func (s *SkyWalkingConfig) GinMiddleware(r *gin.Engine) gin.HandlerFunc {
	return sgin.Middleware(r, s.tracer)
}

func (s *SkyWalkingConfig) Close() error {
	s.reporter.Close()
	return nil
}
