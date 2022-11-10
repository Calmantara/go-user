package creditcard

type CreditCard struct {
	Cvv     string `json:"cvv"`
	Expired string `json:"expired"`
	Name    string `json:"name"`
	Number  string `json:"number"`
	Type    string `json:"type"`
}
