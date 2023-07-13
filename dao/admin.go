package dao

import (
	"errors"
	"github.com/e421083458/go_gateway/dto"
	"github.com/e421083458/go_gateway/public"
	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
	"time"
)

type Admin struct {
	Id        int       `json:"id" gorm:"primary_key" description:"自增主键"`
	UserName  string    `json:"user_name" gorm:"column:user_name" description:"管理员用户名"`
	Salt      string    `json:"salt" gorm:"column:salt" description:"盐"`
	PassWord  string    `json:"pass_word" gorm:"column:password" description:"密码"`
	UpdateAt  time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at" description:"创建时间"`
	IsDelete  int       `json:"is_delete" gorm:"column:is_delete" description:"是否删除"`
}

func (t *Admin) TableName() string {
	return "gateway_admin"
}

func (t *Admin) LoginCheck(c *gin.Context, tx *gorm.DB, param *dto.AdminLoginInput) (*Admin, error) {
	adminInfo, err := t.Find(c, tx, (&Admin{UserName: param.UserName, IsDelete: 0}))
	if err != nil {
		return nil, errors.New("用户信息不存在")
	}
	saltPassWord := public.GenSaltPassword(adminInfo.Salt, param.PassWord)
	if adminInfo.PassWord != saltPassWord {
		return nil, errors.New("密码错误，请重新输入")
	}
	return nil, nil
}

func (t *Admin) Find(c *gin.Context, tx *gorm.DB, search *Admin) (*Admin, error) {
	out := &Admin{}
	err := tx.SetCtx(public.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (t *Admin) Save(c *gin.Context, tx *gorm.DB) error {

	return tx.SetCtx(public.GetGinTraceContext(c)).Save(t).Error

}
