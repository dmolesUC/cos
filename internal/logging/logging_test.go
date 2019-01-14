package logging

import (
	"fmt"
	"testing"

	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

// ------------------------------------------------------------
// Helper types

type Prettifiable struct {
	Val interface{}
}

func (p Prettifiable) String() string {
	return fmt.Sprintf("Prettifiable{ Val: %v }", p.Val)
}

func (p Prettifiable) Pretty() string {
	return fmt.Sprintf("prettified %v", p.Val)
}

type StringableWriter interface {
	Write(p []byte) (n int, err error)
	String() string
}

type FailWriter struct {

}

func (f FailWriter) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("failed to write %v", p)
}

func (f FailWriter) String() string {
	return "FailWriter{}"
}
