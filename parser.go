package messenger

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

//Config is the struct for the configuration file
type Config struct {
	VerifyToken string `yaml:"verify_token"`
	AccessToken string `yaml:"access_token"`
	AppSecret   string `yaml:"app_secret"`
}

//ReadYml parses the config yml file and format into the Config struct
func (x *Config) ReadYml() *Config {
	configFile, err := filepath.Abs("./bot.config.yml")
	if err != nil {
		log.Printf("ERROR READING THE CONFIG FILE: %s", err)
	}

	yamlFile, err := ioutil.ReadFile(configFile)

	if err != nil {
		log.Println("Could not find the config file. Please make sure it is created", err)
		os.Exit(-1)
	}

	yaml.Unmarshal(yamlFile, &x)

	return x
}

//GetTokens returns the verifytoken, accesstoken and the appsecret from the config file
func GetTokens() (string, string, string) {
	var c Config

	configObj := c.ReadYml()
	verifyToken, accessToken, appSecret := configObj.VerifyToken, configObj.AccessToken, configObj.AppSecret
	return verifyToken, accessToken, appSecret
}
