package objects

import (
	"fmt"
	"net/url"

	"github.com/ncw/swift"

	"github.com/dmolesUC3/cos/internal/logging"
	"github.com/dmolesUC3/cos/internal/protocols"
)

// SwiftObject is an OpenStack Swift implementation of Object
type SwiftObject struct {
	container       string
	objectName      string
	cnxParams       protocols.SwiftConnectionParams
	logger          logging.Logger
	swiftConnection *swift.Connection
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
	cnx, err := obj.connection()
	if err != nil {
		return 0, err
	}
	file, _, err := cnx.ObjectOpen(obj.container, obj.objectName, false, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		err = file.Close()
	}()

	contentLength, err := file.Length()
	if err != nil {
		return 0, err
	}

	logger := obj.logger
	totalBytes := int64(0)
	rangeCount := (contentLength + rangeSize - 1) / rangeSize
	for rangeIndex := int64(0); rangeIndex < rangeCount; rangeIndex++ {
		// byte ranges are 0-indexed and inclusive
		startInclusive := rangeIndex * rangeSize
		var endInclusive int64
		if (rangeIndex + 1) < rangeCount {
			endInclusive = startInclusive + rangeSize - 1
		} else {
			endInclusive = contentLength - 1
		}
		expectedBytes := (endInclusive + 1) - startInclusive
		logger.Detailf("range %d of %d: retrieving %d bytes (%d - %d)\n", rangeIndex, rangeCount, expectedBytes, startInclusive, endInclusive)

		byteRange := make([]byte, expectedBytes)
		actualBytes, err := file.Read(byteRange)
		expectedBytes64 := int64(actualBytes)
		totalBytes = totalBytes + expectedBytes64
		if err != nil {
			return totalBytes, err
		}
		if expectedBytes64 != expectedBytes {
			logger.Infof("range %d of %d: expected %d bytes (%d - %d), got %d\n", rangeIndex, rangeCount, expectedBytes, startInclusive, endInclusive, actualBytes)
		}
		err = handleBytes(byteRange)
		if err != nil {
			return totalBytes, err
		}
	}
	return totalBytes, nil
}

// ------------------------------------------------------------
// Unexported functions

func (obj *SwiftObject) connection() (*swift.Connection, error) {
	return nil, fmt.Errorf("connection() not implemented")
}
