# zap_mate
The `zap_mate` is a simple packaging tools of [zap](https://github.com/uber-go/zap). It combines zap, viper, lumberjack and many other tools to integrate it, making zap configuration lighter, more convenient and faster, so as to better assist developers.



See the example folder for more details ...


Because asyn logs will cause a lot of memory escape, zap_mate does not recommend or support the use of asyn logs.

If you have a scene with very high real-time performance, then it is recommended that you implement the asyn log yourself. In fact, it is very easy.

- Example
    ```go
type mylogger struct{
    field string
    msg  string
    *zap.Logger
}

    var logger=make(chan *zmylogger,1000)
    go func (){
        for {
            select zl<-logger:
                zl.Info(zl.msg,zl.field)
        }
    }()
    ```
