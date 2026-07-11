package config

import (
	"fmt"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

func (c *Config) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.Host, c.Port, c.User, c.Password, c.DBName)
}
