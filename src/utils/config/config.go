package config

type ConfigApp struct {
	Port int
}

func NewConfigApp() *ConfigApp {
	return &ConfigApp{
		Port: 8000,
	}
}
