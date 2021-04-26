package main

import (
	golcm "github.com/lcm-proj/lcm/lcm-go/lcm"
	"log"
	"os"
	"time"
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
	StepHeight:  0.1,
	VDes:        [3]float32{0, 0, 0},
}

func main() {
	MQTTInit()
	opencvInit()
	subscribeMQTT("control")

	lcm, err := golcm.NewProvider("udpm://239.255.76.67:7667?ttl=255")
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
	time.Sleep(5 * time.Second)

	INFO.Print("Start")
	robotLcm.ControlMode = 11
	setRobotMode(*robotLcm)
	time.Sleep(1 * time.Second)

	for true {
		publishMQTT("control", putJson(dataImg{Image: getImgBase64()}))
		mqttPayload := <-mqttChannel
		INFO.Print("Received MQTT Message:")
		INFO.Printf("\ttopic: %s\n", mqttPayload[0])
		INFO.Printf("\tmessage: %s\n", mqttPayload[1])
		// TODO: get vdes from message
		robotLcm.VDes = [3]float32{0, 0, 0}
		lcmPayload, err := robotLcm.Encode()
		if err != nil {
			CRITICAL.Print(err)
		}
		INFO.Printf("Publish LCM Message:")
		INFO.Printf("\tVDes:\t%f\t%f\t%f", robotLcm.VDes[0], robotLcm.VDes[1], robotLcm.VDes[2])
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
