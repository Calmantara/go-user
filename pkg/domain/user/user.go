package user

import (
	"strings"

	"github.com/Calmantara/go-user/lib/entity"
	creditcard "github.com/Calmantara/go-user/pkg/domain/credit-card"
	"github.com/Calmantara/go-user/pkg/domain/photo"
)

type User struct {
	Hased           bool                        `json:"-" gorm:"-"`
	PassPhoto       bool                        `json:"-" gorm:"-"`
	Id              uint64                      `json:"id"`
	Name            string                      `json:"name"`
	Address         string                      `json:"address"`
	Email           string                      `json:"email" binding:"required"`
	Password        string                      `json:"password,omitempty" `
	Photos          []*photo.Photo              `json:"photos,omitempty" gorm:"foreignKey:UserId;references:Id"`
	CreditCardToken *creditcard.CreditCardToken `json:"-" gorm:"foreignKey:UserId;references:Id"`
	CreditCard      *creditcard.CreditCard      `json:"creditcard,omitempty" gorm:"-"`
	entity.DefaultColumn
}

type UserQuery struct {
	Q  string `form:"q"`
	Ob string `form:"ob,default=name"`
	Sb string `form:"sb,default=asc"`
	Of int    `form:"of,default=0"`
	Lt int    `form:"lt,default=30"`
}

var (
	validOrder = map[string]bool{
		"name":  true,
		"email": true}
	validSort = map[string]bool{
		"asc":  true,
		"desc": true}
)

func (u *UserQuery) ValidateForm() bool {
	u.Sb = strings.ToLower(u.Sb)
	if !validOrder[u.Ob] || !validSort[u.Sb] || u.Lt < 0 || u.Of < 0 {
		return false
	}
	if u.Q != "" {
		u.Q = "%" + u.Q + "%"
	}
	return true
}

func (User) TableName() string {
	return "users"
}
