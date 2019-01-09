# cos

A tool for checking cloud object storage.

## Roadmap

- ✅ fixity checking: expected vs. actual
- 🔲 streaming download performance
  - throughput
  - time download to nowhere 
  - time download to file
  - include fixity check
- 🔲 reliability
  - same file
  - different files
  - retries

## Running

### Fixity checking with `cos check`

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

## For developers

`cos` is a [Go 1.11 module](https://github.com/golang/go/wiki/Modules). 

As such, it requires Go 1.11 or later, and should be cloned _outside_
`$GOPATH/src`.

### Building

From the project root:

- to build `cos`, writing the executable to the source directory, use `go build`.
- to build `cos` and install it in `$GOPATH/bin`, use `go install`.

#### Cross-compiling

To cross-compile for Linux (Intel, 64-bit):

```
GOOS=linux GOARCH=amd64 go build -o <output file>
```

### Configuring JetBrains GoLand IDE

In **Preferences > Go > Go Modules (vgo)**, check “Enable Go Modules (vgo)
integration“. The “Vgo Executable” field should default to “Project SDK”
(1.11.x).

