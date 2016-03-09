package main

import (
	"strconv"
	"time"

	"github.com/arbboter/go-logger/logger"
)

func main() {
	TestMain()

}

func log(i int) {
	logger.Debug("Debug>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
	logger.Info("Info>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
	logger.Warn("Warn>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
	logger.Error("Error>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
	logger.Fatal("Fatal>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
	logger.Key("Key>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
	logger.Keyf("Key>>>>>>>>>>>>>>>>>>>>>>>>>%v", i)
}

func TestMain() {
	logger.Init("logger", logger.DEBUG)
	nNum := 100000000
	for nNum > 0 {
		for i := 10; i > 0; i-- {
			go log(i)
		}
		nNum -= 1
		time.Sleep(30 * time.Millisecond)
	}
	time.Sleep(15 * time.Second)
}
