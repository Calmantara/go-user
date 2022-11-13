package creditcard

import "github.com/Calmantara/go-user/lib/entity"

var (
	CreditCardTypeMap = map[string]bool{
		"GOLD":     true,
		"PREMIUM":  true,
		"PLATINUM": true,
	}
)

type CreditCardToken struct {
	Id     uint64 `json:"id"`
	UserId uint64 `json:"user_id"`
	Token  string `json:"token"`
	entity.DefaultColumn
}
type CreditCard struct {
	Cvv     string `json:"cvv,omitempty" binding:"required"`
	Expired string `json:"expired" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Number  string `json:"number" binding:"required"`
	Type    string `json:"type" binding:"required"`
}

type CreditClaim struct {
	Subject    string     `json:"sub"`
	Issuer     string     `json:"iss"`
	Audience   string     `json:"aud"`
	Type       string     `json:"type"`
	IssuedAt   int64      `json:"iat"`
	CreditCard CreditCard `json:"creditCard"`
}
