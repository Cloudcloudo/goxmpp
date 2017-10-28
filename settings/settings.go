package settings

import (
	"fmt"
	"os"
)

type Settings struct {
	Domain 			string
	WebPort			string
	XMPPPort		string
	ListenAddress	string

	DbName			string
	DbUser			string
	DbPassword		string
	DbHost			string
	DbPort			string
	DbSSL			string
}

func GetSettings() Settings {
	settings := Settings{
		Domain: os.Getenv("DOMAIN"),
		WebPort: os.Getenv("WEB_PORT"),
		XMPPPort: os.Getenv("XMPP_PORT"),
		ListenAddress: os.Getenv("ADDRESS"),

		DbName: os.Getenv("DB_NAME"),
		DbUser: os.Getenv("DB_USER"),
		DbPassword: os.Getenv("DB_PASS"),
		DbHost: os.Getenv("DB_HOST"),
		DbPort: os.Getenv("DB_PORT"),
		DbSSL: os.Getenv("DB_SSL"),
	}
	if settings.Domain == "" {
		settings.Domain = "localhost"
	}
	if settings.WebPort == "" {
		settings.WebPort = "8080"
	}
	if settings.XMPPPort == "" {
		settings.XMPPPort = "5222"
	}
	if settings.ListenAddress == "" {
		settings.ListenAddress = "0.0.0.0"
	}
	if settings.DbName == "" {
		settings.DbName = "goxmpp"
	}
	if settings.DbUser == "" {
		settings.DbUser = "goxmpp_user"
	}
	if settings.DbPassword == "" {
		settings.DbPassword = "aZ82w2E-aXwNch5"
	}
	if settings.DbHost == "" {
		settings.DbHost = "localhost"
	}
	if settings.DbPort == "" {
		settings.DbPort = "5432"
	}
	if settings.DbSSL == "" {
		settings.DbSSL = "disable"
	}

	return settings
}

func GetDBConnectionString() string {
	c := GetSettings()

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DbUser, c.DbPassword, c.DbHost, c.DbPort, c.DbName, c.DbSSL)
}