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
	nsStart := time.Now().UnixNano()

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
		logger.Detailf("read %d total bytes", totalBytes)
		err = file.Close()
	}()

	contentLength, err := file.Length()
	if err != nil {
		return 0, err
	}

	nsLastUpdate := int64(time.Second) // skip first update
	for totalBytes < contentLength {
		byteRange := make([]byte, rangeSize)
		actualBytes, err := file.Read(byteRange)
		eof := err == io.EOF
		totalBytes = totalBytes + int64(actualBytes)

		nsNow := time.Now().UnixNano()
		nsSinceLastUpdate := nsNow - nsLastUpdate
		if nsSinceLastUpdate > int64(time.Second) || eof {
			nsLastUpdate = nsNow
			nsElapsed := nsNow - nsStart
			nsPerByte := float64(nsElapsed) / float64(totalBytes)
			nsRemaining := int64(float64(contentLength - totalBytes) * nsPerByte)
			logProgress(logger, totalBytes, actualBytes, nsElapsed, nsRemaining)
		}

		err = handleBytes(byteRange)
		if err != nil {
			return totalBytes, err
		}
		if eof {
			break
		}
	}
	return totalBytes, nil
}

func logProgress(logger logging.Logger, totalBytes int64, actualBytes int, nsElapsed int64, nsRemaining int64) {
	elapsedStr := formatNanos(nsElapsed)
	remainingStr := formatNanos(nsRemaining)
	logger.Detailf("read %d of %d bytes (%v elapsed, %v remaining)\n", totalBytes, actualBytes, elapsedStr, remainingStr)
}

func formatNanos(ns int64) string {
	hours := ns / int64(time.Hour)
	remainder := ns % int64(time.Hour)
	minutes := remainder / int64(time.Minute)
	remainder = ns % int64(time.Minute)
	seconds := remainder / int64(time.Second)
	if hours > 0 {
		return fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds)
	}
	if minutes > 0 {
		return fmt.Sprintf("%dm %ds", minutes, seconds)
	}
	return fmt.Sprintf("%ds", seconds)
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
