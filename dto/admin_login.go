package dto

import (
	"github.com/e421083458/go_gateway/public"
	"github.com/gin-gonic/gin"
	"time"
)

type AdminSessionInfo struct {
	ID        int       `json:"id"`
	UserName  string    `json:"username"`
	LoginTime time.Time `json:"login_time"`
}

type AdminLoginInput struct {
	UserName string `json:"username" form:"username" commnet:"姓名" example:"admin" validate:"required,is_valid_username"` //管理员用户名
	PassWord string `json:"password" form:"password" commnet:"密码" example:"123456" validate:"required"`                  //密码
}

func (param *AdminLoginInput) BindValidParam(c *gin.Context) error {
	return public.DefaultGetValidParams(c, param)
}

type AdminLoginOutput struct {
	Token string `json:"token" form:"token" commnet:"token" example:"token" validate:""`
}
