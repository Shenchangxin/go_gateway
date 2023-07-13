package controller

import (
	"encoding/json"
	"github.com/e421083458/go_gateway/dao"
	"github.com/e421083458/go_gateway/dto"
	"github.com/e421083458/go_gateway/middleware"
	"github.com/e421083458/go_gateway/public"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"time"
)

type AdminLoginController struct {
}

func AdminLoginRegister(group *gin.RouterGroup) {
	adminLogin := &AdminLoginController{}
	group.POST("/login", adminLogin.AdminLogin)
	group.GET("/login_out", adminLogin.AdminLoginOut)
}

//AdminLogin godoc
//@Summary 管理员登陆
//@Description 管理员登陆
//@Tags 管理员接口
//@ID /admin_login/login
//@Accept json
//@Produce json
//@Param body body dto.AdminLoginInput true "body"
//@Success 200 {object} minddleware.Response{data=dto.AdminLoginOutput}
//"success"
//@Router /admin_login/login [post]

func (adminlogin *AdminLoginController) AdminLogin(c *gin.Context) {
	params := &dto.AdminLoginInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	//1. params.UserName 取得管理员信息admininfo
	//2. admininfo.salt+param.PassWord sha256 => saltPassWord
	//3. saltPassWord==admininfo.password
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	admin := &dao.Admin{}
	admin, err = admin.LoginCheck(c, tx, params)
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}

	//设置session
	sessInfo := &dto.AdminSessionInfo{
		ID:        admin.Id,
		UserName:  admin.UserName,
		LoginTime: time.Now(),
	}
	sessBts, err := json.Marshal(sessInfo)
	if err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}
	sess := sessions.Default(c)
	sess.Set(public.AdminInfoSessionKey, sessBts)
	sess.Save()

	out := &dto.AdminLoginOutput{Token: params.UserName}
	middleware.ResponseSuccess(c, out)
}

//AdminLoginOut godoc
//@Summary 管理员退出登陆
//@Description 管理员退出登陆
//@Tags 管理员接口
//@ID /admin_login/login_out
//@Accept json
//@Produce json
//@Success 200 {object} minddleware.Response{data=string}
//"success"
//@Router /admin_login/login_out [get]

func (adminlogin *AdminLoginController) AdminLoginOut(c *gin.Context) {
	sess := sessions.Default(c)
	sess.Delete(public.AdminInfoSessionKey)
	sess.Save()
	middleware.ResponseSuccess(c, "")
}
