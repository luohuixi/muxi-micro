package ginx

import (
	"github.com/gin-gonic/gin"
	t_http "github.com/muxi-Infra/muxi-micro/pkg/transport/http"
	"net/http"
	"testing"
)

func TestHandler(t *testing.T) {
	g := NewDefaultEngine()
	api := g.Group("/v1")
	RegisterUsersv1(api, &UserController{})
}

type UserController struct{}

type CreateUserReq struct {
	UserName string `json:"user_name"`
}
type CreateUserResp struct {
	UserName string `json:"user_name"`
}

func RegisterUsersv1(g *gin.RouterGroup, c User) {
	api := g.Group("/user")
	api.GET("/get", WrapReq(c.create))
}

type User interface {
	create(ctx *gin.Context, req CreateUserReq) t_http.Response
}

//利用@AutoGen标识符来保证重新生成的时候不会将代码内部逻辑破坏
func (c *UserController) create(ctx *gin.Context, req CreateUserReq) t_http.Response {
	//填写你的逻辑
	return t_http.Response{
		HttpCode: http.StatusOK,
		Code:     0,
		Message:  "",
		Data:     CreateUserResp{},
	}
}
