package suite

import (
	"fmt"

	. "code.cloudfoundry.org/bytefmt"
	"github.com/dmolesUC3/cos/internal/logging"

	"github.com/dmolesUC3/cos/pkg"

	"github.com/dmolesUC3/cos/internal/objects"
)

// ------------------------------------------------------------
// TestTask

type TestTask interface {
	Title() string
	Invoke(target objects.Target) (ok bool, err error)
}

func AllTasks() []TestTask {
	return []TestTask{
		newCrvdTask(0),
		newCrvdTask(KILOBYTE),
		newCrvdTask(MEGABYTE),
		newCrvdTask(128 * MEGABYTE),
		newCrvdTask(GIGABYTE),
		newCrvdTask(8 * GIGABYTE),
	}
}

// ------------------------------------------------------------
// Task implementations

func newCrvdTask(size int64) TestTask {
	title := fmt.Sprintf("create/retrieve/verify/delete %v file", logging.FormatBytes(size))
	runTest := func(target objects.Target) (ok bool, err error) {
		crvd := pkg.NewCrvd(target, "", size, pkg.DefaultRandomSeed)
		err = crvd.CreateRetrieveVerifyDelete()
		return err == nil, err
	}
	return newTestTask(title, runTest)
}

// ------------------------------------------------------------
// Unexported types

type runTest func(target objects.Target) (ok bool, err error)

type testTask struct {
	title   string
	runTest runTest
}

func newTestTask(title string, runTest runTest) TestTask {
	return &testTask{title, runTest}
}

func (t *testTask) Title() string {
	return t.title
}

func (t *testTask) Invoke(target objects.Target) (ok bool, err error) {
	return t.runTest(target)
}
