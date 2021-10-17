// Code generated by protoc-gen-govalidator. DO NOT EDIT.

package example

import (
	context "context"
	_ "github.com/SafetyCulture/s12-proto/s12/flags/permissions"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	log "log"
	jwtclaims "sc-go.io/pkg/jwtclaims"
)

// ExamplePermissionsUnaryInterceptor is a gRPC unary server interceptor that validates the S12 JWT claims
// for defined permissions for a service method. Returns PermissionDenied status on permission error.
func ExamplePermissionsUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		c, _ := ctx.Value(jwtclaims.ContextKeyS12JWTClaims).(jwtclaims.S12JWTClaims)
		_ = c
		if info.FullMethod == "/example.Example/Unary" {
			if !c.HasPermission(jwtclaims.Permission("write:users")) {
				log.Println("s12perm: claims does contain the required permissions")
				return ctx, status.Errorf(codes.PermissionDenied, "Permission Denied")
			}
		}
		return handler(ctx, req)
	}
}

// ExamplePermissionsStreamInterceptor is a gRPC stream server interceptor that validates the S12 JWT claims
// for defined permissions for a service method. Returns PermissionDenied status on permission error.
func ExamplePermissionsStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		c, _ := stream.Context().Value(jwtclaims.ContextKeyS12JWTClaims).(jwtclaims.S12JWTClaims)
		_ = c
		if info.FullMethod == "/example.Example/ServerStream" {
			if !c.HasPermission(jwtclaims.Permission("write:users")) {
				log.Println("s12perm: claims does contain the required permissions")
				return status.Errorf(codes.PermissionDenied, "Permission Denied")
			}
		}
		return handler(srv, stream)
	}
}

// NoPermissionsPermissionsUnaryInterceptor is a gRPC unary server interceptor that validates the S12 JWT claims
// for defined permissions for a service method. Returns PermissionDenied status on permission error.
func NoPermissionsPermissionsUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		c, _ := ctx.Value(jwtclaims.ContextKeyS12JWTClaims).(jwtclaims.S12JWTClaims)
		_ = c
		return handler(ctx, req)
	}
}

// NoPermissionsPermissionsStreamInterceptor is a gRPC stream server interceptor that validates the S12 JWT claims
// for defined permissions for a service method. Returns PermissionDenied status on permission error.
func NoPermissionsPermissionsStreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		c, _ := stream.Context().Value(jwtclaims.ContextKeyS12JWTClaims).(jwtclaims.S12JWTClaims)
		_ = c
		return handler(srv, stream)
	}
}
