package engine

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"testing"
)

func TestEngine(t *testing.T) {
	// Lê o arquivo de configuração
	configData, err := ioutil.ReadFile("data/input1_output1_config1.json")
	if err != nil {
		log.Fatal(err)
	}
	config := []map[string]interface{}{}
	err = json.Unmarshal(configData, &config)
	if err != nil {
		log.Fatal(err)
	}

	inputData, err := ioutil.ReadFile("data/input1.json")
	if err != nil {
		log.Fatal(err)
	}
	input := make(map[string]interface{})
	err = json.Unmarshal(inputData, &input)
	if err != nil {
		log.Fatal(err)
	}

	output := processData(input, config)

	outputJSON, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(outputJSON))
}
