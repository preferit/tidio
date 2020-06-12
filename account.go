package tidio

func NewAccount(username, role string) *Account {
	return &Account{
		Username: username,
	}
}

type Account struct {
	Username string
}

// ----------------------------------------

type Accounts interface {
	LoadAccountByKey(*Account, string) error
}

// ----------------------------------------

type APIKeys map[string]*Account
