package main

import (
	"encoding/json"
	golcm "github.com/lcm-proj/lcm/lcm-go/lcm"
	"log"
	"os"
)

const (
	dogName     = "myAwesomeDog"
	dogId       = 1
	dogLocation = "316"
	lcmChannel  = "test"
)

var INFO = log.New(os.Stdout, "[INFO] ", log.LstdFlags)
var CRITICAL = log.New(os.Stderr, "[CRIT] ", log.LstdFlags|log.Lshortfile)

var mqttChannel = make(chan [2]string)
var lcmPublisher chan<- []byte
var lcmErrs <-chan error

var robotLcm = &ExlcmRobotControlLcmt{
	ControlMode: 12,
	StepHeight:  0.0,
	VDes:        [3]float32{0, 0, 0},
}

func main() {
	MQTTInit()
	opencvInit()

	lcm, err := golcm.NewProvider("udpm://239.255.76.67:7667?ttl=1")
	if err != nil {
		CRITICAL.Print("Create lcm failed")
		panic(err)
	}
	defer lcm.Destroy()
	lcmPublisher, lcmErrs = lcm.Publisher(lcmChannel)
	defer close(lcmPublisher)
	go watchLcmError(lcmErrs)

	INFO.Print("Stand Up")
	robotLcm.ControlMode = 12
	setRobotMode(*robotLcm)
	//time.Sleep(5 * time.Second)

	INFO.Print("Start")
	robotLcm.ControlMode = 11
	robotLcm.StepHeight = 0.1
	setRobotMode(*robotLcm)
	//time.Sleep(1 * time.Second)

	go subscribeMQTT("control")

	for true {
		publishMQTT("robot_upload", putJson(dataImg{Image: getImgBase64()}))
		mqttPayload := <-mqttChannel
		INFO.Print("Received MQTT Message:")
		INFO.Printf("\ttopic: %s\n", mqttPayload[0])
		INFO.Printf("\tmessage: %s\n", mqttPayload[1])
		// TODO: get vdes from message
		control := ExlcmRobotControlLcmt{}
		received := RespJson{Data: &control}
		err := json.Unmarshal([]byte(mqttPayload[1]), &received)
		if err != nil {
			CRITICAL.Print(err)
		}
		robotLcm = received.Data.(*ExlcmRobotControlLcmt)
		//robotLcm.VDes = [3]float32{0, 0.3, 0}
		lcmPayload, err := robotLcm.Encode()
		if err != nil {
			CRITICAL.Print(err)
		}
		INFO.Printf("Publish LCM Message:")
		INFO.Printf("\tControlMode: %v", robotLcm.ControlMode)
		INFO.Printf("\tGaitType: %v", robotLcm.GaitType)
		INFO.Printf("\tVDes: %2.2f %2.2f %2.2f", robotLcm.VDes[0], robotLcm.VDes[1], robotLcm.VDes[2])
		INFO.Printf("\tStepHeight: %2.2f", robotLcm.StepHeight)
		lcmPublisher <- lcmPayload
	}
}

func setRobotMode(lcmt ExlcmRobotControlLcmt) {
	lcmPayload, err := lcmt.Encode()
	if err != nil {
		CRITICAL.Print(err)
	}
	lcmPublisher <- lcmPayload
}
