package main

import (
	"encoding/base64"
	"gocv.io/x/gocv"
	"image"
)

var webcam *gocv.VideoCapture
var img gocv.Mat
var smallImg gocv.Mat

func opencvInit() {
	webcam, _ = gocv.VideoCaptureDevice(0)
	img = gocv.NewMat()
	smallImg = gocv.NewMat()
	for i := 0; i < 5; i++ { // 前几帧可能不清晰 故跳过
		webcam.Read(&img)
	}
}

func getImgBase64() string {
	webcam.Read(&img)
	gocv.Resize(img, &smallImg, image.Point{X: 224, Y: 224}, 0, 0, 1)
	encode, err := gocv.IMEncode(".jpg", smallImg)
	if err != nil {
		CRITICAL.Print("Error encoding image")
	}
	imgBase64 := base64.StdEncoding.EncodeToString(encode)
	//fmt.Println(imgBase64)
	//window := gocv.NewWindow("Hello")
	//window.IMShow(smallImg)
	//window.WaitKey(3000)
	return imgBase64
}
