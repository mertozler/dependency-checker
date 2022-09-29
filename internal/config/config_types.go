package config

type Config struct {
	Redis      *Redis
	Host       *Host
	Mail       *Mail
	Github     *Github
	Javascript *Javascript
	Dependency *Dependency
	MailSender *MailSender
}

type Redis struct {
	Host     string
	Password string
}

type Host struct {
	Port string
}

type MailSender struct {
	Hour int
}
type Javascript struct {
	RegistryURL string
}

type Github struct {
	API string
}

type Mail struct {
	From     string
	Password string
	Host     string
	Port     int
}

type Dependency struct {
	Path []string
}
