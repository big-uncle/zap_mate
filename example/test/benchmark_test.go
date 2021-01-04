package example

import (
	"testing"

	"github.com/big-uncle/zap_mate"
)

func Benchmark_Checklog(t *testing.B) {

	logger := zap_mate.NewZapMateLogger("../test.yaml", "default")

	logger.SetAsyncer(100000)

	sugar := logger.Sugar()

	num := 0

	for num < 1000 {

		sugar.AsyncInfof("Hi , boy!")

		logger.Info("Hi , boy!")

		logger.AsyncInfo("I am zap_mate!")

		sugar.Error("I am Sugar!")

		num++
	}

	logger.Flush()

}

//go test -v -bench=. benchmark_test.go

//go test -v -bench=. -benchmem
