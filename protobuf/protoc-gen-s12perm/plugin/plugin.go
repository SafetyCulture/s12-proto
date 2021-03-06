// Copyright (c) 2020 SafetyCulture Pty Ltd. All Rights Reserved.

package plugin

import (
	"fmt"
	"strings"

	"google.golang.org/protobuf/compiler/protogen"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/SafetyCulture/s12-proto/s12/flags/permissions"
)

// Standard library dependencies.
const (
	contextPackage = protogen.GoImportPath("context")
	logPackage     = protogen.GoImportPath("log")
)

// Other library dependencies.
const (
	grpcPackage     = protogen.GoImportPath("google.golang.org/grpc")
	s12permPackage  = protogen.GoImportPath("github.com/SafetyCulture/s12-proto/s12/flags/permissions")
	s12utilsPackage = protogen.GoImportPath("github.com/SafetyCulture/s12-utils-go/utils")
)

// GenerateFile generates the .perm.pb.go file
func GenerateFile(p *protogen.Plugin, file *protogen.File) *protogen.GeneratedFile {
	if len(file.Services) == 0 {
		return nil
	}

	filename := file.GeneratedFilenamePrefix + ".perm.pb.go"
	g := p.NewGeneratedFile(filename, file.GoImportPath)

	g.P("// Code generated by protoc-gen-govalidator. DO NOT EDIT.")
	g.P()
	g.P("package ", file.GoPackageName)
	g.P()

	// blank import to register extention for gRPC Reflection API
	g.Import(s12permPackage)

	for _, srv := range file.Services {
		genSrvInterceptor(g, srv)
	}

	return g
}

func genSrvInterceptor(g *protogen.GeneratedFile, srv *protogen.Service) {
	genUnaryInterceptor(g, srv)
	genStreamInterceptor(g, srv)
}

func genUnaryInterceptor(g *protogen.GeneratedFile, srv *protogen.Service) {
	g.P("// ", srv.GoName, "PermissionsUnaryInterceptor is a gRPC unary server interceptor that validates the S12 JWT claims")
	g.P("// for defined permissions for a service method. Returns PermissionDenied status on permission error.")
	g.P("func ", srv.GoName, "PermissionsUnaryInterceptor() ", grpcPackage.Ident("UnaryServerInterceptor"), " {")
	g.P("return func(ctx ", contextPackage.Ident("Context"), ", req interface{}, info *", grpcPackage.Ident("UnaryServerInfo"), ", handler ", grpcPackage.Ident("UnaryHandler"), ") (interface{}, error) {")
	g.P("c, _ := ctx.Value(", s12utilsPackage.Ident("ContextKeyS12JWTClaims"), ").(", s12utilsPackage.Ident("S12JWTClaims"), ")")
	g.P("_ = c")

	for _, md := range srv.Methods {
		if md.Desc.IsStreamingClient() || md.Desc.IsStreamingServer() {
			continue
		}
		flags := getStringSliceExtension(md, permissions.E_RequiredFlags)
		if len(flags) == 0 {
			continue
		}

		sname := fmt.Sprintf("/%s/%s", srv.Desc.FullName(), md.Desc.Name())
		g.P("if info.FullMethod == \"", sname, "\" {")
		var perms []string
		for _, perm := range flags {
			perms = append(perms, fmt.Sprintf("%s(%q)", g.QualifiedGoIdent(s12utilsPackage.Ident("Permission")), perm))
		}
		g.P("if !c.HasPermission(", strings.Join(perms, ", "), ") {")
		g.P(logPackage.Ident("Println"), "(\"s12perm: claims does contain the required permissions\")")
		g.P("return ctx, ", s12utilsPackage.Ident("ErrPermissionDenied"))
		g.P("}")
		g.P("}")
	}

	g.P("return handler(ctx, req)")
	g.P("}")
	g.P("}")
	g.P()
}

func genStreamInterceptor(g *protogen.GeneratedFile, srv *protogen.Service) {
	g.P("// ", srv.GoName, "PermissionsStreamInterceptor is a gRPC stream server interceptor that validates the S12 JWT claims")
	g.P("// for defined permissions for a service method. Returns PermissionDenied status on permission error.")
	g.P("func ", srv.GoName, "PermissionsStreamInterceptor() ", grpcPackage.Ident("StreamServerInterceptor"), " {")
	g.P("return func(srv interface{}, stream ", grpcPackage.Ident("ServerStream"), ", info *", grpcPackage.Ident("StreamServerInfo"), ", handler ", grpcPackage.Ident("StreamHandler"), ") error {")
	g.P("c, _ := stream.Context().Value(", s12utilsPackage.Ident("ContextKeyS12JWTClaims"), ").(", s12utilsPackage.Ident("S12JWTClaims"), ")")
	g.P("_ = c")

	// (rogchap) A lot of this is repeated from the unary interceptor; a refactor to make this more DRY
	// is possible, however it would remove some of the readability, so keeping the repetition.
	for _, md := range srv.Methods {
		if !md.Desc.IsStreamingClient() && !md.Desc.IsStreamingServer() {
			continue
		}
		flags := getStringSliceExtension(md, permissions.E_RequiredFlags)
		if len(flags) == 0 {
			continue
		}

		sname := fmt.Sprintf("/%s/%s", srv.Desc.FullName(), md.Desc.Name())
		g.P("if info.FullMethod == \"", sname, "\" {")
		var perms []string
		for _, perm := range flags {
			perms = append(perms, fmt.Sprintf("%s(%q)", g.QualifiedGoIdent(s12utilsPackage.Ident("Permission")), perm))
		}
		g.P("if !c.HasPermission(", strings.Join(perms, ", "), ") {")
		g.P(logPackage.Ident("Println"), "(\"s12perm: claims does contain the required permissions\")")
		g.P("return ", s12utilsPackage.Ident("ErrPermissionDenied"))
		g.P("}")
		g.P("}")
	}

	g.P("return handler(srv, stream)")
	g.P("}")
	g.P("}")
	g.P()
}

func getStringSliceExtension(md *protogen.Method, xt protoreflect.ExtensionType) []string {
	if opts := md.Desc.Options(); opts != nil {
		ext := proto.GetExtension(opts, xt)
		if v, ok := ext.([]string); ok {
			return v
		}
	}
	return nil
}
