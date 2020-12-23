package main

import (
	"github.com/big-uncle/zap_mate"
)

func main() {

	zaplog := zap_mate.NewLogger("./test.yml", "default")

	zaplog.Info("Hi, body!")

	zaplog.Warn("I am zap_mate!")

}
