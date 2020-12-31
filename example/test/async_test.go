package example

import (
	"testing"

	"github.com/big-uncle/zap_mate"
)

func TestCheckFile(t *testing.T) {

	logger := zap_mate.NewZapMateLogger("../test.yaml", "default")

	logger.SetAsyncer(10)

	logger.AsyncDebug("Hi, boy!")

	logger.AsyncInfo("I am zap_mate!")

	logger.Flush()

	sugar := logger.Sugar()

	sugar.Error("I am Sugar!")

}
