# gRPC Mock Generator

A `protoc` plugin that generates a mock implementation of a gRPC service.

## Features

- [x] auto-generated mock response messages
- [x] unary rpc support
- [ ] server streaming support
- [x] string type
- [x] int types
- [ ] float types
- [x] inner messages
- [ ] bytes type
- [ ] oneof type
- [ ] map type
- [ ] option overrides*

## Usage

```bash
protoc -I. --gogo_out=plugin=grpc=:/out --gogrpcmock_out=:/out src/*.proto
```

```go
grpcServer := grpc.NewServer()
myservice.RegisterMyServiceServer(s, &myservice.MyServiceMock{})
lis, _ := net.Listen("tcp", 50501)
grpcServer.Serve(lis)
```

## Options

> WIP
