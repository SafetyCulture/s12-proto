
.PHONY: generate
generate:
	protoc \
	-I. \
	--go_out=paths=source_relative:. \
	s12/protobuf/proto/*.proto s12/flags/permissions/*.proto

CXX = g++
CPPFLAGS += -I/usr/local/include -pthread
CXXFLAGS += -std=c++17
LDFLAGS += -L/usr/local/lib -lprotoc -lprotobuf -lpthread -ldl
protoc-gen-cruxclient: protobuf/protoc-gen-cruxclient/cruxclient_generator.o protobuf/protoc-gen-cruxclient/legacy_generator.o protobuf/protoc-gen-cruxclient/api_generator.o
	$(CXX) $^ $(LDFLAGS) -o $@

.PHONY: install-govalidator
install-govalidator:
	go install github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-govalidator

.PHONY: install-logger
	go install github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-logger

.PHONY: install-s12perm
install-s12perm:
	go install ./protobuf/protoc-gen-s12perm

.PHONY: install-cruxclient
install-cruxclient: protoc-gen-cruxclient
	install protoc-gen-cruxclient /usr/local/bin/protoc-gen-cruxclient

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
	go test github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-govalidator/example -v

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
	--plugin=protoc-gen-grpc=/usr/local/bin/grpc_cpp_plugin \
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
