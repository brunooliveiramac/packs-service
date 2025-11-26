package config

import "os"

type DBConfig struct {
	URL      string
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func DefaultDB() DBConfig {
	return DBConfig{
		Host:     "localhost",
		Port:     "5435",
		User:     "postgres",
		Password: "postgres",
		Database: "packs",
		SSLMode:  "disable",
	}
}

func LoadDB() DBConfig {
	cfg := DefaultDB()
	if v := os.Getenv("DATABASE_URL"); v != "" { cfg.URL = v }
	if v := os.Getenv("PGHOST"); v != "" { cfg.Host = v }
	if v := os.Getenv("PGPORT"); v != "" { cfg.Port = v }
	if v := os.Getenv("PGUSER"); v != "" { cfg.User = v }
	if v := os.Getenv("PGPASSWORD"); v != "" { cfg.Password = v }
	if v := os.Getenv("PGDATABASE"); v != "" { cfg.Database = v }
	if v := os.Getenv("PGSSLMODE"); v != "" { cfg.SSLMode = v }
	return cfg
}


