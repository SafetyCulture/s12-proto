
.PHONY: generate
generate:
	protoc \
	-I./protobuf/s12proto/:. \
	--gogo_out=Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor:../../../ \
	protobuf/s12proto/*.proto

.PHONY: install
install:
	go install github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-govalidator


.PHONY: example
example: install
	protoc \
	-I./protobuf/protoc-gen-govalidator/example \
	-I$(GOPATH)/src \
	--gogo_out=:protobuf/protoc-gen-govalidator/example \
	--govalidator_out=:protobuf/protoc-gen-govalidator/example \
	protobuf/protoc-gen-govalidator/example/*.proto