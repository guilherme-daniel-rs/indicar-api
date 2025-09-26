package configs

import (
	"fmt"
	"os"
	"reflect"

	"github.com/spf13/viper"
)

var configuration *Config

type Config struct {
	Database database
	JWT      jwt
	AWS      aws
}

type database struct {
	User     string `mapstructure:"DB_USER"`
	Password string `mapstructure:"DB_PASSWORD"`
	Host     string `mapstructure:"DB_HOST"`
	Port     string `mapstructure:"DB_PORT"`
	Name     string `mapstructure:"DB_NAME"`
}

type jwt struct {
	Secret string `mapstructure:"JWT_SECRET" default:"your-secret-key"`
}

type aws struct {
	Region          string `mapstructure:"AWS_REGION" default:"us-east-1"`
	AccessKeyID     string `mapstructure:"AWS_ACCESS_KEY_ID"`
	SecretAccessKey string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	S3Bucket        string `mapstructure:"AWS_S3_BUCKET" default:"indicar-evaluation-photos"`
}

func getMappedEnvs(configStruct reflect.Type) []string {
	result := make([]string, 0)

	for i := 0; i < configStruct.NumField(); i++ {
		field := configStruct.Field(i)
		if configName := field.Tag.Get("mapstructure"); configName != "" {
			result = append(result, configName)
		}
		if field.Type.Kind() == reflect.Struct {
			result = append(result, getMappedEnvs(field.Type)...)
		}
	}
	return result
}

func setDefaultValues(configStruct reflect.Type) {
	for i := 0; i < configStruct.NumField(); i++ {
		field := configStruct.Field(i)
		configName := field.Tag.Get("mapstructure")
		defaultValue := field.Tag.Get("default")

		if configName != "" && defaultValue != "" {
			viper.SetDefault(configName, defaultValue)
		}

		if field.Type.Kind() == reflect.Struct {
			setDefaultValues(field.Type)
		}
	}
}

func Load() error {
	configuration = &Config{}

	environment := os.Getenv("GO_ENV")
	if environment == "" {
		fmt.Println("[Method: Config.Load()] Your GO_ENV was not filled, configure it on environment or env file and try again.")
	}
	envFile := ".env-"
	envPath := "."
	if environment == "" {
		envFile += "development"
	} else if environment == "test" {
		envFile = envFile + environment
		envPath = "../../test/"
	} else {
		envFile += environment
	}

	viper.AddConfigPath(envPath)
	viper.SetConfigName(envFile)
	viper.SetConfigType("env")

	setDefaultValues(reflect.TypeOf(Config{}))

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("[Method: Config.Load()]", envFile, "not found, load by environment variables")

			viper.AutomaticEnv()
			mapped := getMappedEnvs(reflect.TypeOf(Config{}))
			for _, env := range mapped {
				viper.BindEnv(env)
			}
		} else {
			return err
		}
	} else {
		fmt.Println("[Method: Config.Load()] Using config file:", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&configuration); err != nil {
		return err
	}

	if err := viper.Unmarshal(&configuration.Database); err != nil {
		return err
	}

	return nil
}

func Get() *Config {
	return configuration
}
