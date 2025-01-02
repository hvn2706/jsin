package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"

	"jsin/logger"
	"jsin/pkg/common"
	"jsin/pkg/constants"
)

const (
	tmpConfigFileName    = "config.tmp.yml"
	actualConfigFileName = "config.yml"
)

// GlobalCfg global variable to access app configuration without passing it around
var GlobalCfg Config

type Database struct {
	MySQLConfig     MySQLConfig `mapstructure:"mysql"`
	MySQLTestConfig MySQLConfig `mapstructure:"mysql_test"`
}

type MySQLConfig struct {
	Host string `mapstructure:"db_host"`
	Port string `mapstructure:"db_port"`
	User string `mapstructure:"username"`
	Pass string `mapstructure:"password"`
	Name string `mapstructure:"db_name"`

	MaxOpenCons        int   `mapstructure:"max_open_cons"`
	MaxIdleCons        int   `mapstructure:"max_idle_cons"`
	ConnMaxIdleTimeSec int64 `mapstructure:"conn_max_idle_time_sec"`
	ConnMaxLifetimeSec int64 `mapstructure:"conn_max_life_time_sec"`
}

// Config for app configuration
type Config struct {
	Logger          logger.LoggerConfig `mapstructure:"logger"`
	Database        Database            `mapstructure:"database"`
	ExternalService External            `mapstructure:"external"`
	TelegramBot     TelegramBot         `mapstructure:"telegram_bot"`
	HelpContent     string              `mapstructure:"help_command_content"`
}

type TelegramBot struct {
	Token               string `mapstructure:"token"`
	Debug               bool   `mapstructure:"debug"`
	Offset              int    `mapstructure:"offset"`
	Timeout             int    `mapstructure:"timeout"`
	CreatCronJobContent string `mapstructure:"create_cronjob_command_content"`
	CronJobImageCaption string `mapstructure:"cronjob_image_caption"`
}

type External struct {
	S3 S3Config `mapstructure:"s3"`
}

type S3Config struct {
	Cloudflare S3CloudflareConfig `mapstructure:"cloudflare"`
}

type S3CloudflareConfig struct {
	Bucket               string `mapstructure:"bucket"`
	Uri                  string `mapstructure:"uri"`
	AccountId            string `mapstructure:"account_id"`
	Token                string `mapstructure:"token"`
	AccessKeyID          string `mapstructure:"access_key_id"`
	SecretAccessKey      string `mapstructure:"secret_access_key"`
	JurisdictionSpecific string `mapstructure:"jurisdiction_specific"`
}

func Load() {
	if !common.CheckIfFileExist(actualConfigFileName) {
		err := os.Rename(tmpConfigFileName, actualConfigFileName)
		if err != nil {
			fmt.Println("Error renaming file:", err)
		} else {
			fmt.Println("File renamed successfully.")
			defer func() {
				err := os.Rename(actualConfigFileName, tmpConfigFileName)
				if err != nil {
					logger.Errorf("Error renaming file: %v", err)
				}
			}()
		}
	}

	vip := viper.New()
	vip.SetConfigName(constants.ConstConfig)
	vip.SetConfigType(constants.Yml)
	vip.AddConfigPath(constants.RootPath) // ROOT

	vip.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vip.AutomaticEnv()

	err := vip.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// Workaround https://github.com/spf13/viper/issues/188#issuecomment-399518663
	// to allow read from environment variables when Unmarshal
	for _, key := range vip.AllKeys() {
		var (
			js     interface{}
			val    = vip.Get(key)
			valStr = fmt.Sprintf("%v", val)
		)

		err := json.Unmarshal([]byte(valStr), &js)

		if err != nil {
			vip.Set(key, val)
		} else {
			vip.Set(key, js)
		}
	}

	fmt.Printf("===== Config file used: %+v \n", vip.ConfigFileUsed())

	GlobalCfg = Config{}
	err = vip.Unmarshal(&GlobalCfg)
	if err != nil {
		panic(err)
	}
}
