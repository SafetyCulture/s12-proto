
.PHONY: generate
generate:
	protoc \
	-I./protobuf/s12proto/:. \
	--gogo_out=Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor:../../../ \
	protobuf/s12proto/*.proto

CXX = g++
CPPFLAGS += -I/usr/local/include -pthread
CXXFLAGS += -std=c++11
LDFLAGS += -L/usr/local/lib -lprotoc -lprotobuf -lpthread -ldl
protoc-gen-cppservice: protobuf/protoc-gen-cppservice/cppservice_generator.o
	$(CXX) $^ $(LDFLAGS) -o $@

.PHONY: install-govalidator
install-govalidator:
	go install github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-govalidator

.PHONY: install-logger
	go install github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-logger

.PHONY: install-cppservice
install-cppservice: protoc-gen-cppservice
	install protoc-gen-cppservice /usr/local/bin/protoc-gen-cppservice

.PHONY: install
install: install-govalidator install-logger

.PHONY: govalidator
govalidator: install-govalidator
	protoc \
	-I./protobuf/protoc-gen-govalidator/example \
	-I$(GOPATH)/src \
	--gogo_out=:protobuf/protoc-gen-govalidator/example \
	--govalidator_out=:protobuf/protoc-gen-govalidator/example \
	protobuf/protoc-gen-govalidator/example/*.proto

.PHONY: logger
logger: install-logger
	protoc \
	-I./protobuf/protoc-gen-logger/example \
	-I$(GOPATH)/src \
	--gogo_out=:protobuf/protoc-gen-logger/example \
	--logger_out=:protobuf/protoc-gen-logger/example \
	protobuf/protoc-gen-logger/example/*.proto

.PHONY: cppservice
cppservice: install-cppservice
	protoc \
	-I./protobuf/protoc-gen-cppservice/example \
	-I$(GOPATH)/src \
	--plugin=protoc-gen-grpc=/usr/local/bin/grpc_cpp_plugin \
	--cpp_out=:protobuf/protoc-gen-cppservice/example \
	--grpc_out=:protobuf/protoc-gen-cppservice/example \
	--cppservice_out=:protobuf/protoc-gen-cppservice/example \
	protobuf/protoc-gen-cppservice/example/*.proto

.PHONY: example
example: govalidator logger