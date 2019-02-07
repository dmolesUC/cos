package suite

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"math/rand"
	"sort"
	"time"

	. "code.cloudfoundry.org/bytefmt"

	"github.com/dmolesUC3/cos/internal/logging"

	. "github.com/dmolesUC3/cos/pkg"

	"github.com/dmolesUC3/cos/internal/objects"
)

// ------------------------------------------------------------
// TestTask

type TestTask interface {
	Title() string
	Invoke(target objects.Target) (ok bool, detail string)
}

func AllTasks(sizeMax int64, countMax uint64) []TestTask {
	tasks := singleFileTasks(sizeMax)
	tasks = append(tasks, multipleFileTasks(countMax)...)
	return tasks
}

// ------------------------------------------------------------
// Task implementations

func singleFileTasks(sizeMax int64) []TestTask {
	tasks := []TestTask{singleFileTask(0)}
	for _, unit := range []int64{BYTE, KILOBYTE, MEGABYTE, GIGABYTE} {
		if unit > sizeMax {
			break
		}
		for _, multiplier := range []int64{1, 16, 256} {
			size := multiplier * unit
			if size > sizeMax {
				break
			}
			tasks = append(tasks, singleFileTask(size))
		}
	}
	return tasks
}

func singleFileTask(size int64) TestTask {
	title := fmt.Sprintf("create/retrieve/verify/delete %v file", logging.FormatBytes(size))
	runTest := func(target objects.Target) (ok bool, detail string) {
		crvd := NewCrvd(target, "", size, DefaultRandomSeed)
		err := crvd.CreateRetrieveVerifyDelete()
		if err == nil {
			return true, ""
		} else {
			return false, err.Error()
		}
	}
	return newTestTask(title, runTest)
}

func multipleFileTasks(countMax uint64) []TestTask {
	var tasks []TestTask
	for i := 2; i <= 12; i++ {
		count := uint64(1) << uint64(2*i)
		if count > countMax {
			break
		}
		tasks = append(tasks, multipleFileTask("prefix", count))
	}
	return tasks
}

func multipleFileTask(prefix string, count uint64) TestTask {
	title := fmt.Sprintf("create %d files under single prefix", count)

	var contentLength int64 = DefaultContentLengthBytes
	body := make([]byte, contentLength)
	rand.Read(body)
	bodyProvider := func() io.Reader {
		return bytes.NewReader(body)
	}


	runTest := func(target objects.Target) (ok bool, detail string) {
		var keysToDelete []string
		defer func() {
			for _, k := range keysToDelete {
				_ = target.Object(k).Delete()
			}
		}()

		// TODO: any out-of-the-box timing tools?
		times := make([]int64, count)

		for i := uint64(0); i < count; i++ {
			start := time.Now().UnixNano()

			key := fmt.Sprintf("%v/file-%d.bin", prefix, i)
			keysToDelete = append(keysToDelete, key)

			crvd := Crvd{
				Object:        target.Object(key),
				ContentLength: contentLength,
				BodyProvider:  bodyProvider,
			}
			err := crvd.CreateRetrieveVerify()
			if err != nil {
				return false, err.Error()
			}

			times[i] = time.Now().UnixNano() - start
		}

		first := times[0]
		last := times[count-1]
		sort.Slice(times, func(i, j int) bool { return times[i] < times[j] })
		fastest := times[0]
		slowest := times[count-1]
		median := int64(math.Round(float64(times[count/2]+times[count/2-1]) / 2))

		return true, fmt.Sprintf("first: %v, last: %v, fastest: %v, slowest: %v, median: %v",
			logging.FormatNanos(first),
			logging.FormatNanos(last),
			logging.FormatNanos(fastest),
			logging.FormatNanos(slowest),
			logging.FormatNanos(median),
		)
	}

	return newTestTask(title, runTest)
}

// ------------------------------------------------------------
// Unexported types

type runTest func(target objects.Target) (ok bool, detail string)

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

func (t *testTask) Invoke(target objects.Target) (ok bool, detail string) {
	return t.runTest(target)
}
