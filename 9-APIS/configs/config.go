package configs

import (
	"github.com/go-chi/jwtauth"
	"github.com/spf13/viper"
)

type conf struct {
	DBDriver      string `mapstructure:"DBDriver"`
	DBHost        string `mapstructure:"DBHost"`
	DBPort        string `mapstructure:"DBPort"`
	DBUser        string `mapstructure:"DBUser"`
	DBPassword    string `mapstructure:"DBPassword"`
	DBName        string `mapstructure:"DBName"`
	WebServerPort string `mapstructure:"WebServerPort"`
	JWTSecret     string `mapstructure:"JWTSecret"`
	JWTExpiresIn  int    `mapstructure:"JWTExpiresIn"`
	TokenAuth     *jwtauth.JWTAuth
}

func LoadConfig(path string) (*conf, error) {
	var cfg *conf
	viper.SetConfigName("app_config")
	viper.SetConfigType("env")
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&cfg)
	if err != nil {
		panic(err)
	}
	cfg.TokenAuth = jwtauth.New("HS256", []byte(cfg.JWTSecret), nil)

	return cfg, nil
}
