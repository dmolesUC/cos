# cos

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

`cos` is a [Go 1.11 module](https://github.com/golang/go/wiki/Modules). 

As such, it requires Go 1.11 or later, and should be cloned _outside_
`$GOPATH/src`.

### Building

From the project root:

- to build `cos`, writing the executable to the source directory, use `go build`.
- to build `cos` and install it in `$GOPATH/bin`, use `go install`.

### Running

Amazon example with credentials:

```
AWS_ACCESS_KEY_ID=<access key> \
AWS_SECRET_ACCESS_KEY=<secret access key> \
cos check https://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive.svg
```

Minio example with credentials:

```
AWS_ACCESS_KEY_ID=<access key> \
AWS_SECRET_ACCESS_KEY=<secret access key> \
cos check http://127.0.0.1:9000/mrt-test/inusitatum.png -a md5 -x cadf871cd4135212419f488f42c62482`
```

### Configuring JetBrains GoLand IDE

In **Preferences > Go > Go Modules (vgo)**, check “Enable Go Modules (vgo)
integration“. The “Vgo Executable” field should default to “Project SDK”
(1.11.x).

