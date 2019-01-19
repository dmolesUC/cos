package logging

import "regexp"

type Pretty interface {
	Pretty() string
}

func Prettify(a ...interface{}) []interface{} {
	var pretty []interface{}
	for _, v := range a {
		if p, ok := v.(Pretty); ok {
			pretty = append(pretty, p.Pretty())
		} else {
			pretty = append(pretty, v)
		}
	}
	return pretty
}

func Untabify(text string, indent string) string {
	// TODO: support multi-level indent
	return regexp.MustCompile(`(?m)^[\t ]+`).ReplaceAllString(text, indent)
}

func PrettyStrP(s *string) string {
	if s == nil {
		return "<nil>"
	}
	return *s
}
