package logger

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
)

// Log 获取log对象
func Log(name ...string) *glog.Logger {
	if len(name) > 0 && name[0] != "" {
		return g.Log(name[0]).Skip(1).Line()
	}
	return g.Log().Skip(1).Line()
}

// Print 普通打印，不带任何标签
func Print(v ...interface{}) {
	Log("access").Print(v)
}

// Println 普通打印，不带任何标签
func Println(v ...interface{}) {
	Log("access").Println(v)
}

// Info INFO打印，带[INFO]标签
func Info(v ...interface{}) {
	Log("access").Info(v)
}

// Infof INFO打印，带[INFO]标签
func Infof(format string, v ...interface{}) {
	Log("access").Infof(format, v...)
}

// Debug Debug打印，带[Debug]标签
func Debug(v ...interface{}) {
	Log("access").Debug(v)
}

// Debugf Debug打印，带[Debug]标签
func Debugf(format string, v ...interface{}) {
	Log("access").Debugf(format, v...)
}

// Error Error打印，带[Error]标签
func Error(v ...interface{}) {
	Log("error").Error(v)
}

// Errorf Error打印，带[Error]标签
func Errorf(format string, v ...interface{}) {
	Log("error").Errorf(format, v...)
}
