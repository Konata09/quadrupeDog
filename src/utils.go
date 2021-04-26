package main

import (
	"encoding/json"
	"time"
)

const (
	GETCONTROLBYCAM = "getControlByCam"
)

type DogJson struct {
	DogId       int         `json:"dog_id"`
	DogName     string      `json:"dog_name"`
	DogLocation string      `json:"dog_location"`
	Timestamp   int64       `json:"timestamp"`
	Status      bool        `json:"status"`
	StatusMsg   string      `json:"status_msg"`
	Type        string      `json:"type"`
	Data        interface{} `json:"data"`
}

type dataImg struct {
	Image string `json:"image"`
}

func putJson(data interface{}) string {
	dogjson := &DogJson{
		DogId:       dogId,
		DogName:     dogName,
		DogLocation: dogLocation,
		Timestamp:   time.Now().Unix(),
		Status:      true,
		StatusMsg:   "OK",
		Type:        GETCONTROLBYCAM,
		Data:        data,
	}
	marshal, err := json.Marshal(dogjson)
	if err != nil {
		CRITICAL.Print("Error encoding json")
		return ""
	}
	return string(marshal)
}

func watchLcmError(errs <-chan error) {
FOR_SELECT:
	for {
		select {
		case err, ok := <-errs:
			if !ok {
				break FOR_SELECT
			}
			CRITICAL.Print("LCM Publisher Error")
			panic(err)
		}
	}
}
