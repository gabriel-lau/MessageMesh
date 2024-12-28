package models

type Account struct {
	Username  string `json:"username"`
	PublicKey string `json:"publicKey"`
}

func (a Account) GetType() string {
	return "Account"
}
