package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"go.uber.org/zap"

	_ "net/http/pprof"

	"github.com/big-uncle/zap_mate"
)

//func main() {
//
//	start := time.Now()
//	logger := zap_mate.NewZapMateLogger("./example/test.yaml", "default").SetAsyncer(10)
//
//	logger.AsyncDebug("Hi, boy!", zap.String("ww", "ss"), zap.String("ww", "ss"))
//
//	//logger.Info("I am zap_mate!")
//	log.Printf("耗时【%v】", time.Since(start))
//	logger.Flush()
//
//	time.Sleep(6 * time.Second)
//	log.Printf("耗时【%v】", time.Since(start))
//	//sugar := logger.Sugar()
//	//sugar.Info("")
//	//sugar.Infof("")
//	//sugar.Infow("")
//	//log.Printf("耗时【%v】", time.Since(start))
//
//}

//var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var logger *zap_mate.ZapMateLogger

func main() {
	//flag.Parse()
	//if *cpuprofile != "" {
	//	f, err := os.Create(*cpuprofile)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//	pprof.StartCPUProfile(f)
	//	defer pprof.StopCPUProfile()
	//}
	go func() {
		logger = zap_mate.NewZapMateLogger("./example/test.yaml", "default").SetAsyncer(100000)
		num := 0
		start := time.Now()
		for num < 100000 {
			//start := time.Now()
			//logger := zap_mate.NewZapMateLogger("./test.yaml", "default").SetAsyncer(10)
			//设置异步//传递 loginfo 指针
			//logger.AsyncDebug("Hi, boy!", zap.String("NAME", "JERO"), zap.String("addr", "xian")) //  NumGC = 7

			//logger.AsyncDebug("Hi, boy!") //NumGC = 4
			//传递 loginfo结构体
			//logger.AsyncDebug("Hi, boy!") NumGC = 3
			//log.Printf("%p", logger)
			//log.Printf("%p", logger.asynclogger)
			//sk := logger.Named("ss")

			//log.Printf("%p", sk)
			//log.Printf("%p", sk.asynclogger)

			//logger.AsyncInfo("source switch")
			//sk.AsyncDebug("Hi, boy!", zap.String("NAME", "JERO"), zap.Int("addr", num))     //10W  NumGC = 7
			//sk.AsyncInfo("Hi, boy!", zap.String("NAME", "JERO"), zap.Int("addr", num)) //10W  NumGC = 7
			//sk.Info("w", zap.String("NAME", "JERO"), zap.Int("addr", num))
			//logger.Info("aa", zap.String("NAME", "JERO"), zap.Int("addr", num))
			//使用了sync.pool传递指针 的确少了不少    只进行了 4次GC      100W =16 gc
			//b := logger.Sugar()
			//sk.Info("222222")
			//b.Info("111", 222, "哎")
			//b.Infow("www", "name", "age")
			//b.Infof("我爱你[%v]", "窝窝")
			//不设置异步  俩者擦不多一样
			//logger.AsyncDebug("Hi, boy!") //NumGC =2
			//if num < 30 {
			//	//runtime.GC()
			//}

			//logger.Debug("Hi, boy!") //NumGC =2

			//logger.Info("I am zap_mate!")
			//logger.Flush()

			//log.Printf("耗时【%v】", time.Since(start))
			//time.Sleep(100 * time.Millisecond)
			num++
		}

		//log.Println(len(zap_mate.Point))

		//for
		//logger.Flush()
		log.Printf("耗时【%v】", time.Since(start))

	}()
	//go tool pprof --alloc_space http://localhost:18899/debug/pprof/heap
	//go tool pprof -inuse_space http://localhost:18899/debug/pprof/heap

	//go tool pprof -inuse_space -cum -svg http://localhost:18899/debug/pprof/heap > ./Desktop/heap_inuse3.svg
	//go func() {
	http.HandleFunc("/log", Index)
	http.ListenAndServe("0.0.0.0:18899", nil)

	//}()
}

//ServeHTTP(ResponseWriter, *Request)
func Index(w http.ResponseWriter, r *http.Request) {
	wrilog()
	fmt.Fprintf(w, "success!")
}
func wrilog() {
	num := 0
	start := time.Now()
	//su := logger.Sugar()

	for num < 1 {
		//logger.AsyncInfo("source switch")
		//logger.asynclogger.Info("source switch")

		//su.Infof("source switch{%s}", "hahah")
		//su.Debugw("error: {%s}", "a", "b")
		//su.Debugw("error: {%s}", "a", "b")
		//su.AsyncInfof("111[%s]", "ss")
		//su.AsyncInfof("111[%s]", "ss111")
		////logger.Debug("1121")
		//logger.Logger.Sugar().Named()
		//logger.AsyncInfo()
		//logger.Logger.Named()
		//logger.Named()
		//su.Named()

		orilog := logger
		orilog = orilog.With(zap.String("nb", "plus"))
		//orilog.With("SSS",)
		orilog.Info("SSSSSSSSSSS")
		a := orilog.Named("AA").Named("BB").Named("CC").Named("DD")
		orisu := a.Sugar()

		a.Info("SS")
		orilog.Info("SSSSSSSSSSS")
		orisu.Info("111")
		orisu.With()
		ns := orisu.Named("----")
		ns.Info("wwww")
		oril := ns.Desugar()

		PP := oril.Named("KKKK")
		LL := ns.Named("lll")
		ns.Info("wwwwsss")
		PP.Info("000")
		LL.Info("mmm")
		oril.Info("11111")
		//2020-12-31 15:19:07.746 INFO    SSSSSSSSSSS     {"nb": "plus"}
		//2020-12-31 15:19:07.746 INFO    AA.BB.CC.DD     SS      {"nb": "plus"}
		//2020-12-31 15:19:07.746 INFO    SSSSSSSSSSS     {"nb": "plus"}
		//2020-12-31 15:19:07.746 INFO    AA.BB.CC.DD     111     {"nb": "plus"}
		//2020-12-31 15:19:07.746 INFO    AA.BB.CC.DD.----        wwww    {"nb": "plus"}
		//2020-12-31 15:19:07.746 INFO    AA.BB.CC.DD.----        wwwwsss {"nb": "plus"}
		//2020-12-31 15:19:07.746 INFO    AA.BB.CC.DD.----.KKKK   000     {"nb": "plus"}
		//2020-12-31 15:19:07.746 INFO    AA.BB.CC.DD.----.lll    mmm     {"nb": "plus"}
		//2020-12-31 15:19:07.746 INFO    AA.BB.CC.DD.----        11111   {"nb": "plus"}

		//a.AsyncInfo("SS")
		//
		//orisu.AsyncInfo("111")
		//
		//ns := orisu.Named("----")
		//ns.AsyncInfo("wwww")
		//oril := ns.Desugar()
		//ns.Named("lll")
		//ns.AsyncInfo("wwwwsss")
		//orisu.AsyncInfo("000")
		//oril.AsyncInfo("mmm")
		//ns.Info("mmm")
		//
		//logger111 := zap_mate.NewZapLogger("./example/test.yaml", "default")
		//
		//l2 := logger111.With(zap.String("www", "aaa"))
		//l2.Info("21")
		num++
	}
	log.Printf("耗时【%v】", time.Since(start))

}
