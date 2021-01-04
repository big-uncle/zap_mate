package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"runtime/pprof"
	"time"

	_ "net/http/pprof"

	"github.com/big-uncle/zap_mate"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

var logger *zap_mate.ZapMateLogger

func main() {

	flag.Parse()

	if *cpuprofile != "" {

		f, err := os.Create(*cpuprofile)

		if err != nil {

			return

		}

		pprof.StartCPUProfile(f)

		defer pprof.StopCPUProfile()

	}

	logger = zap_mate.NewZapMateLogger("./example/test.yaml", "default")

	logger.SetAsyncer(100000)

	http.HandleFunc("/", log)

	http.ListenAndServe("0.0.0.0:8899", nil)

}

//go tool pprof --alloc_space http://localhost:8899/debug/pprof/heap

//go tool pprof -inuse_space http://localhost:8899/debug/pprof/heap

//go tool pprof -inuse_space -cum -svg http://localhost:8899/debug/pprof/heap > ./Desktop/heap_inuse3.svg

func log(w http.ResponseWriter, r *http.Request) {

	sugar := logger.Sugar()

	num := 0

	start := time.Now()

	for num < 1000 {

		sugar.AsyncInfof("Hi , boy!")

		//logger.Info("Hi , boy!")

		num++

	}

	logger.Flush()

	fmt.Printf("TS:[%v]\n", time.Since(start))

	fmt.Fprintf(w, "TS:[%v]", time.Since(start))

}
