package config

type Config struct {
	Algorithm string `yaml:"type"`
	// HealthCheckInterval string   `yaml:"healthCheckInterval"`
	Servers    []Server `yaml:"backends"`
	ListenPort string   `yaml:"port"`
}

type Server struct {
	Url            string `yaml:"url"`
	HealthEndPoint string `yaml:"health_end_point"`
}
