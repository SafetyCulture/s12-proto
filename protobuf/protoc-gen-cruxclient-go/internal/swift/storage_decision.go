// Copyright (c) 2026 SafetyCulture Pty Ltd. All Rights Reserved.

package swift

// Messages with a cost of fields over this limit keep their fields in a storage class.
const totalFieldCostRequiringStorage = 17

type analyzeResult struct {
	usesStorage bool
	costAsField int
}

var useStorageResult = analyzeResult{usesStorage: true, costAsField: 1}

// storageAnalyzer decides which messages use heap storage. It keeps its cache for one run.
type storageAnalyzer struct {
	analysisCache  map[string]analyzeResult
	recursiveCache map[string]bool
}

func newStorageAnalyzer() *storageAnalyzer {
	return &storageAnalyzer{
		analysisCache:  map[string]analyzeResult{"google.protobuf.Any": useStorageResult},
		recursiveCache: map[string]bool{},
	}
}

func (a *storageAnalyzer) shouldUseHeapStorage(m *message) bool {
	return a.analyze(m).usesStorage
}

func (a *storageAnalyzer) estimate(f *field) int {
	if f.isRepeated() {
		return 1
	}
	switch f.typ {
	case typeGroup, typeMessage:
		return a.analyze(f.messageType).costAsField
	}
	return 1
}

func (a *storageAnalyzer) analyze(m *message) analyzeResult {
	if result, ok := a.analysisCache[m.fullName]; ok {
		return result
	}
	result := a.helper(m)
	a.analysisCache[m.fullName] = result
	return result
}

func (a *storageAnalyzer) helper(m *message) analyzeResult {
	if a.containsRecursiveSingularField(m) {
		return useStorageResult
	}
	fieldsCost := 0
	for _, f := range m.fields {
		if f.oneofIndex != nil {
			continue
		}
		fieldsCost += a.estimate(f)
		if fieldsCost >= totalFieldCostRequiringStorage {
			return useStorageResult
		}
	}
	for _, o := range m.oneofs {
		oneofCost := 0
		for _, f := range o.fields {
			if c := a.estimate(f); c > oneofCost {
				oneofCost = c
			}
			if fieldsCost+oneofCost >= totalFieldCostRequiringStorage {
				return useStorageResult
			}
		}
		fieldsCost += oneofCost
	}
	return analyzeResult{costAsField: fieldsCost}
}

func (a *storageAnalyzer) containsRecursiveSingularField(m *message) bool {
	if cached, ok := a.recursiveCache[m.fullName]; ok {
		return cached
	}
	initialFile := m.file
	var recursionHelper func(d *message, stack []*message) bool
	recursionHelper = func(d *message, stack []*message) bool {
		stack = append(append([]*message(nil), stack...), d)
		result := false
		for _, f := range d.fields {
			if f.isRepeated() {
				continue
			}
			if f.typ != typeMessage && f.typ != typeGroup {
				continue
			}
			messageType := f.messageType
			if messageType == nil || messageType.file != initialFile {
				continue
			}
			if known, ok := a.recursiveCache[messageType.fullName]; ok && !known {
				continue
			}
			first := -1
			for i, s := range stack {
				if s == messageType {
					first = i
					break
				}
			}
			if first >= 0 {
				for _, s := range stack[first:] {
					a.recursiveCache[s.fullName] = true
					a.analysisCache[s.fullName] = useStorageResult
				}
				if first == 0 {
					result = true
					break
				}
				continue
			}
			if recursionHelper(messageType, stack) {
				result = true
				break
			}
		}
		// Keep the upstream cache overwrite behavior.
		a.recursiveCache[d.fullName] = result
		return result
	}
	return recursionHelper(m, nil)
}
