package testutils

type Config struct {
	Host   string `default:"localhost"`
	Port   int    `default:"3000"`
	Sentry Sentry
}

type Sentry struct {
	Url    string `default:"https://xx.xx.xx"`
	Labels map[string]string
}
