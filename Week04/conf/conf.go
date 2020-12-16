package conf

type Config struct {
	Addr string
	Dsn  string
}

func Load() Config {
	return Config{}
}
