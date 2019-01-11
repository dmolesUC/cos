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

	nsLastUpdate := int64(0)
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
			sElapsed := nsElapsed / int64(time.Second)
			bps := float64(totalBytes) / float64(sElapsed)
			sRemaining := float64(contentLength - totalBytes) / bps
			logger.Detailf("read %d of %d bytes (%ds elapsed, %.0fs remaining\n", totalBytes, actualBytes, nsElapsed, sRemaining)
		}

		err = handleBytes(byteRange)
		if err != nil {
			return totalBytes, err
		}
		if eof {
			break
		}
	}

	//rangeCount := (contentLength + rangeSize - 1) / rangeSize
	//for rangeIndex := int64(0); rangeIndex < rangeCount; rangeIndex++ {
	//	// byte ranges are 0-indexed and inclusive
	//	startInclusive := rangeIndex * rangeSize
	//	var endInclusive int64
	//	if (rangeIndex + 1) < rangeCount {
	//		endInclusive = startInclusive + rangeSize - 1
	//	} else {
	//		endInclusive = contentLength - 1
	//	}
	//	expectedBytes := (endInclusive + 1) - startInclusive
	//	logger.Detailf("range %d of %d: retrieving %d bytes (%d - %d)\n", rangeIndex, rangeCount, expectedBytes, startInclusive, endInclusive)
	//
	//	byteRange := make([]byte, expectedBytes)
	//	actualBytes, err := file.Read(byteRange)
	//	actualBytes64 := int64(actualBytes)
	//	totalBytes = totalBytes + actualBytes64
	//	if err != nil {
	//		return totalBytes, err
	//	}
	//	if actualBytes64 != expectedBytes {
	//		return totalBytes, fmt.Errorf("range %d of %d: expected %d bytes (%d - %d), got %d\n", rangeIndex, rangeCount, expectedBytes, startInclusive, endInclusive, actualBytes)
	//	}
	//	err = handleBytes(byteRange)
	//	if err != nil {
	//		return totalBytes, err
	//	}
	//}
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
