package mysql

type Options struct {
	User               string
	Password           string
	Host               string
	Net                string
	DB                 string
	MaxFileTimeMinutes int
	MaxOpenConns       int
}

func NewOption(user, password, host, db string) Options {
	return Options{User: user, Password: password, Host: host, DB: db, Net: "tcp", MaxFileTimeMinutes: 3, MaxOpenConns: 10}
}
