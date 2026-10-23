// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

// storageGen generates the storage class of a message that keeps its fields on the heap.
type storageGen struct {
	fields []fieldGenerator
	isAny  bool
}

func (s *storageGen) storageVisibility() string {
	if s.isAny {
		return "internal"
	}
	return "fileprivate"
}

// providesEqualTo reports whether the storage compares itself.
func (s *storageGen) providesEqualTo() bool { return s.isAny }

func (s *storageGen) generateTypeDeclaration(p *printer) {
	if s.isAny {
		p.print("typealias _StorageClass = AnyMessageStorage")
		return
	}
	p.print("fileprivate class _StorageClass {")
	p.withIndentation(func() {
		for _, f := range s.fields {
			f.generateStorage(p)
		}
		p.print("\n  // This property is used as the initial default value for new instances of the type.\n  // The type itself is protecting the reference to its storage via CoW semantics.\n  // This will force a copy to be made of this reference when the first mutation occurs;\n  // hence, it is safe to mark this as `nonisolated(unsafe)`.\n  static nonisolated(unsafe) let defaultInstance = _StorageClass()\n\nprivate init() {}\n")
		p.print("init(copying source: _StorageClass) {")
		p.withIndentation(func() {
			for _, f := range s.fields {
				f.generateStorageClassClone(p)
			}
		})
		p.print("}")
	})
	p.print("}")
}

func (s *storageGen) generateUniqueStorage(p *printer) {
	p.print(s.storageVisibility() + " mutating func _uniqueStorage() -> _StorageClass {")
	p.withIndentation(func() {
		p.print("if !isKnownUniquelyReferenced(&_storage) {")
		p.printIndented("_storage = _StorageClass(copying: _storage)")
		p.print("}", "return _storage")
	})
	p.print("}")
}

func (s *storageGen) generatePreTraverse(p *printer) {
	if s.isAny {
		p.print("try _storage.preTraverse()")
	}
}
