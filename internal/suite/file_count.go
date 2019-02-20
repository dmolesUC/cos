package suite

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"math/rand"
	"sort"
	"time"

	"github.com/dmolesUC3/cos/internal/logging"

	. "github.com/dmolesUC3/cos/pkg"

	"github.com/dmolesUC3/cos/internal/objects"
)

const (
	log2CountMin = 9
	log2CountMax = 24
	CountMaxDefault =  uint64(1) << uint64(log2CountMax)
)

func FileCountCases(countMax uint64) []Case {
	var cases []Case
	for i := log2CountMin; i <= log2CountMax; i++ {
		count := uint64(1) << uint64(i)
		if count > countMax {
			break
		}
		cases = append(cases, FileCountCase("prefix", count))
	}
	return cases
}

func FileCountCase(prefix string, count uint64) Case {
	title := fmt.Sprintf("create %d files under single prefix", count)

	var contentLength int64 = DefaultContentLengthBytes
	body := make([]byte, contentLength)
	rand.Read(body)
	bodyProvider := func() io.Reader {
		return bytes.NewReader(body)
	}

	execution := func(target objects.Target) (ok bool, detail string) {
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

	return newCase(title, execution)
}
