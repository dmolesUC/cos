package keys

type simpleKeyList struct {
	name string
	desc string
	keys []string
}

func (s *simpleKeyList) Name() string {
	return s.name
}

func (s *simpleKeyList) Desc() string {
	return s.desc
}

func (s *simpleKeyList) Keys() []string {
	return s.keys
}

func (s *simpleKeyList) Count() int {
	return len(s.Keys())
}

