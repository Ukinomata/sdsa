package config

import "fmt"

type PostgresConfig struct {
	user     string `yaml:"user" env-required:"true"`
	password string `yaml:"password" env-default:""`
	host     string `yaml:"host" env-default:"localhost"`
	dbname   string `yaml:"dbname" env-required:"true"`
	sslmode  string `yaml:"sslmode" env-default:"disable"`
}

func (config *PostgresConfig) GetDataSourceName() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=%s",
		config.user,
		config.password,
		config.host,
		config.dbname,
		config.sslmode)
}
