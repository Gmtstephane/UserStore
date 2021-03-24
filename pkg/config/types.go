package config

type Config struct {
	Database Dbconfig    `yaml:"database"`
	Redis    RedisConfig `yaml:"redis"`
}

type Dbconfig struct {
	Name     string `yaml:"name"`
	UserName string `yaml:"user_name"`
	Uri      string `yaml:"uri"`
	Port     uint   `yaml:"port"`
	Password string `yaml:"password"`
}

type RedisConfig struct {
	Uri  string `yaml:"uri"`
	Port uint   `yaml:"port"`
}
