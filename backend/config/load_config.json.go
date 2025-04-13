// Package config loads all the config in accordance to running environment
// Configuration varies for dev, pp, prod and prod environment
// Viper is used to load all the configs
package config

import (
	"fmt"
	"github.com/amityadav9314/goinkgrid/constants"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"os"
	"runtime"
)

var Config *viper.Viper
var ReadCredFromSecretManager = true
var EnvSecretsMapping = map[string]string{
	"prodpp": "/dev/" + ",/dev/amityadav/gosample/hashed",
	"prod":   "/dev/" + ",/dev/amityadav/gosample/hashed",
	"dev":    "/dev/" + ",/dev/amityadav/gosample/hashed",
}

type Setter interface {
	setConfig(environment string)
}

// reads config via config/{env}/config.json + base_config.json files
type FileConfig struct {
}

// reads config via aws secret manager
type SecretManager struct {
}

func GetStringMapInt(key string) map[string]int {
	return cast.ToStringMapInt(Config.Get(key))
}

func DoInit(environment string) {
	isProdEnv := environment == "prodpp" || environment == "prod" || environment == "qa"
	var configSetter Setter
	if ReadCredFromSecretManager && isProdEnv {
		configSetter = new(SecretManager)
	} else {
		configSetter = new(FileConfig)
	}
	configSetter.setConfig(environment)
}

func (fm *FileConfig) setConfig(environment string) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("error in fetching config directory")
	}

	configPath := filename + "/../"
	Config = viper.New()

	LoadConfig(configPath, "base_config", "")
	LoadConfig(configPath, "config", environment)
}

// reads config via aws secret manager
func (sm *SecretManager) setConfig(environment string) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("error in fetching config directory")
	}

	configPath := filename + "/../"
	Config = viper.New()

	//LoadConfig(configPath, "base_config", "")
	LoadConfig(configPath, "config", environment)

}

func LoadConfig(configPath, configFile, environment string) {
	SetConfig(configPath, configFile, "json", environment)
	if err := Config.MergeInConfig(); err != nil {
		panic(fmt.Errorf("fatal error in reading config file %s", err))
	}
}

func SetConfig(configPath, configFileName, configFileType, environment string) {
	Config.SetConfigName(configFileName)
	Config.SetConfigType(configFileType)
	Config.AddConfigPath(configPath + environment)
}

func LoadAppEnv() string {
	env := os.Getenv(constants.VarLogFile)
	if len(env) == 0 {
		return constants.EnvProd
	}
	return env
}
