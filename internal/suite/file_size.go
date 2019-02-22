package suite

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"unicode"

	"github.com/dmolesUC3/cos/internal/objects"

	. "code.cloudfoundry.org/bytefmt"

	"github.com/dmolesUC3/cos/internal/logging"
	. "github.com/dmolesUC3/cos/pkg"
)

const (
	SizeMaxDefault = 256 * GIGABYTE
)

func FileSizeCases(sizeMax int64) []Case {
	tasks := []Case{FileSizeCase(0)}
	for _, unit := range []int64{BYTE, KILOBYTE, MEGABYTE, GIGABYTE, TERABYTE} {
		if unit > sizeMax {
			break
		}
		for _, multiplier := range []int64{1, 16, 256} {
			size := multiplier * unit
			if size > sizeMax {
				break
			}
			tasks = append(tasks, FileSizeCase(size))
		}
	}
	return tasks
}

func FileSizeCase(size int64) Case {
	title := fmt.Sprintf("create/retrieve/verify/delete %v file", logging.FormatBytes(size))
	execution := func(target objects.Target) (ok bool, detail string) {
		crvd := NewCrvd(target, "", size, DefaultRandomSeed)
		err := crvd.CreateRetrieveVerifyDelete()
		if err == nil {
			return true, ""
		} else {
			return false, err.Error()
		}
	}
	return newCase(title, execution)
}

func ParseSizeMax(sizeStr string) (int64, error) {
	sizeIsNumeric := strings.IndexFunc(sizeStr, unicode.IsLetter) == -1
	if sizeIsNumeric {
		return strconv.ParseInt(sizeStr, 10, 64)
	}

	bytes, err := ToBytes(sizeStr)
	if err == nil && bytes > math.MaxInt64 {
		return 0, fmt.Errorf("specified size %d bytes exceeds maximum %d", bytes, math.MaxInt64)
	}
	return int64(bytes), err
}