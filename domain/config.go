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
}

type ConfigService interface {
	GetConfig() Config
}
