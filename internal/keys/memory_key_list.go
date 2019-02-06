package keys

// ------------------------------------------------------------
// MemoryKeyList

type MemoryKeyList struct {
	name string
	desc string
	keys []string
}

func (s *MemoryKeyList) Name() string {
	return s.name
}

func (s *MemoryKeyList) Desc() string {
	return s.desc
}

func (s *MemoryKeyList) Keys() []string {
	return s.keys
}

func (s *MemoryKeyList) Count() int {
	return len(s.Keys())
}

