package main

import (
	"github.com/big-uncle/zap_mate"
	"log"
	"time"
)

func main() {

	start := time.Now()
	logger := zap_mate.NewZapMateLogger("./example/test.yaml", "default")
	var num int
	for num < 10 {
		logger.Debug("Hi, boy!")

		logger.Info("I am zap_mate!")

		num++
	}
	logger.Flush()

	log.Printf("耗时【%v】", time.Since(start))

}
