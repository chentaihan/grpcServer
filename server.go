package grpcServer

import (
	"context"
	"google.golang.org/grpc"
	"reflect"
	"time"
)

type IRpcServer interface {
}

type RpcServerMgr struct {
	server *grpc.Server
	log    ILog
}

func NewRpcRegister(server *grpc.Server, log ILog) *RpcServerMgr {
	if log == nil {
		log = &LogDefault{}
	}
	return &RpcServerMgr{
		server: server,
		log:    log,
	}
}

func (rpcReg *RpcServerMgr) Register(service IRpcServer) bool {
	interfaceInfo := NewInterfaceInfo(service)
	if !interfaceInfo.IsServer() {
		return false
	}
	handlerMethods := rpcReg.createMethodDesc(interfaceInfo.Methods, interfaceInfo)
	sd := &grpc.ServiceDesc{
		ServiceName: interfaceInfo.Name,
		HandlerType: (*IRpcServer)(nil),
		Methods:     handlerMethods,
		Streams:     []grpc.StreamDesc{},
		Metadata:    "server.proto",
	}
	rpcReg.server.RegisterService(sd, service)
	return true
}

func (rpcReg *RpcServerMgr) requestHandler(req interface{}, method *MethodInfo, service *InterfaceInfo) func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	return func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
		start := time.Now()
		if err := dec(req); err != nil {
			rpcReg.log.Error("requestHandler error, reason:", err.Error())
			return nil, err
		}
		url := method.Url
		if interceptor == nil {
			args := []reflect.Value{reflect.ValueOf(method.Object), reflect.ValueOf(ctx), reflect.ValueOf(req)}
			result := method.Method.Call(args)
			resp := result[0].Interface()
			err := result[1].Interface()
			latency := float64(time.Now().Sub(start))/float64(time.Second)
			if err == nil {
				rpcReg.log.Infof("requestHandler success method=%s, url=%s, req=%v, resp=%v, time=%.5fs", method.MethodName, url, req, resp, latency)
				return resp, nil
			}
			rpcReg.log.Errorf("requestHandler failed method=%s, url=%s, req=%v, resp=%v, time=%.5fs", method.MethodName, url, req, resp, latency)
			return resp, err.(error)
		}

		info := &grpc.UnaryServerInfo{
			Server:     srv,
			FullMethod: url,
		}
		handler := func(ctx context.Context, req interface{}) (interface{}, error) {
			args := []reflect.Value{reflect.ValueOf(method.Object), reflect.ValueOf(ctx), reflect.ValueOf(req)}
			result := method.Method.Call(args)
			resp := result[0].Interface()
			err := result[1].Interface()
			latency := float64(time.Now().Sub(start))/float64(time.Second)
			if err == nil {
				rpcReg.log.Infof("requestHandler success method=%s, url=%s, req=%v, resp=%v, time=%.5fs", method.MethodName, url, req, resp, latency)
				return resp, nil
			}
			rpcReg.log.Errorf("requestHandler failed method=%s, url=%s, req=%v, resp=%v, time=%.5fs", method.MethodName, url, req, resp, latency)
			return resp, err.(error)
		}
		return interceptor(ctx, req, info, handler)
	}
}

func (rpcReg *RpcServerMgr) createMethodDesc(methodInfo []*MethodInfo, service *InterfaceInfo) []grpc.MethodDesc {
	var methodList []grpc.MethodDesc
	for _, mInfo := range methodInfo {
		method := grpc.MethodDesc{
			MethodName: mInfo.MethodName,
			Handler:    rpcReg.requestHandler(mInfo.ReqParam, mInfo, service),
		}
		methodList = append(methodList, method)
	}
	return methodList
}
