go-logger 是golang 的日志库 ，基于对golang内置log的封装。
用法类似java日志工具包log4j

打印日志有6个方法 Debug，Info，Warn，Error，Fatal，Key日志级别由低到高，对应的支持
Debugf，Infof，Warnf，Errorf，Fatalf，Keyf格式化的方式 

日志特性：
1.参考：https://github.com/donnie4w/go-logger
2.支持日志滚动，且老日志在末文件位置
3.封装了一层，可直接使用

示例：
    初始化：logger.Init("logger", logger.DEBUG) // 指定日志模块名+日志等级

打印日志：
func log(i int) {
    logger.Debug("Debug>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
    logger.Info("Info>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
    logger.Warn("Warn>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
    logger.Error("Error>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
    logger.Fatal("Fatal>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
    logger.Key("Key>>>>>>>>>>>>>>>>>>>>>>>>>" + strconv.Itoa(i))
    logger.Keyf("Key>>>>>>>>>>>>>>>>>>>>>>>>>%v",i)
}