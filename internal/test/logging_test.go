package test

import (
	"fmt"
)

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
