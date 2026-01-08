package exporter

type Config struct {
	Type   string         `yaml:"type"`
	Config map[string]any `yaml:"config"`
}
