package config

import (
	"encoding/json"
	"fmt"
	"jsin/pkg/constants"
	"os"
	"strings"

	"github.com/spf13/viper"
	"jsin/logger"
	"jsin/pkg/common"
)

const (
	tmpConfigFileName    = "config.tmp.yml"
	actualConfigFileName = "config.yml"
)

// GlobalCfg global variable to access app configuration without passing it around
var GlobalCfg Config

// ===== Init structs =====

// ServerListen for specifying host & port
type ServerListen struct {
	Host string `mapstructure:"host"`
	Port uint16 `mapstructure:"port"`
}

// ServerConfig for configure HTTP & gRPC host & port
type ServerConfig struct {
	HTTP   ServerListen `mapstructure:"http"`
	ApiKey string       `mapstructure:"api_key"`
}

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
	Server   ServerConfig        `mapstructure:"server"`
	Logger   logger.LoggerConfig `mapstructure:"logger"`
	Database Database            `mapstructure:"database"`
}

// ===== Util func =====

func (s ServerListen) String() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

// ListenString for listen to 0.0.0.0
func (s ServerListen) ListenString() string {
	return fmt.Sprintf(":%d", s.Port)
}

func Load() {
	os.Chdir("../..")

	if !common.CheckIfFileExist(actualConfigFileName) {
		err := os.Rename(tmpConfigFileName, actualConfigFileName)
		if err != nil {
			fmt.Println("Error renaming file:", err)
		} else {
			fmt.Println("File renamed successfully.")
			defer os.Rename(actualConfigFileName, tmpConfigFileName)
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
	return
}
