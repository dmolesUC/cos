# coscheck

A tool for checking cloud object storage.

## Roadmap

- fixity checking: expected vs. actual
  - [SHA-256](https://golang.org/pkg/crypto/sha256/)
  - [MD5](https://golang.org/pkg/crypto/md5/)
  - notes:
    - use [WriteAtBuffer](https://docs.aws.amazon.com/sdk-for-go/api/aws/#WriteAtBuffer)
      to deal w/S3 API's need for a random-access writer?
    - use [GetObjectInput.Range](https://docs.aws.amazon.com/sdk-for-go/api/service/s3/#GetObjectInput)
      to limit the size of the buffer?
- streaming download performance
  - throughput
  - time download to nowhere 
  - time download to file
  - include fixity check
- reliability
  - same file
  - different files
  - retries

## For developers

`coscheck` is a [Go 1.11 module](https://github.com/golang/go/wiki/Modules). 

As such, it requires Go 1.11 or later, and should be cloned _outside_
`$GOPATH/src`.

### Building

From the project root:

- to build `coscheck`, writing the executable to the source directory, use `go build`.
- to build `coscheck` and install it in `$GOPATH/bin`, use `go install`.

### Configuring JetBrains GoLand IDE

In **Preferences > Go > Go Modules (vgo)**, check “Enable Go Modules (vgo)
integration“. The “Vgo Executable” field should default to “Project SDK”
(1.11.x).

