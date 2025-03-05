package config

type Config struct {
	Servers []ServerConfig `yaml:"servers"`
}

type ServerConfig struct {
	Address string       `yaml:"address"`
	Users   []UserConfig `yaml:"users"`
}

type UserConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}
