package user

import (
	creditcard "github.com/Calmantara/go-user/pkg/domain/credit-card"
	"github.com/Calmantara/go-user/pkg/domain/photo"
)

type User struct {
	Id         uint64                 `json:"id"`
	Name       string                 `json:"name" binding:"required"`
	Address    string                 `json:"address" binding:"required"`
	Email      string                 `json:"email" binding:"required"`
	Password   string                 `json:"password" binding:"required"`
	Photos     []*photo.Photo         `json:"photos,omitempty" binding:"required" gorm:"foreignKey:Id;references:UserId"`
	CreditCard *creditcard.CreditCard `json:"creditcard,omitempty" binding:"required" gorm:"foreignKey:Id;references:UserId"`
}

type UserQuery struct {
}

func (User) TableName() string {
	return "users"
}
