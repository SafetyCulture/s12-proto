
ARCH := $(shell uname -p)
ifeq ($(ARCH), arm)
	SYS_ROOT = /opt/homebrew/
else
	SYS_ROOT = /usr/local/
endif

CXX = g++
CPPFLAGS += -I$(SYS_ROOT)include -pthread
LDFLAGS += -L$(SYS_ROOT)lib -lprotoc -lprotobuf -lpthread -ldl
CXXFLAGS += -std=c++17

protoc-gen-cruxclient: \
protobuf/protoc-gen-cruxclient/cruxclient_generator.o \
protobuf/protoc-gen-cruxclient/api_generator.o \
protobuf/protoc-gen-cruxclient/mock_service_generator.o
	$(CXX) $^ $(LDFLAGS) -o $@

.PHONY: install-logger
	go install github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-logger

.PHONY: install-cruxclient
install-cruxclient: protoc-gen-cruxclient
	install protoc-gen-cruxclient $(SYS_ROOT)/bin/protoc-gen-cruxclient

.PHONY: install
install: install-logger

PLAYGROUND=./protobuf/protoc-gen-govalidator/playground
PLAYGROUND_GEN=$(PLAYGROUND)/gen

playground-req-object: # Generate the code-generator-request object
ifeq ($(shell which protoc-gen-debug),)
	go install github.com/lyft/protoc-gen-star/protoc-gen-debug@latest
endif
	protoc \
		-I$(PLAYGROUND) \
		-I$(GOPATH)/src \
		-I. \
  		--plugin=protoc-gen-debug=$(shell which protoc-gen-debug) \
  		--debug_out="${PLAYGROUND}:${PLAYGROUND}" \
  		$(PLAYGROUND)/*.proto


.PHONY: logger
logger: install-logger
	protoc \
	-I./protobuf/protoc-gen-logger/example \
	-I$(GOPATH)/src \
	--gogo_out=:protobuf/protoc-gen-logger/example \
	--logger_out=:protobuf/protoc-gen-logger/example \
	protobuf/protoc-gen-logger/example/*.proto

.PHONY: cruxclient
cruxclient: install-cruxclient
	protoc \
	-I./protobuf/protoc-gen-cruxclient/proto \
	--plugin=protoc-gen-grpc=$(SYS_ROOT)bin/grpc_cpp_plugin \
	--cpp_out=:protobuf/protoc-gen-cruxclient/generated \
	--grpc_out=:protobuf/protoc-gen-cruxclient/generated \
	--cruxclient_out=:protobuf/protoc-gen-cruxclient/generated \
	protobuf/protoc-gen-cruxclient/proto/routeguide/v1/*.proto

.PHONY: install-gogrpcmock
install-gogrpcmock:
	go install github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-gogrpcmock

.PHONY: gogrpcmock
gogrpcmock: install-gogrpcmock
	protoc \
	-I./protobuf/protoc-gen-gogrpcmock/example \
	-I./protobuf \
	--gogo_out=plugin=grpc=:protobuf/protoc-gen-gogrpcmock/example \
	--gogrpcmock_out=:protobuf/protoc-gen-gogrpcmock/example \
	protobuf/protoc-gen-gogrpcmock/example/*.proto


.PHONY: example
example: logger
