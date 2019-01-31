package objects

import (
	"fmt"
	"io"

	"code.cloudfoundry.org/bytefmt"
	"github.com/ncw/swift"

	"github.com/dmolesUC3/cos/internal/logging"
	"github.com/dmolesUC3/cos/internal/streaming"
)

const (
	dloSizeThreshold = int64(2 * bytefmt.GIGABYTE)
)

type SwiftObject struct {
	Endpoint  *SwiftTarget
	Container string
	Name      string
}

// ------------------------------
// Object implementation

func (obj *SwiftObject) Pretty() string {
	return fmt.Sprintf("swift://%v/%v", obj.Container, obj.Name)
}

func (obj *SwiftObject) String() string {
	return obj.Pretty()
}

func (obj *SwiftObject) GetEndpoint() Target {
	return obj.Endpoint
}

func (obj *SwiftObject) ContentLength() (length int64, err error) {
	cnx, err := obj.Endpoint.Connection()
	if err != nil {
		return 0, err
	}
	info, _, err := cnx.Object(obj.Container, obj.Name)
	if err != nil {
		return 0, err
	}
	return info.Bytes, nil
}

func (obj *SwiftObject) DownloadRange(startInclusive, endInclusive int64, buffer []byte) (n int64, err error) {
	cnx, err := obj.Endpoint.Connection()
	if err != nil {
		return 0, err
	}
	rangeStr := fmt.Sprintf("bytes=%d-%d", startInclusive, endInclusive)
	headers := map[string]string{"Range": rangeStr}
	file, _, err := cnx.ObjectOpen(obj.Container, obj.Name, false, headers)
	if err != nil {
		return 0, err
	}
	err = streaming.ReadExactly(file, buffer)
	if err != nil {
		return 0, err
	}
	return int64(len(buffer)), nil
}

func (obj *SwiftObject) Create(body io.Reader, length int64) (err error) {
	cnx, err := obj.Endpoint.Connection()
	if err != nil {
		return err
	}

	logger := logging.DefaultLogger()
	var out io.WriteCloser
	if length <= dloSizeThreshold { // 2 GiB
		out, err = cnx.ObjectCreate(obj.Container, obj.Name, false, "", "", nil)
	} else {
		logger.Tracef(
			"Object size %d is greater than single-object maximum %d; creating dynamic large object\n",
			length, dloSizeThreshold,
		)
		dloOpts := swift.LargeObjectOpts{
			Container:  obj.Container,
			ObjectName: obj.Name,
			ChunkSize:  streaming.DefaultRangeSize, // 5 MiB
		}
		out, err = cnx.DynamicLargeObjectCreateFile(&dloOpts)
	}
	if err != nil {
		logger.Tracef("Error opening upload stream: %v\n", err)
		return err
	}

	defer func() {
		err := out.Close()
		if err != nil {
			logger.Tracef("Error closing upload stream: %v\n", err)
		}
	}()

	buffer := make([]byte, streaming.DefaultRangeSize)
	written, err := io.CopyBuffer(out, body, buffer)
	if err != nil {
		logger.Tracef("Error writing to upload stream: %v\n", err)
	}
	logger.Tracef("Wrote %d bytes to %v\n", written, obj)
	return err
}

func (obj *SwiftObject) Delete() (err error) {
	cnx, err := obj.Endpoint.Connection()
	if err != nil {
		return err
	}

	// TODO: detect DynamicLargeObjects
	err = cnx.ObjectDelete(obj.Container, obj.Name)
	logger := logging.DefaultLogger()
	logger.Detailf("Deleting %v\n", obj)
	if err == nil {
		logger.Detailf("Deleted %v\n", obj)
	} else {
		logger.Detailf("Deleting %v failed: %v", obj, err)
	}
	return err

}

