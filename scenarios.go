package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func LoadParameters(filepath string) {
	log.Printf("Load parameters for file:%v", filepath)
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	var data interface{}
	json.Unmarshal(file, &data)
	msg := data.(map[string]interface{})

	for k, v := range msg {
		SetParameter(k, v.(string))
	}
}
