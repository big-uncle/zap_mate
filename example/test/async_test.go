package example

import (
	"github.com/big-uncle/zap_mate"
	"testing"
)

func TestCheckFile(t *testing.T) {

	logger := zap_mate.NewZapMateLogger("../test.yaml", "default")

	logger.SetAsyncer(10)

	logger.AsyncDebug("Hi, boy!")

	logger.AsyncInfo("I am zap_mate!")

	logger.Flush()

	logger.Warn("Who are you?")

	sugar := logger.Sugar()

	sugar.Error("I am Sugar!")

	sugar.DPanic("How are you?")

}
