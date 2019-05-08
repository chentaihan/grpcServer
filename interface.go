package grpcServer

import (
	"reflect"
	"runtime"
	"strings"
)

type InterfaceInfo struct {
	Name    string
	Methods []*MethodInfo
}

type MethodInfo struct {
	Object     interface{}	 //对象
	ReqParam   interface{}   //请求参数
	Method     reflect.Value //处理请求对应的方法
	MethodName string        //方法名称
	Url        string        //方法对应url，格式：/类名/方法名
}

func NewInterfaceInfo(in interface{}) *InterfaceInfo {
	itf := &InterfaceInfo{}
	itf.Parse(in)
	return itf
}

/**
解析对象的所有grpc接口方法
 */
func (itf *InterfaceInfo) Parse(in interface{}) bool {
	inType := reflect.TypeOf(in)
	itf.Name = inType.Name()
	if inType.Kind() == reflect.Ptr {
		itf.Name = inType.Elem().Name()
	}

	num := inType.NumMethod()
	for i := 0; i < num; i++ {
		method := inType.Method(i)
		//必须是public方法
		if method.Name[0] < 'A' || method.Name[0] > 'Z' {
			continue
		}

		methodType := method.Type
		//共两个参数
		paramNum := methodType.NumIn()
		if paramNum != 3 {
			continue
		}

		//第1个参数必须是context(第0个参数是对象本身)
		param1 := methodType.In(1)
		if param1.Kind() == reflect.Ptr {
			param1 = param1.Elem()
		}
		if param1.Name() != "Context" {
			continue
		}

		param2 := methodType.In(2)
		if param2.Kind() == reflect.Ptr {
			param2 = param2.Elem()
		}
		mInfo := &MethodInfo{
			Object:   in,
			ReqParam: reflect.New(param2).Interface(),
			Method:   method.Func,
		}
		mInfo.MethodName = itf.getMethodName(mInfo.Method)
		mInfo.Url = itf.createUrl(itf.Name, mInfo.MethodName)
		itf.Methods = append(itf.Methods, mInfo)
	}

	return itf.IsServer()
}

//生成接口url[/类名/方法名]
func (itf *InterfaceInfo) createUrl(serviceName string, methodName string) string {
	return "/" + serviceName + "/" + methodName
}

//获取方法名称
func (itf *InterfaceInfo) getMethodName(method reflect.Value) string {
	methodName := runtime.FuncForPC(method.Pointer()).Name()
	fields := strings.FieldsFunc(methodName, func(sep rune) bool {
		for _, s := range []rune{'/', '.', '-'} {
			if sep == s {
				return true
			}
		}
		return false
	})
	index := len(fields) - 1
	//类方法比函数后面多了一个"-fm"
	if strings.Index(methodName, "-") > 0 {
		index = len(fields) - 2
	}
	return fields[index]
}

func (itf *InterfaceInfo) IsServer() bool {
	if itf.Name == "" || len(itf.Methods) == 0 {
		return false
	}
	return true
}
