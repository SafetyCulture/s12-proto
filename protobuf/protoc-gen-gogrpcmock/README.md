# gRPC Mock Generator

A `protoc` plugin that generates a mock implementation of a gRPC service.

## Features

- [x] auto-generated mock response messages
- [x] unary rpc support
- [ ] server streaming support
- [x] string type
- [x] int types
- [x] float types
- [x] inner messages
- [x] enum types
- [x] bool types
- [ ] bytes type
- [x] oneof type
- [ ] map type
- [x] option overrides

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

## Options & Example

```proto
message Example {
  string id = 1;
  // Randomly generates one of the static strings
  string status = 2 [ (grpcmock.field) = {string : "in progress", string : "complete"} ];
  string description = 3;
  string email = 4;
  // Randomly generates one word; this is the default for a string type
  string not_email = 5 [ (grpcmock.field) = {word : true} ];
  string phone = 6;
  repeated string word = 7;
  string url = 8;
  int32 single_number = 9;
  repeated int64 repeated_number = 10;
  int32 lat = 11;
  int32 lng = 12;
  // Randomly generates between 1 and 5 words
  string words = 13 [ (grpcmock.field) = {words : true} ];
  // Randomly generates a set number of words
  string wordsn = 14 [ (grpcmock.field) = {wordsn : 10} ];
  // Randomly generates an integer between 0 and n
  int32 intn = 15 [ (grpcmock.field) = {intn : 5} ];
  // Generates a random full name
  string fullname = 16 [ (grpcmock.field) = { fullname : true } ];
  // Generates a random first name
  string firstname = 17 [ (grpcmock.field) = { firstname : true } ];
  // Generates a random last name
  string lastname = 18 [ (grpcmock.field) = { lastname : true } ];
  // Generates a single paragraph of random words
  string paragraph = 19 [ (grpcmock.field) = {paragraph : true} ];
  // Between 1 and 5 paragraphs of random words
  string paragraphs = 20 [ (grpcmock.field) = {paragraphs : true} ];
  // Set number of paragraphs
  string paragraphsn = 21 [ (grpcmock.field) = {paragraphsn : 2} ];
  string uuid = 22 [ (grpcmock.field) = {uuid : true} ];
  string email_address = 23 [ (grpcmock.field) = {email : true} ];
  string phone_number = 24 [ (grpcmock.field) = {phone : true} ];
  string company = 25 [ (grpcmock.field) = {company : true} ];
  string brand = 26 [ (grpcmock.field) = {brand : true} ];
  string product = 27 [ (grpcmock.field) = {product : true} ];
  string color = 28 [ (grpcmock.field) = {color : true} ];
  // You can also add a prefix to the generated output
  string hexcolor = 29 [ (grpcmock.field) = {prefix : "#", hexcolor : true} ];
  double latitude = 30;
  double longitude = 31;
  float floatn = 32 [ (grpcmock.field) = {floatn : 3} ];
  bool boolean = 33;
}
```
