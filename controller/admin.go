package controller

import (
	"encoding/json"
	"fmt"
	"github.com/e421083458/go_gateway/dao"
	"github.com/e421083458/go_gateway/dto"
	"github.com/e421083458/go_gateway/middleware"
	"github.com/e421083458/go_gateway/public"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type AdminController struct {
}

func AdminRegister(group *gin.RouterGroup) {
	adminLogin := &AdminController{}
	group.GET("/admin_info", adminLogin.AdminInfo)
	group.GET("/change_pwd", adminLogin.ChangePwd)
}

//AdminInfo godoc
//@Summary 管理员信息
//@Description 管理员信息
//@Tags 管理员接口
//@ID /admin/admin_info
//@Accept json
//@Produce json
//@Param body body dto.AdminLoginInput true "body"
//@Success 200 {object} minddleware.Response{data=dto.AdminInfoOutput}
//"success"
//@Router /admin/admin_info [get]

func (admin *AdminController) AdminInfo(c *gin.Context) {

	sess := sessions.Default(c)
	sessInfo := sess.Get(public.AdminInfoSessionKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), adminSessionInfo); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	//1. params.UserName 取得管理员信息admininfo
	//2. admininfo.salt+param.PassWord sha256 => saltPassWord

	out := &dto.AdminInfoOutput{
		ID:           adminSessionInfo.ID,
		Name:         adminSessionInfo.UserName,
		LoginTime:    adminSessionInfo.LoginTime,
		Avatar:       "",
		Introduction: "I am a super administrator",
		Roles:        []string{"admin"},
	}
	middleware.ResponseSuccess(c, out)
}

//ChangePwd godoc
//@Summary 修改密码
//@Description 修改密码
//@Tags 管理员接口
//@ID /admin/change_pwd
//@Accept json
//@Produce json
//@Param body body dto.AdminLoginInput true "body"
//@Success 200 {object} minddleware.Response{data=string} "success"
//@Router /admin/change_pwd [post]

func (admin *AdminController) ChangePwd(c *gin.Context) {

	params := &dto.ChangePwdInput{}
	if err := params.BindValidParam(c); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	//1. session 读取用户信息到结构体 sessInfo
	//2. sessInfo.ID读取数据库信息adminInfo
	//3. params.password+adminInfo.salt sha256 saltPassword
	//4. saltPassword==>adminInfo.password执行数据保存

	//session 读取用户信息到结构体 sessInfo
	sess := sessions.Default(c)
	sessInfo := sess.Get(public.AdminInfoSessionKey)
	adminSessionInfo := &dto.AdminSessionInfo{}
	if err := json.Unmarshal([]byte(fmt.Sprint(sessInfo)), adminSessionInfo); err != nil {
		middleware.ResponseError(c, 2000, err)
		return
	}
	//从数据库中读取adminInfo
	tx, err := lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c, 2001, err)
		return
	}
	adminInfo := &dao.Admin{}
	adminInfo, err = adminInfo.Find(c, tx, (&dao.Admin{UserName: adminSessionInfo.UserName}))
	if err != nil {
		middleware.ResponseError(c, 2002, err)
		return
	}
	//生成加密后的密码并保存数据库
	saltPassword := public.GenSaltPassword(adminInfo.Salt, params.Password)
	adminInfo.PassWord = saltPassword
	if err := adminInfo.Save(c, tx); err != nil {
		middleware.ResponseError(c, 2003, err)
		return
	}

	middleware.ResponseSuccess(c, "")
}
