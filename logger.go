package client

// Logger logger
type Logger interface {
	// Warningf 警告
	Warningf(format string, args ...interface{})
	// Warningln 警告
	Warningln(args ...interface{})
	// Errorf 错误
	Errorf(format string, args ...interface{})
	// Errorln 错误
	Errorln(args ...interface{})
	// Infof 信息
	Infof(format string, args ...interface{})
	// Infoln 消息
	Infoln(args ...interface{})
	// Printf 打印
	Printf(format string, args ...interface{})
	// Println 打印
	Println(args ...interface{})
	// Debugf 测试
	Debugf(format string, args ...interface{})
	// Debugln 测试
	Debugln(args ...interface{})
}
