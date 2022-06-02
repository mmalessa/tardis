package app

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func InitConfig(environment string) error {
	var configDir string

	if err := loadDotEnv(); err != nil {
		return err
	}

	config.WithOptions(config.ParseEnv)
	config.AddDriver(yaml.Driver)

	switch environment {
	case "dev":
		configDir = filepath.Join("/go/src/", AppName, "config")
	default:
		configDir = filepath.Join("/etc/", AppName)
	}
	log.Infof("Set configdir =\"%s\"", configDir)

	// load default configuration
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		errorMsg := fmt.Sprintf("Config directory \"%s\" does not exists", configDir)
		log.Fatal(errorMsg)
		return errors.New(errorMsg)
	}
	if err := loadConfigFilesFromDirectory(configDir); err != nil {
		return err
	}

	// // load configuration for environment
	// envConfigDir := filepath.Join(configDir, environment)
	// if _, err := os.Stat(envConfigDir); os.IsNotExist(err) {
	// 	errorMsg := fmt.Sprintf("Config directory for environment \"%s\" does not exists", environment)
	// 	log.Warning(errorMsg)
	// 	return nil
	// }
	// if err := loadConfigFilesFromDirectory(envConfigDir); err != nil {
	// 	return err
	// }
	return nil
}

func loadConfigFilesFromDirectory(directory string) error {
	log.Infof("Load config liles from directory \"%s\"", directory)
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return err
	}
	for _, file := range files {
		fileName := file.Name()
		if strings.HasSuffix(fileName, ".yaml") {
			fullPathFileName := filepath.Join(directory, fileName)
			if err := config.LoadFiles(fullPathFileName); err != nil {
				return err
			}
			log.Infof("Load config file: \"%s\"", fullPathFileName)
		}
	}
	return nil

}

//FIXME - do it better
func loadDotEnv() error {

	fromFiles := false

	if _, err := os.Stat(".env.dev"); err == nil {
		log.Info("ENV from .env.dev")
		godotenv.Load(".env.dev")
		fromFiles = true
	}

	if _, err := os.Stat(".env.prod"); err == nil {
		log.Info("ENV from .env.prod")
		godotenv.Load(".env.prod")
		fromFiles = true
	}

	if _, err := os.Stat(".env"); err == nil {
		log.Info("ENV from .env")
		godotenv.Load(".env")
		fromFiles = true
	}

	if !fromFiles {
		log.Info(".env doesn't exist. Let us have faith that there are appropriate environmental variables.")
	}

	return nil
}
