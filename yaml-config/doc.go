package YAMLConfig

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	. "github.com/azer/on-change"
	. "github.com/azer/debug"
)

type (
	Document map[string]map[string]string
	RawDocument map[string]interface{}
)

type YAMLConfig struct {
	Filename string
	Document Document
	callback func(map[string]map[string]string)
}

func (config *YAMLConfig) Load () {
	raw, err := Read(config.Filename)

	if err != nil {
		Debug("Failed to read and parse %s", config.Filename)
		return
	}

	config.Document = Normalize(raw)

	go config.callback(config.Document)
}

func (config *YAMLConfig) EnableAutoReload () {
	OnChange(config.Filename, config.Load)
}

func NewYAMLConfig (filename string, callback func(map[string]map[string]string)) *YAMLConfig {
	Debug("Creating a new YAML config from %s", filename)

	config := &YAMLConfig{filename, nil, callback}
	config.Load()

	return config
}

func Normalize (raw RawDocument) Document {
	Debug("Loading...")

	config := make(Document)

	for hostname, options := range raw {
		switch t := options.(type) {
		case string:
			config[hostname] = make(map[string]string)
			config[hostname]["/"] = t
		case map[interface{}]interface{}:
			config[hostname] = make(map[string]string)

			for key, value := range t {
				path, ok1 := key.(string)
				uri, ok2 := value.(string)
                if ok1 && ok2 {
					config[hostname][path] = uri
				}
			}
		}
	}

	return config
}

func Read (filename string) (RawDocument, error) {
	Debug("Reading %s", filename)

	raw := make(RawDocument)
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(content, &raw)

	if err != nil {
		return nil, err
	}

	return raw, nil
}
