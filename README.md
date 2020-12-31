# zap_mate

***The `zap_mate` is a simple packaging tools of [zap](https://github.com/uber-go/zap). It combines zap, viper, lumberjack and many other tools to integrate it, making zap configuration lighter, more convenient and faster, so as to better assist developers.***

See the example folder for more details ...

## Usage
### Use the origin ZAP logger
- Example
    ```go

    logger := zap_mate.NewZapLogger("../test.yaml", "default")

    logger.Debug("Hi, boy!")

    logger.Info("I am zap_mate!")
    
    ```
### Use the mate logger
If you want to use mate logger:
> The mate logger supports async and sync write logs, and extends all feature of origin zap.
- Example
    ```go

    logger := zap_mate.NewZapMateLogger("../test.yaml", "default")

	logger.SetAsyncer(10)

	logger.AsyncDebug("Hi, boy!")

    logger.AsyncInfo("I am zap_mate!")
  
    logger.Info("oh oh oh")
    
    logger.Flush()

   	logger.Warn("I am a logger of zap_mate")

	sugar := logger.Sugar()

	sugar.Error("I am Sugar!")
    
    ```
  > Note: **`Child node cannot affect parent nodes,but child node all feature of extends parent node!`**