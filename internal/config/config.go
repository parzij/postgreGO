package config

import (
	"fmt"
	"os"
)

type Config struct {
	HostDB string
	PortDB string
	UserDB string
	PasswordDB string
	NameDB string
	SSLmodeDB string
}

func LoadConfig() (*Config, error) {
	config := &Config {
		HostDB: os.Getenv("HOST_DB"),
		PortDB: os.Getenv("PORT_DB"),
		UserDB: os.Getenv("USER_DB"),
		PasswordDB: os.Getenv("PASSWORD_DB"),
		NameDB: os.Getenv("NAME_DB"),
		SSLmodeDB: os.Getenv("SSLMODE_DB"),
	}
	
	if config.HostDB == "" ||
		config.PortDB == "" ||
		config.UserDB == "" ||
		config.PasswordDB == "" ||
		config.NameDB == "" ||
		config.SSLmodeDB == "" {
			return nil, fmt.Errorf("one or more required parametrs are missing")
		}

	return config, nil
}