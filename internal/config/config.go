package config

import "os"

type Config struct {
	Host string
	Port string
}

func Load() Config {
	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return Config{
		Host: host,
		Port: port,
	}
}

func (c Config) Address() string {
	return c.Host + ":" + c.Port
}
