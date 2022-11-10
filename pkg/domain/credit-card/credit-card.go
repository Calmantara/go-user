package creditcard

type CreditCard struct {
	Cvv     string `json:"cvv" binding:"required"`
	Expired string `json:"expired" binding:"required"`
	Name    string `json:"name" binding:"required"`
	Number  string `json:"number" binding:"required"`
	Type    string `json:"type" binding:"required"`
}
