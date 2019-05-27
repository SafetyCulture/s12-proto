
.PHONY: generate
generate:
	protoc \
	-I./protobuf/s12proto/:. \
	--gogo_out=Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor:protobuf/s12proto/ \
	protobuf/s12proto/validator.proto

CXX = g++
CPPFLAGS += -I/usr/local/include -pthread
CXXFLAGS += -std=c++11
LDFLAGS += -L/usr/local/lib -lprotoc -lprotobuf -lpthread -ldl
protoc-gen-cruxclient: protobuf/protoc-gen-cruxclient/cruxclient_generator.o
	$(CXX) $^ $(LDFLAGS) -o $@

.PHONY: install-govalidator
install-govalidator:
	go install github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-govalidator

.PHONY: install-logger
	go install github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-logger

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

.PHONY: cruxclient
cruxclient: install-cruxclient
	protoc \
	-I./protobuf/protoc-gen-cruxclient/example \
	--plugin=protoc-gen-grpc=/usr/local/bin/grpc_cpp_plugin \
	--cpp_out=:protobuf/protoc-gen-cruxclient/example \
	--grpc_out=:protobuf/protoc-gen-cruxclient/example \
	--cruxclient_out=:protobuf/protoc-gen-cruxclient/example \
	protobuf/protoc-gen-cruxclient/example/*.proto

.PHONY: example
example: govalidator logger