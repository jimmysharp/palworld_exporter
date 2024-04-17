package config

type Config struct {
	ListenAddress string
	ScrapeURI     string
	HttpUsername  string
	HttpPassword  string

	LogLevel  string
	LogFormat string
}
