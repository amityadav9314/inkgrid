package config

var environment string

func SetEnv(env string) {
	environment = env
}

func GetEnv() string {
	return environment
}
