
.PHONY: generate
generate:
	protoc \
	-I. \
	--go_out=paths=source_relative:. \
	s12/protobuf/proto/*.proto s12/flags/permissions/*.proto

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
protobuf/protoc-gen-cruxclient/legacy_generator.o \
protobuf/protoc-gen-cruxclient/api_generator.o \
protobuf/protoc-gen-cruxclient/mock_service_generator.o
	$(CXX) $^ $(LDFLAGS) -o $@

.PHONY: install-govalidator
install-govalidator:
	go install ./protobuf/protoc-gen-govalidator

.PHONY: install-logger
	go install github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-logger

.PHONY: install-s12perm
install-s12perm:
	go install ./protobuf/protoc-gen-s12perm

.PHONY: install-cruxclient
install-cruxclient: protoc-gen-cruxclient
	install protoc-gen-cruxclient $(SYS_ROOT)/bin/protoc-gen-cruxclient

.PHONY: install
install: install-govalidator install-logger

.PHONY: govalidator
govalidator: install-govalidator
	protoc \
	-I./protobuf/protoc-gen-govalidator/example \
	-I$(GOPATH)/src \
	-I. \
	--go_out=paths=source_relative:protobuf/protoc-gen-govalidator/example \
	--govalidator_out=paths=source_relative:protobuf/protoc-gen-govalidator/example \
	protobuf/protoc-gen-govalidator/example/*.proto

.PHONY: govalidator-test
govalidator-test:
	go test github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-govalidator/example

.PHONY: govalidator-valtest
govalidator-valtest: install-govalidator
	protoc \
	-I./protobuf/protoc-gen-govalidator/valtest \
	-I$(GOPATH)/src \
	-I. \
	--go_out=paths=source_relative:protobuf/protoc-gen-govalidator/valtest \
	--govalidator_out=paths=source_relative:protobuf/protoc-gen-govalidator/valtest \
	protobuf/protoc-gen-govalidator/valtest/*.proto

.PHONY: govalidator-valtest-test
govalidator-valtest-test:
	go test github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-govalidator/valtest

.PHONY: s12perm
s12perm: install-s12perm
	protoc \
	-I./protobuf/protoc-gen-s12perm/example \
	-I$(GOPATH)/src \
	-I. \
	--s12perm_out=paths=source_relative:protobuf/protoc-gen-s12perm/example \
	protobuf/protoc-gen-s12perm/example/*proto

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
example: govalidator logger
