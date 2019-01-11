package objects

import (
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/ncw/swift"

	"github.com/dmolesUC3/cos/internal/logging"
	"github.com/dmolesUC3/cos/internal/protocols"
)

const defaultRetries = 3

// SwiftObject is an OpenStack Swift implementation of Object
type SwiftObject struct {
	container       string
	objectName      string
	cnxParams       protocols.SwiftConnectionParams
	logger          logging.Logger
	swiftConnection *swift.Connection
}

func (obj *SwiftObject) String() string {
	return fmt.Sprintf("SwiftObject { container: '%v', objectName: '%v', cnxParams: %v, logger: %v, swiftConnection: %v",
		obj.container,
		obj.objectName,
		obj.cnxParams,
		obj.logger,
		obj.swiftConnection,
	)
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

// StreamDown streams the object down in ranged requests of the specified size, passing
// each byte range retrieved to the specified handler function, in sequence.
func (obj *SwiftObject) StreamDown(rangeSize int64, handleBytes func([]byte) error) (int64, error) {

	logger := obj.logger
	totalBytes := int64(0)

	cnx, err := obj.connection()
	if err != nil {
		return 0, err
	}

	file, _, err := cnx.ObjectOpen(obj.container, obj.objectName, false, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		logger.Detailf("read %d total bytes\n", totalBytes)
		err = file.Close()
	}()

	contentLength, err := file.Length()
	if err != nil {
		return 0, err
	}

	nsStart := time.Now().UnixNano()
	nsLastUpdate := nsStart
	for totalBytes < contentLength {
		byteRange := make([]byte, rangeSize)
		actualBytes, err := file.Read(byteRange)
		eof := err == io.EOF
		totalBytes = totalBytes + int64(actualBytes)

		nsNow := time.Now().UnixNano()
		nsSinceLastUpdate := nsNow - nsLastUpdate
		verbose := logger.Verbose()
		if verbose && (nsSinceLastUpdate > int64(time.Second)) || eof {
			nsLastUpdate = nsNow
			nsElapsed := nsNow - nsStart
			nsPerByte := float64(nsElapsed) / float64(totalBytes)
			estKps := float64(time.Second) / (float64(1024) * nsPerByte)
			nsRemaining := int64(float64(contentLength-totalBytes) * nsPerByte)
			logging.DetailProgress(logger, totalBytes, contentLength, estKps, nsElapsed, nsRemaining)
		}

		err = handleBytes(byteRange[0:actualBytes])
		if err != nil {
			return totalBytes, err
		}
		if eof {
			break
		}
	}
	return totalBytes, nil
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
