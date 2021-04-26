package main

import (
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

func main() {
	MQTTInit()
	go subscribeMQTT("control")

	lcm, err := golcm.New()
	if err != nil {
		CRITICAL.Print("Create lcm failed")
		panic(err)
	}
	defer lcm.Destroy()
	lcmPublisher, errs := lcm.Publisher(lcmChannel)
	defer close(lcmPublisher)
	go watchLcmError(errs)
	i := 0
	for true {
		i++
		publishMQTT("control", putJson(dataImg{Image: getImgBase64()}))
		mqttPayload := <-mqttChannel
		INFO.Print("Received MQTT Message:")
		INFO.Printf("\ttopic: %s\n", mqttPayload[0])
		INFO.Printf("\tmessage: %s\n", mqttPayload[1])

		robotLcm := &ExlcmRobotControlLcmt{
			ControlMode: int32(i),
		}
		lcmPayload, err := robotLcm.Encode()
		if err != nil {
			CRITICAL.Print(err)
		}
		INFO.Printf("sending data")
		lcmPublisher <- lcmPayload
	}
}
