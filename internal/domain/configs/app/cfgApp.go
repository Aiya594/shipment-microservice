package cfgApp

import "os"

type ConfigApp struct {
	Port string
}

func LoadConfigApp() *ConfigApp {
	return &ConfigApp{Port: os.Getenv("PORT")}
}
