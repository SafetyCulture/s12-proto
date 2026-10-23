// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

import (
	"fmt"
	"sort"
)

const unrecognizedCaseName = "UNRECOGNIZED"

const requiredEnumConformances = "Swift.CaseIterable"

// enumGen generates a Swift enum from a proto enum.
type enumGen struct {
	g            *generator
	desc         *enum
	aliasInfo    *aliasInfo
	mainSorted   []*enumValue
	relativeName string
	fullName     string
}

func (g *generator) newEnumGen(e *enum) *enumGen {
	info := newAliasInfo(e)
	sorted := append([]*enumValue(nil), info.mainValues...)
	sort.SliceStable(sorted, func(i, j int) bool { return sorted[i].number < sorted[j].number })
	return &enumGen{
		g: g, desc: e, aliasInfo: info, mainSorted: sorted,
		relativeName: g.namer.relativeNameEnum(e),
		fullName:     g.namer.fullNameEnum(e),
	}
}

func (e *enumGen) generateMainEnum(p *printer) {
	g := e.g
	visibility := g.opts.visibilitySnippet()
	comments := g.enumComments(e.desc)
	closedRaw := ""
	if e.desc.isClosed() {
		closedRaw = "Int, "
	}
	enumDecl := fmt.Sprintf("%snonisolated enum %s: %s%sEnum, %s {", visibility, e.relativeName, closedRaw, g.namer.runtimePrefix(), requiredEnumConformances)
	if g.opts.enumGeneration != enumGenerationNone && !e.desc.isClosed() {
		p.print("", comments+g.opts.enumAnnotation()+enumDecl)
	} else {
		p.print("", comments+enumDecl)
	}
	p.withIndentation(func() {
		if !e.desc.isClosed() {
			p.print(visibility + "typealias RawValue = Int")
		}
		e.generateCasesOrAliases(p)
		p.print("", visibility+"init() {")
		p.printIndented("self = " + g.namer.dottedRelativeNameEnumValue(e.desc.values[0]))
		p.print("}")
		if !e.desc.isClosed() {
			p.print()
			e.generateInitRawValue(p)
			p.print()
			e.generateRawValueProperty(p)
		}
		e.maybeGenerateCaseIterable(p)
	})
	p.print("", "}")
}

func (e *enumGen) maybeGenerateCaseIterable(p *printer) {
	if e.desc.isClosed() {
		return
	}
	visibility := e.g.opts.visibilitySnippet()
	p.print(
		"",
		"// The compiler won't synthesize support with the "+unrecognizedCaseName+" case.",
		fmt.Sprintf("%sstatic let allCases: [%s] = [", visibility, e.fullName),
	)
	p.withIndentation(func() {
		for _, v := range e.aliasInfo.mainValues {
			p.print(e.g.namer.dottedRelativeNameEnumValue(v) + ",")
		}
	})
	p.print("]")
}

func (e *enumGen) generateRuntimeSupport(p *printer) {
	p.print(
		"",
		fmt.Sprintf("nonisolated extension %s: %s_ProtoNameProviding {", e.fullName, e.g.namer.runtimePrefix()),
	)
	p.withIndentation(func() { e.generateProtoNameProviding(p) })
	p.print("}")
}

func (e *enumGen) generateCasesOrAliases(p *printer) {
	g := e.g
	visibility := g.opts.visibilitySnippet()
	for _, v := range g.namer.uniquelyNamedValues(e.aliasInfo, e.desc) {
		comments := g.enumValueComments(v)
		if comments != "" {
			p.print()
		}
		relativeName := g.namer.relativeNameEnumValue(v)
		switch {
		case e.aliasInfo.original(v) != nil:
			aliasOfName := g.namer.relativeNameEnumValue(e.aliasInfo.original(v))
			p.print(fmt.Sprintf("%s%sstatic let %s = %s", comments, visibility, relativeName, aliasOfName))
		case e.desc.isClosed():
			p.print(fmt.Sprintf("%scase %s = %d", comments, relativeName, v.number))
		default:
			p.print(fmt.Sprintf("%scase %s // = %d", comments, relativeName, v.number))
		}
	}
	if !e.desc.isClosed() {
		p.print("case " + unrecognizedCaseName + "(Int)")
	}
}

func (e *enumGen) generateProtoNameProviding(p *printer) {
	g := e.g
	visibility := g.opts.visibilitySnippet()
	rt := g.namer.runtimePrefix()
	if g.opts.hides(hiddenEnumValues) {
		p.print(fmt.Sprintf("%sstatic let _protobuf_nameMap = %s_NameMap()", visibility, rt))
		return
	}
	w := newNameInstructionWriter()
	for _, v := range e.mainSorted {
		jsonName, hasJSON := v.jsonOption()
		aliases := e.aliasInfo.aliases(v)
		switch {
		case hasJSON:
			if aliases != nil {
				w.writeAliased(v, aliases)
			}
			w.writeUnique(v.number, v.name, jsonName)
		case aliases != nil:
			w.writeAliased(v, aliases)
		default:
			w.writeSame(v.number, v.name)
		}
	}
	if w.err != nil {
		g.err = fmt.Errorf("%s: %w", e.desc.fullName, w.err)
		return
	}
	p.print(fmt.Sprintf("%sstatic let _protobuf_nameMap = %s_NameMap(bytecode: %s)", visibility, rt, w.bytecode.stringLiteral()))
}

func (e *enumGen) generateInitRawValue(p *printer) {
	visibility := e.g.opts.visibilitySnippet()
	p.print(visibility + "init?(rawValue: Int) {")
	p.withIndentation(func() {
		p.print("switch rawValue {")
		for _, v := range e.mainSorted {
			p.print(fmt.Sprintf("case %d: self = %s", v.number, e.g.namer.dottedRelativeNameEnumValue(v)))
		}
		p.print("default: self = ." + unrecognizedCaseName + "(rawValue)")
		p.print("}")
	})
	p.print("}")
}

func (e *enumGen) generateRawValueProperty(p *printer) {
	visibility := e.g.opts.visibilitySnippet()
	neededCases := len(e.mainSorted) + 1
	useMultipleSwitches := neededCases > maxCasesInSwitch
	p.print(visibility + "var rawValue: Int {")
	p.withIndentation(func() {
		if useMultipleSwitches {
			for i, v := range e.mainSorted {
				if i%maxCasesInSwitch == 0 {
					if i > 0 {
						p.print("default: break", "}")
					}
					p.print("switch self {")
				}
				p.print(fmt.Sprintf("case %s: return %d", e.g.namer.dottedRelativeNameEnumValue(v), v.number))
			}
			p.print("case ." + unrecognizedCaseName + "(let i): return i")
			p.print("default: break\n}\n\n// Can't get here, all the cases are listed in the above switches.\n// See https://github.com/apple/swift-protobuf/issues/904 for more details.\nfatalError()")
			return
		}
		p.print("switch self {")
		for _, v := range e.mainSorted {
			p.print(fmt.Sprintf("case %s: return %d", e.g.namer.dottedRelativeNameEnumValue(v), v.number))
		}
		p.print("case ." + unrecognizedCaseName + "(let i): return i")
		p.print("}")
	})
	p.print("}")
}
