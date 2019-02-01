package keys

type keyList struct {
	name string
	desc string
	keys []string
}

func (s *keyList) Name() string {
	return s.name
}

func (s *keyList) Desc() string {
	return s.desc
}

func (s *keyList) Keys() []string {
	return s.keys
}

func (s *keyList) Count() int {
	return len(s.Keys())
}

