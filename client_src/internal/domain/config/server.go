package config

type Config struct {
	Servers []ServerConfig `yaml:"servers"`
	Watcher WatcherConfig  `yaml:"watcher,omitempty"`
}

type ServerConfig struct {
	Address string       `yaml:"address"`
	Users   []UserConfig `yaml:"users"`
}

type UserConfig struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type WatcherConfig struct {
	Reconnect ReconnectConfig `yaml:"reconnect,omitempty"`
}

type ReconnectConfig struct {
	Timeout int `yaml:"timeout"`
	Retries int `yaml:"retries"`
}
