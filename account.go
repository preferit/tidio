package tidio

func NewAccount(username, role string) *Account {
	return &Account{
		Username: username,
	}
}

type Account struct {
	Username string
}
