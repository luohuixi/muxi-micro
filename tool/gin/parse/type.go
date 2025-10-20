package parse

type Api struct {
	T           string // type
	ServiceName string
	Service     []*Service
	Server      *Server
}

type (
	Server struct {
		Prefix string
		Group  string
	}
	Service struct {
		Doc     *Doc
		Handler string
		Method  *Method
	}
	Doc struct {
		Summary     string
		Description string
		Accept      string
		Produce     string
	}
	Method struct {
		Method string
		Route  string
		Req    string
		Resp   string
	}
)
