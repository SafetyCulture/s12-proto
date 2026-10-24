package example_test

import (
	"context"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sc-go.io/pkg/credentials"
	"sc-go.io/pkg/jwtclaims"

	"github.com/SafetyCulture/s12-proto/protobuf/protoc-gen-s12perm/example"
)

// checkResult asserts a deny is PermissionDenied and an allow carries no error.
func checkResult(t *testing.T, err error, shouldErr bool) {
	t.Helper()
	switch {
	case shouldErr && err == nil:
		t.Error("An error is expected")
	case shouldErr:
		if code := status.Code(err); code != codes.PermissionDenied {
			t.Errorf("Error code = %v, want %v", code, codes.PermissionDenied)
		}
	case err != nil:
		t.Errorf("No error is expected but received: %v", err)
	}
}

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
			checkResult(t, err, tt.shouldErr)
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
			name: "admin permission should satisfy a write permission",
			scope: []string{
				"user",
			},
			permissions: []jwtclaims.Permission{
				jwtclaims.Permission("admin:users"),
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
			// The S12-JWT middleware stores prm verbatim; BuildScope does not
			// model that path.
			s12jwt := jwtclaims.S12JWTClaims{Permissions: tt.permissions}
			ctx := context.WithValue(context.Background(), credentials.ContextKeyCredentialsScope, credentials.Scope(tt.scope))
			ctx = context.WithValue(ctx, jwtclaims.ContextKeyS12JWTClaims, s12jwt)
			info := &grpc.UnaryServerInfo{
				FullMethod: "/example.Example/Unary",
			}
			_, err := interceptor(ctx, nil, info, noopHandler)
			checkResult(t, err, tt.shouldErr)
		})
	}
}

// stubServerStream carries only a context: the generated interceptors read
// nothing else from the stream, so the embedded grpc.ServerStream stays nil.
type stubServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s stubServerStream) Context() context.Context { return s.ctx }

func TestExampleWithScopesPermissionsStreamInterceptor(t *testing.T) {
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
			name: "admin permission should satisfy a write permission",
			scope: []string{
				"user",
			},
			permissions: []jwtclaims.Permission{
				jwtclaims.Permission("admin:users"),
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
			name: "empty permissions should return an error",
			scope: []string{
				"user",
			},
			permissions: []jwtclaims.Permission{},
			shouldErr:   true,
		},
		{
			// The admin scope bypasses the permission check but does not
			// satisfy the user scope this method also requires.
			name: "admin scope without the user scope should return an error",
			scope: []string{
				"admin",
			},
			permissions: []jwtclaims.Permission{},
			shouldErr:   true,
		},
		{
			name: "admin scope with the user scope should bypass the permission check",
			scope: []string{
				"admin",
				"user",
			},
			permissions: []jwtclaims.Permission{},
			shouldErr:   false,
		},
		{
			name: "missing the required user scope should return an error",
			scope: []string{
				"service",
			},
			permissions: []jwtclaims.Permission{
				"write:users",
			},
			shouldErr: true,
		},
	}

	interceptor := example.ExampleWithScopesPermissionsStreamInterceptor()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s12jwt := jwtclaims.S12JWTClaims{Permissions: tt.permissions}
			ctx := context.WithValue(context.Background(), credentials.ContextKeyCredentialsScope, credentials.Scope(tt.scope))
			ctx = context.WithValue(ctx, jwtclaims.ContextKeyS12JWTClaims, s12jwt)
			info := &grpc.StreamServerInfo{
				FullMethod: "/example.ExampleWithScopes/ServerStream",
			}
			handled := false
			handler := func(_ interface{}, _ grpc.ServerStream) error {
				handled = true
				return nil
			}
			err := interceptor(nil, stubServerStream{ctx: ctx}, info, handler)
			checkResult(t, err, tt.shouldErr)
			if wantHandled := !tt.shouldErr; handled != wantHandled {
				t.Errorf("Handler ran = %v, want %v", handled, wantHandled)
			}
		})
	}
}

// TestExamplePermissionsStreamInterceptor covers Example/ServerStream, which
// requires write:users and declares no scope, so an admin-scoped token with no
// claims is allowed.
func TestExamplePermissionsStreamInterceptor(t *testing.T) {
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
			name: "admin permission should satisfy a write permission",
			scope: []string{
				"user",
			},
			permissions: []jwtclaims.Permission{
				jwtclaims.Permission("admin:users"),
			},
			shouldErr: false,
		},
		{
			name: "admin scope should bypass the permission check",
			scope: []string{
				"admin",
			},
			permissions: []jwtclaims.Permission{},
			shouldErr:   false,
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
			name: "empty permissions should return an error",
			scope: []string{
				"user",
			},
			permissions: []jwtclaims.Permission{},
			shouldErr:   true,
		},
	}

	interceptor := example.ExamplePermissionsStreamInterceptor()
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s12jwt := jwtclaims.S12JWTClaims{Permissions: tt.permissions}
			ctx := context.WithValue(context.Background(), credentials.ContextKeyCredentialsScope, credentials.Scope(tt.scope))
			ctx = context.WithValue(ctx, jwtclaims.ContextKeyS12JWTClaims, s12jwt)
			info := &grpc.StreamServerInfo{
				FullMethod: "/example.Example/ServerStream",
			}
			handled := false
			handler := func(_ interface{}, _ grpc.ServerStream) error {
				handled = true
				return nil
			}
			err := interceptor(nil, stubServerStream{ctx: ctx}, info, handler)
			checkResult(t, err, tt.shouldErr)
			if wantHandled := !tt.shouldErr; handled != wantHandled {
				t.Errorf("Handler ran = %v, want %v", handled, wantHandled)
			}
		})
	}
}
