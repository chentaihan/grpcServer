# grpcServer

## grpc管理器

### 目的

`简化代码：原生grpc方法，每个方法都要写一堆的代码。使用grpcServer，
你只需要关注grpc方法的实现，其他的代码都不需要写`

### 原理
`通过反射实现自动注册grpc方法`