package objects

import (
	"fmt"
	"io"
	"net/url"
	"time"

	"github.com/ncw/swift"

	"github.com/dmolesUC3/cos/internal/logging"
	"github.com/dmolesUC3/cos/internal/protocols"
	"github.com/dmolesUC3/cos/internal/streaming"
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

func (obj *SwiftObject) StreamDown(rangeSize int64, handleBytes func([]byte) error) (int64, error) {
	cnx, err := obj.connection()
	if err != nil {
		return 0, err
	}

	// this will 404 if the object doesn't exist
	contentLength, err := obj.ContentLength()
	if err != nil {
		return 0, err
	}

	fillRange := func(byteRange *streaming.ByteRange) (int64, error) {
		startInclusive := byteRange.StartInclusive
		endInclusive := byteRange.EndInclusive
		rangeStr := fmt.Sprintf("bytes=%d-%d", startInclusive, endInclusive)

		headers := map[string]string { "Range" : rangeStr }

		file, _, err := cnx.ObjectOpen(obj.container, obj.objectName, false, headers)
		if err != nil {
			return 0, err
		}
		bytesRead, err := io.ReadFull(file, byteRange.Buffer)
		return int64(bytesRead), err
	}

	streamer, err := streaming.NewStreamer(rangeSize, contentLength, &fillRange)
	if err != nil {
		return 0, err
	}

	return streamer.StreamDown(obj.logger, handleBytes)
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
