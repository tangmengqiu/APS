package src

import (
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type gc struct {
	DingUrl string
}

var GlobalConfig gc

// Config define
type Config struct {
	Name string
}

// Init prepare config
func InitConfig(cfg string) error {
	c := Config{
		Name: cfg,
	}

	if err := c.initConfig(); err != nil {
		return err
	}
	//read config
	GlobalConfig.DingUrl = viper.GetString("ding.url")

	c.watchConfig()

	return nil
}

func (c *Config) initConfig() error {
	if c.Name != "" {
		viper.SetConfigFile(c.Name)
	} else {
		viper.AddConfigPath("conf")
		viper.SetConfigName("conf")
	}
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()            // Read Env value
	viper.SetEnvPrefix("APISERVER") // Env prefix : APISERVER
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig(); err != nil { // viper read config
		return err
	}

	return nil
}

// watch config change
func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logrus.Info("Config file changed: %s", e.Name)
	})
}
