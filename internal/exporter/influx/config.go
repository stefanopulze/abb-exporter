package influx

type Options struct {
	Host         string `yaml:"host" env:"INFLUX_HOST"`
	Token        string `yaml:"token" env:"INFLUX_TOKEN"`
	Database     string `yaml:"database" env:"INFLUX_DATABASE"`
	Organization string `yaml:"organization" env:"INFLUX_ORGANIZATION"`
}

func OptionsFromMap(m map[string]any) Options {
	return Options{
		Host:         m["host"].(string),
		Token:        m["token"].(string),
		Database:     m["database"].(string),
		Organization: m["organization"].(string),
	}
}
