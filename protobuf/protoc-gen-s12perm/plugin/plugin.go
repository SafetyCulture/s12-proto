// Copyright (c) 2020 SafetyCulture Pty Ltd. All Rights Reserved.

package plugin

import (
	"fmt"
	"strings"

	perms "github.com/SafetyCulture/s12-proto/protobuf/s12proto"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
)

type s12perm struct {
	*generator.Generator

	imports generator.PluginImports

	ctxPkg   generator.Single
	logPkg   generator.Single
	grpcPkg  generator.Single
	utilsPkg generator.Single
}

func New() generator.Plugin {
	return &s12perm{}
}

func (g *s12perm) Name() string {
	return "s12perm"
}

func (g *s12perm) Init(gen *generator.Generator) {
	g.Generator = gen
}

func (g *s12perm) GenerateImports(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}
	g.imports.GenerateImports(file)
}

func (g *s12perm) Generate(file *generator.FileDescriptor) {
	if len(file.FileDescriptorProto.Service) == 0 {
		return
	}

	g.imports = generator.NewPluginImports(g.Generator)
	g.logPkg = g.imports.NewImport("log")
	g.ctxPkg = g.imports.NewImport("context")
	g.utilsPkg = g.imports.NewImport("github.com/SafetyCulture/s12-utils-go/utils")
	g.grpcPkg = g.imports.NewImport("google.golang.org/grpc")

	for _, service := range file.FileDescriptorProto.Service {
		g.P("func ", service.GetName(), "PermissionsUnaryInterceptor() ", g.grpcPkg.Use(), ".UnaryServerInterceptor {")
		g.In()
		g.P("return func(ctx ", g.ctxPkg.Use(), ".Context, req interface{}, info *", g.grpcPkg.Use(), ".UnaryServerInfo, handler ", g.grpcPkg.Use(), ".UnaryHandler) (interface{}, error) {")
		g.In()
		g.P("claims, _ := ctx.Value(", g.utilsPkg.Use(), ".ContextKeyS12JWTClaims).(", g.utilsPkg.Use(), ".S12JWTClaims)")
		g.P("_ = claims")

		for _, method := range service.Method {
			if method.GetServerStreaming() || method.GetClientStreaming() {
				// TODO: [RC] create streaming interceptor.
				continue
			}
			flags := getPermissions(method)
			if len(flags) == 0 {
				continue
			}
			g.P("if info.FullMethod == \"/", file.GetPackage(), ".", service.GetName(), "/", method.GetName(), "\" {")
			g.In()

			var perms []string
			for _, perm := range flags {
				perms = append(perms, fmt.Sprintf("%s.Permission(%q)", g.utilsPkg.Use(), perm))
			}
			g.P("if !claims.HasPermission(", strings.Join(perms, ", "), ") {")
			g.In()

			g.P(g.logPkg.Use(), ".Println(\"s12perm: claims does contain the required permissions\")")
			g.P("return ctx, ", g.utilsPkg.Use(), ".ErrPermissionDenied")

			g.Out()
			g.P("}")

			g.Out()
			g.P("}\n")
		}

		g.P("return handler(ctx, req)")
		g.Out()
		g.P("}")
		g.Out()
		g.P("}")
	}

}

func getPermissions(method *descriptor.MethodDescriptorProto) []string {
	if method.Options == nil {
		return nil
	}
	v, err := proto.GetExtension(method.Options, perms.E_RequiredFlags)
	if err != nil {
		// option is optional so continue
		return nil
	}
	return v.([]string)
}
