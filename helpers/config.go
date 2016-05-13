package helpers

import (
	"encoding/json"
	"io/ioutil"
)

func ImportConfig(configPath string) (Configuration, error) {
	f, err := ioutil.ReadFile(configPath)
	if err != nil {
		return Configuration{}, err
	}

	var configurationData Configuration
	json.Unmarshal(f, &configurationData)

	return configurationData, nil
}
