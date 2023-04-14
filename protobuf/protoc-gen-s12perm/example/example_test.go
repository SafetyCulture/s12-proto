package example_test

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"sc-go.io/pkg/jwtclaims"

	"sc-go.io/pkg/credentials"

	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-s12perm/example"
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

func TestExamplePermissionsUnaryInterceptor(t *testing.T) {
	tests := [...]struct {
		name        string
		scope       []string
		permissions []jwtclaims.Permission
		shouldErr   bool
	}{
		{
			name: "having the permission should not return an error",
			scope: []string{
				"user",
			},
			permissions: []jwtclaims.Permission{
				"write:users",
			},
			shouldErr: false,
		},
		{
			name: "not having the permission should return an error",
			scope: []string{
				"user",
			},
			permissions: []jwtclaims.Permission{
				"write:folders",
			},
			shouldErr: true,
		},
		{
			name: "admin scope should bypass permission check",
			scope: []string{
				"admin",
			},
			permissions: []jwtclaims.Permission{},
			shouldErr:   false,
		},
		{
			name: "empty permissions with user scope should return an error",
			scope: []string{
				"user",
			},
			permissions: []jwtclaims.Permission{},
			shouldErr:   true,
		},
	}

	interceptor := example.ExamplePermissionsUnaryInterceptor()
	noopHandler := func(_ context.Context, _ interface{}) (interface{}, error) {
		return nil, nil
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s12jwt := jwtclaims.S12JWTClaims{}
			s12jwt.BuildScope(tt.permissions)
			ctx := context.WithValue(context.Background(), credentials.ContextKeyCredentialsScope, credentials.Scope(tt.scope))
			ctx = context.WithValue(ctx, jwtclaims.ContextKeyS12JWTClaims, s12jwt)
			info := &grpc.UnaryServerInfo{
				FullMethod: "/example.Example/Unary",
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
