
.PHONY: generate
generate:
	protoc \
	-I./protobuf/s12proto/:. \
	--gogo_out=Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor:../../../ \
	protobuf/s12proto/*.proto

.PHONY: install
install:
	go install github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-govalidator
	go install github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-logger


.PHONY: govalidator
govalidator: install
	protoc \
	-I./protobuf/protoc-gen-govalidator/example \
	-I$(GOPATH)/src \
	--gogo_out=:protobuf/protoc-gen-govalidator/example \
	--govalidator_out=:protobuf/protoc-gen-govalidator/example \
	protobuf/protoc-gen-govalidator/example/*.proto

.PHONY: logger
logger: install
	protoc \
	-I./protobuf/protoc-gen-logger/example \
	-I$(GOPATH)/src \
	--gogo_out=:protobuf/protoc-gen-logger/example \
	--logger_out=:protobuf/protoc-gen-logger/example \
	protobuf/protoc-gen-logger/example/*.proto

.PHONY: example
example: govalidator logger