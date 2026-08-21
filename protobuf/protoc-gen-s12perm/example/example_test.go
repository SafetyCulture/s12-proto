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

// runPermissionsUnaryInterceptorTests drives interceptor for fullMethod once per case,
// building claims from the case's permissions and scope.
func runPermissionsUnaryInterceptorTests(t *testing.T, interceptor grpc.UnaryServerInterceptor, fullMethod string, tests []permissionsUnaryTest) {
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
				FullMethod: fullMethod,
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

type permissionsUnaryTest struct {
	name        string
	scope       []string
	permissions []jwtclaims.Permission
	shouldErr   bool
}

// any_of: ["write:users", "write:folders"]
func TestExampleWithAnyOfPermissionsPermissionsUnaryInterceptor(t *testing.T) {
	tests := []permissionsUnaryTest{
		{
			name:        "having only the first any_of permission should not return an error",
			scope:       []string{"user"},
			permissions: []jwtclaims.Permission{"write:users"},
			shouldErr:   false,
		},
		{
			name:        "having only the second any_of permission should not return an error",
			scope:       []string{"user"},
			permissions: []jwtclaims.Permission{"write:folders"},
			shouldErr:   false,
		},
		{
			name:        "having both any_of permissions should not return an error",
			scope:       []string{"user"},
			permissions: []jwtclaims.Permission{"write:users", "write:folders"},
			shouldErr:   false,
		},
		{
			name:        "having a higher privilege on an any_of domain should not return an error",
			scope:       []string{"user"},
			permissions: []jwtclaims.Permission{"admin:folders"},
			shouldErr:   false,
		},
		{
			name:        "having neither any_of permission should return an error",
			scope:       []string{"user"},
			permissions: []jwtclaims.Permission{"write:sensors"},
			shouldErr:   true,
		},
		{
			name:        "having a lower privilege on an any_of domain should return an error",
			scope:       []string{"user"},
			permissions: []jwtclaims.Permission{"read:folders"},
			shouldErr:   true,
		},
		{
			name:        "admin scope should bypass the any_of check",
			scope:       []string{"admin"},
			permissions: []jwtclaims.Permission{},
			shouldErr:   false,
		},
		{
			name:        "empty permissions with user scope should return an error",
			scope:       []string{"user"},
			permissions: []jwtclaims.Permission{},
			shouldErr:   true,
		},
	}

	runPermissionsUnaryInterceptorTests(t, example.ExampleWithAnyOfPermissionsPermissionsUnaryInterceptor(),
		"/example.ExampleWithAnyOfPermissions/Unary", tests)
}

// all_of: ["read:users", "read:sensors"] and any_of: ["write:folders", "write:assets"]
func TestExampleWithMixedPermissionsPermissionsUnaryInterceptor(t *testing.T) {
	tests := []permissionsUnaryTest{
		{
			name:        "having every all_of permission and one any_of permission should not return an error",
			scope:       []string{"user"},
			permissions: []jwtclaims.Permission{"read:users", "read:sensors", "write:folders"},
			shouldErr:   false,
		},
		{
			name:        "having every all_of permission and the other any_of permission should not return an error",
			scope:       []string{"user"},
			permissions: []jwtclaims.Permission{"read:users", "read:sensors", "write:assets"},
			shouldErr:   false,
		},
		{
			name:        "missing one all_of permission should return an error",
			scope:       []string{"user"},
			permissions: []jwtclaims.Permission{"read:users", "write:folders"},
			shouldErr:   true,
		},
		{
			name:        "having every all_of permission but neither any_of permission should return an error",
			scope:       []string{"user"},
			permissions: []jwtclaims.Permission{"read:users", "read:sensors"},
			shouldErr:   true,
		},
		{
			name:        "having both any_of permissions but no all_of permission should return an error",
			scope:       []string{"user"},
			permissions: []jwtclaims.Permission{"write:folders", "write:assets"},
			shouldErr:   true,
		},
		{
			name:        "admin scope should bypass both requirements",
			scope:       []string{"admin"},
			permissions: []jwtclaims.Permission{},
			shouldErr:   false,
		},
	}

	runPermissionsUnaryInterceptorTests(t, example.ExampleWithMixedPermissionsPermissionsUnaryInterceptor(),
		"/example.ExampleWithMixedPermissions/Unary", tests)
}

// any_of: ["write:users", "write:folders"] and any_of: ["write:sensors", "write:assets"]
func TestExampleWithTwoAnyOfPermissionsPermissionsUnaryInterceptor(t *testing.T) {
	tests := []permissionsUnaryTest{
		{
			name:        "one permission from each requirement should not return an error",
			scope:       []string{"user"},
			permissions: []jwtclaims.Permission{"write:users", "write:sensors"},
			shouldErr:   false,
		},
		{
			name:        "the other permission from each requirement should not return an error",
			scope:       []string{"user"},
			permissions: []jwtclaims.Permission{"write:folders", "write:assets"},
			shouldErr:   false,
		},
		{
			name:        "both permissions from only the first requirement should return an error",
			scope:       []string{"user"},
			permissions: []jwtclaims.Permission{"write:users", "write:folders"},
			shouldErr:   true,
		},
		{
			name:        "both permissions from only the second requirement should return an error",
			scope:       []string{"user"},
			permissions: []jwtclaims.Permission{"write:sensors", "write:assets"},
			shouldErr:   true,
		},
		{
			name:        "admin scope should bypass both requirements",
			scope:       []string{"admin"},
			permissions: []jwtclaims.Permission{},
			shouldErr:   false,
		},
	}

	runPermissionsUnaryInterceptorTests(t, example.ExampleWithTwoAnyOfPermissionsPermissionsUnaryInterceptor(),
		"/example.ExampleWithTwoAnyOfPermissions/Unary", tests)
}

// required_flags "read:users" and any_of: ["write:folders", "write:assets"]
func TestExampleWithFlagsAndPermissionsPermissionsUnaryInterceptor(t *testing.T) {
	tests := []permissionsUnaryTest{
		{
			name:        "satisfying the flag and one any_of permission should not return an error",
			scope:       []string{"user"},
			permissions: []jwtclaims.Permission{"read:users", "write:folders"},
			shouldErr:   false,
		},
		{
			name:        "satisfying the flag but neither any_of permission should return an error",
			scope:       []string{"user"},
			permissions: []jwtclaims.Permission{"read:users"},
			shouldErr:   true,
		},
		{
			name:        "satisfying an any_of permission but not the flag should return an error",
			scope:       []string{"user"},
			permissions: []jwtclaims.Permission{"write:folders"},
			shouldErr:   true,
		},
		{
			name:        "admin scope should bypass both the flag and the requirement",
			scope:       []string{"admin"},
			permissions: []jwtclaims.Permission{},
			shouldErr:   false,
		},
	}

	runPermissionsUnaryInterceptorTests(t, example.ExampleWithFlagsAndPermissionsPermissionsUnaryInterceptor(),
		"/example.ExampleWithFlagsAndPermissions/Unary", tests)
}
