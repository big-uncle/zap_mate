package main

import (
	"flag"
	"log"
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
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}
	go func() {
		logger = zap_mate.NewZapMateLogger("./example/test.yaml", "default")
		logger.SetAsyncer(100000)
		num := 0
		start := time.Now()
		for num < 100000 {
			logger.AsyncInfo("Hi , boy!")
			num++
		}

		log.Printf("TS:[%v]", time.Since(start))

	}()

	//go tool pprof --alloc_space http://localhost:18899/debug/pprof/heap
	//go tool pprof -inuse_space http://localhost:18899/debug/pprof/heap
	//go tool pprof -inuse_space -cum -svg http://localhost:18899/debug/pprof/heap > ./Desktop/heap_inuse3.svg

	http.ListenAndServe("0.0.0.0:8899", nil)

}
