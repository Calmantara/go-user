package user

import (
	creditcard "github.com/Calmantara/go-user/pkg/domain/credit-card"
	"github.com/Calmantara/go-user/pkg/domain/photo"
)

type User struct {
	Id         uint64                 `json:"id"`
	Name       string                 `json:"name"`
	Address    string                 `json:"address"`
	Email      string                 `json:"email"`
	Password   string                 `json:"password"`
	Photos     []*photo.Photo         `json:"photos,omitempty"`
	CreditCard *creditcard.CreditCard `json:"creditcard,omitempty"`
}

type UserQuery struct {
}
