package config

type mail struct {
	Host     string
	Port     int
	Username string
	Password string
	FromAddr string
	FromMame string
}

var Mail mail

func loadMailConfig() {
	Mail = mail{
		Host:     GetString("mail.host"),
		Port:     GetInt("mail.port"),
		Username: GetString("mail.username"),
		Password: GetString("mail.password"),
		FromAddr: GetString("mail.from_addr"),
		FromMame: GetString("mail.from_name"),
	}
}
