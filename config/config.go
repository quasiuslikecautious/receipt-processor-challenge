package config

type Config struct {
  Debug bool
  Host  string
  Port  string
}

func DefaultConfig() *Config {
  return &Config{
    Debug:  false,
    Host:   "127.0.0.1",
    Port:   "8080",
  }
}
