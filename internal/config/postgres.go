package config

import "fmt"

type Postgres struct {
	Host     string
	Port     int
	User     string
	Password string // todo: to .env?
	DBName   string
}

func (p *Postgres) ConnectionDSN() string {
	return fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable",
		p.Host, p.Port, p.User, p.Password, p.DBName)
}

func (p *Postgres) ConnectionURL() string {
	return fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=disable", p.User, p.Password, p.Host, p.Port, p.DBName)
}
