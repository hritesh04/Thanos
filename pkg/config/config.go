package config

type Config struct {
	Algorithm string `yaml:"type"`
	// HealthCheckInterval string   `yaml:"healthCheckInterval"`
	Servers    []string `yaml:"backends"`
	ListenPort string   `yaml:"port"`
}

type Backend struct {
	Url string `yaml:"url"`
}
