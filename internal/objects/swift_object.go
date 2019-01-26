package objects

import (
	"fmt"
	"io"
	"net/url"

	"code.cloudfoundry.org/bytefmt"
	"github.com/ncw/swift"

	"github.com/dmolesUC3/cos/internal/logging"
	"github.com/dmolesUC3/cos/internal/protocols"
	"github.com/dmolesUC3/cos/internal/streaming"
)

const (
	defaultRetries   = 3
	dloSizeThreshold = int64(2 * bytefmt.GIGABYTE)
	//dloSegmentSize   = int64(bytefmt.GIGABYTE) // TODO: any way to set this?
)

// SwiftObject is an OpenStack Swift implementation of Object
type SwiftObject struct {
	container       string
	objectName      string
	cnxParams       protocols.SwiftConnectionParams
	logger          logging.Logger
	swiftConnection *swift.Connection
}

func (obj *SwiftObject) Protocol() string {
	return protocolSwift
}

func (obj *SwiftObject) Pretty() string {
	format := `SwiftObject { 
		container:      '%v' 
		objectName:     '%v' 
		cnxParams:       %v 
		logger:          %v 
		swiftConnection: %v
	}`
	format = logging.Untabify(format, " ")
	args := logging.Prettify(obj.container, obj.objectName, obj.cnxParams, obj.logger, obj.swiftConnection)
	return fmt.Sprintf(format, args...)
}

func (obj *SwiftObject) Refresh() {
	obj.swiftConnection = nil
}

func (obj *SwiftObject) String() string {
	return fmt.Sprintf("SwiftObject { container: '%v', objectName: '%v', cnxParams: %v, logger: %v, swiftConnection: %v }",
		obj.container,
		obj.objectName,
		obj.cnxParams,
		obj.logger,
		obj.swiftConnection,
	)
}

func (obj *SwiftObject) Logger() logging.Logger {
	return obj.logger
}

// Endpoint returns the Swift authentication URL
func (obj *SwiftObject) Endpoint() *url.URL {
	return obj.cnxParams.AuthURL
}

// Bucket returns the Swift container
func (obj *SwiftObject) Bucket() *string {
	if obj.container == "" {
		return nil
	}
	return &obj.container
}

// Key returns the Swift object name
func (obj *SwiftObject) Key() *string {
	if obj.objectName == "" {
		return nil
	}
	return &obj.objectName
}

func (obj *SwiftObject) ContentLength() (int64, error) {
	cnx, err := obj.connection()
	if err != nil {
		return 0, err
	}
	info, _, err := cnx.Object(obj.container, obj.objectName)
	if err != nil {
		return 0, err
	}
	return info.Bytes, nil
}

func (obj *SwiftObject) DownloadRange(startInclusive, endInclusive int64, buffer []byte) (int64, error) {
	cnx, err := obj.connection()
	if err != nil {
		return 0, err
	}
	rangeStr := fmt.Sprintf("bytes=%d-%d", startInclusive, endInclusive)
	headers := map[string]string{"Range": rangeStr}
	file, _, err := cnx.ObjectOpen(obj.container, obj.objectName, false, headers)
	if err != nil {
		return 0, err
	}
	err = streaming.ReadExactly(file, buffer)
	if err != nil {
		return 0, err
	}
	return int64(len(buffer)), nil
}

func (obj *SwiftObject) Create(body io.Reader, length int64) error {
	cnx, err := obj.connection()
	if err != nil {
		return err
	}
	logger := obj.logger

	// TODO: allow object to include an expected MD5
	var out io.WriteCloser
	if length <= dloSizeThreshold { // 2 GiB
		out, err = cnx.ObjectCreate(obj.container, obj.objectName, false, "", "", nil)
	} else {
		logger.Detailf(
			"object size %d is greater than single-object maximum %d; creating dynamic large object\n",
			length, dloSizeThreshold,
		)
		dloOpts := swift.LargeObjectOpts{
			Container:  obj.container,
			ObjectName: obj.objectName,
			ChunkSize:  streaming.DefaultRangeSize, // 5 MiB
		}
		out, err = cnx.DynamicLargeObjectCreateFile(&dloOpts)
	}
	if err != nil {
		logger.Detailf("error opening upload stream: %v\n", err)
		return err
	}

	defer func() {
		err := out.Close()
		if err != nil {
			logger.Detailf("error closing upload stream: %v\n", err)
		}
	}()

	buffer := make([]byte, streaming.DefaultRangeSize)
	written, err := io.CopyBuffer(out, body, buffer)
	if err != nil {
		logger.Detailf("error writing to upload stream: %v\n", err)
	}
	logger.Detailf("wrote %d bytes to %v/%v\n", written, obj.container, obj.objectName)
	return err
}

func (obj *SwiftObject) Delete() (err error) {
	cnx, err := obj.connection()
	if err != nil {
		return err
	}
	// TODO: detect DynamicLargeObjects
	return cnx.ObjectDelete(obj.container, obj.objectName)
}

// ------------------------------------------------------------
// Unexported functions

func (obj *SwiftObject) connection() (*swift.Connection, error) {
	cnxParams := obj.cnxParams
	authUrl := cnxParams.AuthURL
	if authUrl == nil {
		return nil, fmt.Errorf("authUrl not set in connection parameters: %v", cnxParams)
	}
	authUrlStr := authUrl.String()
	cnx := swift.Connection{
		UserName: cnxParams.UserName,
		ApiKey:   cnxParams.APIKey,
		AuthUrl:  authUrlStr,
		Retries:  defaultRetries,
	}
	return &cnx, nil
}
