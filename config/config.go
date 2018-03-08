package config

import (
	"errors"
	"math/big"
	"os"
	"path/filepath"
	"seth/log"

	"github.com/naoina/toml"
)

// SethConfigFile seth config file name
const SethConfigFile = "seth.conf"

// Config var for config
var Config _Config

// _Config config for App
type _Config struct {
	Name    string   `toml:"name"`
	DataDir string   `toml:"datadir"`
	ChainID *big.Int `toml:"chainid"`
}

func init() {
	workPath, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	configFilePath := filepath.Join(workPath, "conf", SethConfigFile)
	LoadConfig(configFilePath)
}

// LoadConfig load config file
func LoadConfig(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		log.Warn("can't load seth.conf file")
		return err
	}
	defer file.Close()

	err = toml.NewDecoder(file).Decode(&Config)

	if _, ok := err.(*toml.LineError); ok {
		err = errors.New(filepath + ", " + err.Error())
	}
	log.Info("DataDir: %s", Config.DataDir)
	return err

}

// ResolvePath return right path
func ResolvePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	if Config.DataDir == "" {
		return ""
	}
	fullpath := filepath.Join(Config.DataDir, Config.Name)
	fullpath = filepath.Join(fullpath, path)
	return fullpath
}
