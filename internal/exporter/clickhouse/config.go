package clickhouse

type Options struct {
	Addr     string `yaml:"host" env:"CLICKHOUSE_HOST"`
	Database string `yaml:"database" env:"CLICKHOUSE_DATABASE"`
	Username string `yaml:"username" env:"CLICKHOUSE_USERNAME"`
	Password string `yaml:"password" env:"CLICKHOUSE_PASSWORD"`
}

func OptionsFromMap(m map[string]interface{}) Options {
	return Options{
		Addr:     m["host"].(string),
		Database: m["database"].(string),
		Username: m["username"].(string),
		Password: m["password"].(string),
	}
}
