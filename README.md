# Protocol Buffer tools used by SafetyCulture

**WARNING:** Public Repo

## Pre-requisites:

* [Go](https://golang.org/doc/install)
* [Protocol Buffer Compiler](https://grpc.io/docs/protoc-installation/)

### Generating Go code from protobuf

Make sure you have `protoc-gen-go` installed:
https://pkg.go.dev/google.golang.org/protobuf#section-readme
```
$ go install github.com/golang/protobuf/protoc-gen-go@latest
```

Run the generate command:

```
$ make generate
```

To generate and run tests:
```
make generate && make govalidator && make govalidator-valtest
```

#### Making sure the code was generated as expected

Let's use `protoc-gen-s12perm` package as an example which will generate an
example Go code from a protobuf file.

```
$ make s12perm
```

This command should generate an example file written in Go. Now we have to run
the tests to make sure the generated file works as expected.

In this case, the new file should be located at
`protobuf/protoc-gen-s12perm/example/example.perm.pb.go`.

## Generating crux client c++ code from protobuf

* Install [Go](https://golang.org/doc/install)
* Install Protocol Buffer Compiler via Homebrew `brew install protobuf`
* Install gRPC view Homebrew `brew install gRPC`

```
$ make cruxclient
```

This command will first compile the crux code generator as a protoc plugin and install it to the system bin directory. 
Then it will be used to generate custom crux code in c++.  
