package config

import (
	"log"
	"os"

	"k8s.io/apimachinery/pkg/util/yaml"
)

type Config struct {
	Server struct {
		Port int `yaml:"port"`
	}
	K8s struct {
		// k8s api server endpoint
		Host string `yaml:"host"`
		// k8s access token
		Token string `yaml:"token"`
	} `yaml:"k8s"`
}

var defaultConfig Config

func init() {
	fileData, err := os.ReadFile("./config/config.yaml")
	if err != nil {
		log.Panic("Couldn't read config file", err)
	}

	err = yaml.Unmarshal(fileData, &defaultConfig)
	if err != nil {
		log.Panic("Couldn't unmarshal config file", err)
	}
}

func GetConfig() Config {
	return defaultConfig
}

type ApiResult struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

const SUCCESS_CODE = 2000
const ERROR_CODE = 3000

func MakeSucessResult(data interface{}) *ApiResult {
	return &ApiResult{Code: SUCCESS_CODE, Data: data}
}

func MakeErrorResult(err string) *ApiResult {
	return &ApiResult{Code: ERROR_CODE, Data: nil, Error: err}
}
