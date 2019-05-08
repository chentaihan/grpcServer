package grpcServer

import "fmt"

type ILog interface {
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type LogDefault struct {
}

func (*LogDefault) Info(args ...interface{}) {
	fmt.Println(args...)
}

func (*LogDefault) Warn(args ...interface{}) {
	fmt.Println(args...)
}

func (*LogDefault) Error(args ...interface{}) {
	fmt.Println(args...)
}

func (*LogDefault) Infof(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

func (*LogDefault) Warnf(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}

func (*LogDefault) Errorf(format string, args ...interface{}) {
	fmt.Println(fmt.Sprintf(format, args...))
}
