package suite

import (
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"

	"github.com/dmolesUC3/cos/internal/logging"

	"github.com/dmolesUC3/cos/internal/objects"
)

const (
	// if we don't take at least a little time the spinner gets confused
	minTaskTime  = time.Second / time.Duration(8)
	spinCharsStr = "🌑🌒🌓🌔🌕🌖🌗🌘"
)

// ------------------------------------------------------------
// Case

type Case interface {
	Name() string
	RunWithSpinner(index int, target objects.Target, dryRun bool) (detail string)
}

// ------------------------------------------------------------
// Unexported types

type execution func(target objects.Target) (ok bool, detail string)

type caseImpl struct {
	name string
	exec execution
}

func newCase(name string, exec execution) Case {
	return &caseImpl{name, exec}
}

func (c *caseImpl) Name() string {
	return c.name
}

func (c *caseImpl) RunWithSpinner(index int, target objects.Target, dryRun bool) string {
	sp := newSpinner(c.title(index))
	sp.Start()

	elapsed, ok, detail := c.maybeExec(target, dryRun)
	if time.Duration(elapsed) < minTaskTime {
		time.Sleep(minTaskTime - time.Duration(elapsed))
	}

	sp.FinalMSG = c.finalMsg(index, ok, elapsed)
	sp.Stop()

	return detail
}

var spinChars = strings.Split(spinCharsStr, "")
var frameDuration = time.Second / time.Duration(len(spinChars))

func newSpinner(title string) *spinner.Spinner {
	sp := spinner.New(spinChars, frameDuration)
	sp.Suffix = " " + title
	return sp
}

func (c *caseImpl) title(index int) string {
	return fmt.Sprintf("%d. %v", index+1, c.Name())
}

func (c *caseImpl) maybeExec(target objects.Target, dryRun bool) (elapsed int64, ok bool, detail string) {
	start := time.Now().UnixNano()
	if dryRun {
		time.Sleep(minTaskTime)
		ok, detail = true, "(dry run)"
	} else {
		ok, detail = c.exec(target)
	}
	elapsed = time.Now().UnixNano() - start
	return
}

const finalMsgFormat = "%v %d. %v: %v (%v)\n"

func (c *caseImpl) finalMsg(index int, ok bool, elapsed int64) string {
	icon, status := iconAndStatus(ok)
	return fmt.Sprintf(finalMsgFormat, string(icon), index+1, c.Name(), status, logging.FormatNanos(elapsed))
}

func iconAndStatus(ok bool) (rune, string) {
	if ok {
		return '\u2705', "successful"
	}
	return '\u274C', "FAILED"
}
