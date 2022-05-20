package domain

type Config struct {
	Email struct {
		From string `yaml:"from"`
		Pass string `yaml:"pass"`
		To   string `yaml:"to"`
	} `yaml:"email"`
	Activity struct {
		Name      string `yaml:"name"`
		DaysAhead int    `yaml:"daysAhead"`
	} `yaml:"activity"`
	Redis struct {
		Enabled       bool   `yaml:"enabled"`
		ExpireMinutes int    `yaml:"expireMinutes"`
		Url           string `yaml:"url"`
	} `yaml:"redis"`
}

type ConfigService interface {
	GetConfig() Config
}
