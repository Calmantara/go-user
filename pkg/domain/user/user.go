package user

import (
	"github.com/Calmantara/go-user/common/entity"
	creditcard "github.com/Calmantara/go-user/pkg/domain/credit-card"
	"github.com/Calmantara/go-user/pkg/domain/photo"
)

type User struct {
	Id              uint64                      `json:"id"`
	Name            string                      `json:"name" binding:"required"`
	Address         string                      `json:"address" binding:"required"`
	Email           string                      `json:"email" binding:"required"`
	Password        string                      `json:"password,omitempty" binding:"required"`
	Photos          []*photo.Photo              `json:"photos,omitempty" binding:"required" gorm:"foreignKey:UserId;references:Id"`
	CreditCardToken *creditcard.CreditCardToken `json:"-" gorm:"foreignKey:UserId;references:Id"`
	CreditCard      *creditcard.CreditCard      `json:"creditcard,omitempty" gorm:"-"`
	entity.DefaultColumn
}

type UserQuery struct {
}

func (User) TableName() string {
	return "users"
}
