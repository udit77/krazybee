package jsonconfig

import (
	"encoding/json"
	"errors"
	"github.com/krazybee/internals/config"
	"github.com/krazybee/internals/config/model"
	"log"
	"os"
)


const (
	module string = "searchservice"
)

type Provider struct{
	path string
	config *model.Config
}

func NewJsonConfigProvider(path string) (config.Provider, error) {
	filePath := path+"/"+module+".json"
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Printf("[JsonConfig] error in reading json config %v\n",err)
		return nil, errors.New("config provider expects valid path to json file")
	}
	return &Provider{path: path, config: &model.Config{}}, nil
}

func readConfig(cfg *model.Config, path string, module string) error{
	filename := path + "/" + module + ".json"
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		return err
	}
	return nil
}


func (i *Provider) Initialize() error {
	err := readConfig(i.config, i.path, module)
	if err != nil {
		return errors.New("error reading service config")
	}
	return nil
}

func (i *Provider) GetConfig() *model.Config {
	return i.config
}
