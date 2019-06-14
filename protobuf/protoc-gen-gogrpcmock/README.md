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
- [ ] enum types 
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

```proto
message Example {
  // Sets the mock response to a static string
  string name = 1 [(grpcmock.field) = { string: "John Smith" }];
  // Randomly generates one of the static strings
  string name = 2 [(grpcmock.field) = { string: "in progress", string: "complete" }]
}
```
