package meter

type Config struct {
	Name  string `yaml:"name"`
	Slave uint8  `yaml:"slave"`
	Type  string `yaml:"type"`
}
