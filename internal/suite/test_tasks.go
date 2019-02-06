package suite

import (
	"github.com/dmolesUC3/cos/pkg"

	"github.com/dmolesUC3/cos/internal/objects"
)

type TestTask interface {
	Title() string
	Invoke(target objects.Target) (ok bool, err error)
}

func AllTasks() []TestTask {
	return []TestTask{
		crvdTask{},
	}
}

type crvdTask struct {
}

func (t crvdTask) Title() string {
	return "create, retrieve, verify, delete"
}

func (t crvdTask) Invoke(target objects.Target) (ok bool, err error) {
	crvd := pkg.NewDefaultCrvd(target, "")
	err = crvd.CreateRetrieveVerifyDelete()
	return err == nil, err
}
