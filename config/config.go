package config

import (
	"os"
	"time"

	log15 "github.com/inconshreveable/log15"

	// "time"
	"github.com/jinzhu/configor"
)

type (
	// Config Default
	Config struct {
		Server struct {
			Host     string `default:"0.0.0.0" env:"LISTEN_ADDRESS"`
			Port     uint   `default:"8080" env:"LISTEN_PORT"`
			Loglevel string `default:"info" json:"loglevel"`
		}
		Google struct {
			Ios struct {
				ClientID     string `json:"client_id" env:"GOOGLE_IOS_CLIENT_ID"`
				ClientSecret string `json:"client_secret" env:"GOOGLE_IOS_CLIENT_SECRET"`
			}
			Website struct {
				ClientID     string `json:"client_id" env:"GOOGLE_WEBSITE_CLIENT_ID"`
				ClientSecret string `json:"client_secret" env:"GOOGLE_WEBSITE_CLIENT_SECRET"`
			}
		}
		Redis struct {
			Host      string `default:"localhost" env:"REDISHOST"`
			Port      uint   `default:"6379"`
			Password  string
			MaxIdle   int `default:"50"`
			MaxActive int `default:"100"`
			DB        int `default:"1"`
		}
		MongoDB struct {
			User       string        `env:"MONGOUSER"`
			Password   string        `env:"MONGOPASS"`
			Hosts      []interface{} `default:"[\"mongodb:27017\"]" env:"MONGOHOST"`
			Database   string        `default:"casbu" env:"MONGODATABASE"`
			SSLEnabled bool          `json:"sslEnabled"`
		}
		Rabbit struct {
			Host     string `default:"localhost" env:"RABBITMQHOST"`
			Port     uint   `default:"5672"`
			User     string `default:"guest"`
			Password string `default:"guest"`
			Vhost    string `default:""`
			Q        map[string]string
		}
		Application struct {
			URL     string `env:"APP_URL" json:"url"`
			AppRole string `env:"APP_ROLE" json:"role"`
			AppName string `env:"APP_NAME" json:"name"`
		}
		Gcs struct {
			CredentialPath string `env:"GCS_CRED_PATH" json:"credential_path"`
			BucketName     string `env:"GCS_BUCKET_NAME" json:"bucket_name"`
		}
		Gpubsub struct {
			CredentialPath string `env:"GPUBSUB_CRED_PATH" json:"credential_path"`
			ProjectName    string `env:"GPUBSUB_PROJECT_NAME" json:"project_name"`
			Topic          string `env:"GPUBSUB_TOPIC" json:"topic"`
			Subscribe      string `env:"GPUBSUB_SUBSCRIBE" json:"subscribe"`
		}
		Gaurun struct {
			Url    string `env:"GAURUN_URL" json:"url"`
			UrlDev string `env:"GAURUN_URL_DEV" json:"url_dev"`
		}
		Mailgun struct {
			APIKey string `env:"MAILGUN_APIKEY" json:"apiKey"`
			Domain string `env:"MAILGUN_DOMAIN" json:"domain"`
		}
		Chatgpt struct {
			ChatGPTURL     string `env:"CHAT_GPT_URL" json:"chat_gpturl"`
			ChatGPTSkipTLS bool   `env:"CHAT_GPT_SKIP_TLS" json:"chat_gpt_skip_tls"`
			ChatGPTToken   string `env:"CHAT_GPT_TOKEN" json:"chat_gpt_token"`
		}
		Salt string `env:"SALT" json:"salt"`
	}
)

// Load : Reads the config from relevant env json file.
func Load() (config Config) {
	configPath := "config/config.json"
	golang_env := os.Getenv("GOLANG_ENV")
	switch golang_env {
	case "local":
		configPath = "config/config.local.json"
	case "dev", "prod", "staging", "test":
		if golang_env == "prod" {
			log15.Info("LOADING PROUDCTION CONFIGS", "args", os.Args)
		}
		os.Setenv("CONFIGOR_ENV", golang_env)
	default:
		panic("please provide GOLANG_ENV")
	}
	log15.Info("Loading Configuration For :", "golang_env", golang_env)

	if err := configor.New(&configor.Config{
		AutoReload:         true,
		AutoReloadInterval: time.Second,
	}).Load(&config, configPath); err != nil {
		panic(err)
	}

	return config
}
