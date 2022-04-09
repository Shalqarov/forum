package configs

type Config struct {
	BindAddr string `toml:"bind_addr"`

	HostDB     string `toml:"host_db"`
	PortDB     string `toml:"port_db"`
	NameDB     string `toml:"name_db"`
	UserDB     string `toml:"user_db"`
	PasswordDB string `toml:"password_db"`

	CacheHost     string `toml:"cache_host"`
	CachePassword string `toml:"cache_password"`
	CacheAddr     string `toml:"cache_addr"`
}

func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",

		HostDB:     "sqlite3",
		PortDB:     ":5432",
		NameDB:     "forum",
		UserDB:     "sqlite3",
		PasswordDB: "",

		CacheHost:     "localhost",
		CachePassword: "",
		CacheAddr:     ":6379",
	}
}
