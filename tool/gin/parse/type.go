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
		Tag         string
		Produce     string
		Accept      string
		Param       []string
		Success     string
		Failure     string
		Router      string
	}
	Method struct {
		Method string
		Route  string
		Req    string
		Resp   string
	}
)
