package models

type HostConfig struct {
	Host           string
	Port           string
	BuildEnv       string
	AllowedOrigins string
}

func NewHostConfig(host string, port string, BuildEnv string) *HostConfig {

	return &HostConfig{
		Host:     host,
		Port:     port,
		BuildEnv: BuildEnv,
	}
}
