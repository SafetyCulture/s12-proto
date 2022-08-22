// Code generated by protoc-gen-govalidator. DO NOT EDIT.
// versions:
// 	protoc-gen-govalidator v2.5.3
// 	protoc                 v3.21.5
// source: geo.proto

package valtest

import (
	fmt "fmt"
	proto "github.com/SafetyCulture/s12-proto/s12/protobuf/proto"
)

func (m *GeoValidationMessage) Validate() error {
	if m.Latitude != 0 {
		// Range check lower bounds
		if m.Latitude < -90 {
			return fmt.Errorf(`latitude: value must be greater than or equal to -90`)
		}
		// Range check upper bounds
		if m.Latitude > 90 {
			return fmt.Errorf(`latitude: value must be less than or equal to 90`)
		}
	}
	if m.Longitude != 0 {
		// Range check lower bounds
		if m.Longitude < -180 {
			fmt.Printf("[log-only] %s: value must %s: Base64Encoded input: %s\n", "longitude", "be greater than or equal to -180", proto.Base64Encode(proto.FirstCharactersFromString(fmt.Sprintf("%v", m.Longitude), 50)))
		}
		// Range check upper bounds
		if m.Longitude > 180 {
			fmt.Printf("[log-only] %s: value must %s: Base64Encoded input: %s\n", "longitude", "be less than or equal to 180", proto.Base64Encode(proto.FirstCharactersFromString(fmt.Sprintf("%v", m.Longitude), 50)))
		}
	}
	if m.Accuracy != 0 {
		// Range check lower bounds
		if m.Accuracy < 0 {
			fmt.Printf("[log-only] %s: value must %s: Base64Encoded input: %s\n", "accuracy", "be greater than or equal to 0", proto.Base64Encode(proto.FirstCharactersFromString(fmt.Sprintf("%v", m.Accuracy), 50)))
		}
		// Range check upper bounds
		if m.Accuracy > 10000 {
			fmt.Printf("[log-only] %s: value must %s: Base64Encoded input: %s\n", "accuracy", "be less than or equal to 10000", proto.Base64Encode(proto.FirstCharactersFromString(fmt.Sprintf("%v", m.Accuracy), 50)))
		}
	}
	return nil
}
