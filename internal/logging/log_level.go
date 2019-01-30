package logging

type LogLevel int

const (
	Info LogLevel = iota
	Detail
	Trace
	Default = Info
)

func (l LogLevel) String() string {
	if l == Info {
		return "Info"
	}
	if l == Detail {
		return "Detail"
	}
	return "Trace"
}
