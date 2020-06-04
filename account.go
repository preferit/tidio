package tidio

func NewAccount(username, role string) *Account {
	return &Account{
		Username: username,
		Role:     role,
	}
}

type Account struct {
	Username string
	Role     string
}
