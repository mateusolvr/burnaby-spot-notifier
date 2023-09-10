package domain

type Config struct {
	Email struct {
		From string   `yaml:"from"`
		Pass string   `yaml:"pass"`
		To   []string `yaml:"to"`
	} `yaml:"email"`
	ActivityName         string `yaml:"activityName"`
	RecreationalCentreId int    `yaml:"recreationalCentreId"`
	Redis                struct {
		Enabled       bool   `yaml:"enabled"`
		ExpireMinutes int    `yaml:"expireMinutes"`
		Url           string `yaml:"url"`
	} `yaml:"redis"`
}

type ConfigService interface {
	GetConfig() Config
}
