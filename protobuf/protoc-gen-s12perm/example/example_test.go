package example_test

import (
	"context"
	"google.golang.org/grpc"
	"testing"

	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-s12perm/example"
	"sc-go.io/pkg/credentials"
)

func TestExampleWithScopesPermissionsUnaryInterceptor(t *testing.T) {
	tests := [...]struct {
		name      string
		scope     []string
		shouldErr bool
	}{
		{
			name: "valid scope should not return error",
			scope: []string{
				"admin",
			},
			shouldErr: false,
		},
		{
			name: "invalid scope should return an error",
			scope: []string{
				"user",
			},
			shouldErr: true,
		},
		{
			name:      "empty scope should return an error",
			scope:     []string{},
			shouldErr: true,
		},
	}

	interceptor := example.ExampleWithScopesPermissionsUnaryInterceptor()
	noopHandler := func(_ context.Context, _ interface{}) (interface{}, error) {
		return nil, nil
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.WithValue(context.Background(), credentials.ContextKeyCredentialsScope, credentials.Scope(tt.scope))
			info := &grpc.UnaryServerInfo{
				FullMethod: "/example.ExampleWithScopes/Unary",
			}
			_, err := interceptor(ctx, nil, info, noopHandler)
			if tt.shouldErr && err == nil {
				t.Error("An error is expected")
			}
			if !tt.shouldErr && err != nil {
				t.Errorf("No error is expected but received: %v", err)
			}
		})
	}
}