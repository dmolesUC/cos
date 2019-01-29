# cos

A tool for checking cloud object storage.

## Running

### Create/retrieve/verify/delete with `cos crvd`

#### S3 (including Minio)

Amazon (Merritt Stage) example with implicit credentials:

```
cos crvd s3://uc3-s3mrt5001-stg/ -e https://s3-s3-us-west-2.amazonaws.com/
```

Minio example with explicit credentials:

```
AWS_ACCESS_KEY_ID=<access key> \
AWS_SECRET_ACCESS_KEY=<secret access key> \
cos crvd s3://mrt-test/ -e http://localhost:9000/
```

#### OpenStack/Swift

Note that for OpenStack Swift, the credentials must always be specified
explicitly with the SWIFT_API_USER and SWIFT_API_KEY environment variables, and
the bucket URL must be in `swift://<container>/` form, with an explicit
`--endpoint` parameter.

```
SWIFT_API_USER=<user> \
SWIFT_API_KEY=<key> \
cos crvd -v swift://distrib.stage.9001.__c5e/ -e http://cloud.sdsc.edu/auth/v1.0 
```

### Fixity checking with `cos check`

#### S3 (including Minio)

Amazon (Merritt Stage) example with implicit credentials:

```
cos check -v \
  's3://uc3-s3mrt5001-stg/ark:/99999/fk46w9nc06|1|producer/1500MBTestObject.blob' \
  -e 'https://s3-us-west-2.amazonaws.com/' \
  -x d0487cf92819b6f70a4769419348ab51ed77c519664a6262283e0016b9a6235c
```

```
cos check -v \
  's3://uc3-s3mrt5001-prd/ark:/13030/qt30c9r5qj|1|producer/content/supp/FreeSolv_paper.tar.gz' \
  -e 'https://s3-us-west-2.amazonaws.com/' \
  -x c0916ef45d917578e4dcdc3045d9340738d0e750c0ab9f6a8e866aa28da677df
```

(S3 keys for Merritt objects are of the form `<ark>|<version>|<pathname>`.) 

Minio example with explicit credentials:

```
AWS_ACCESS_KEY_ID=<access key> \
AWS_SECRET_ACCESS_KEY=<secret access key> \
cos check http://127.0.0.1:9000/mrt-test/inusitatum.png -a md5 -x cadf871cd4135212419f488f42c62482`
```

Amazon example with explicit credentials:

```
AWS_ACCESS_KEY_ID=<access key> \
AWS_SECRET_ACCESS_KEY=<secret access key> \
cos check https://s3-us-west-2.amazonaws.com/www.dmoles.net/images/fa/archive.svg
```

#### OpenStack/Swift

Note that for OpenStack Swift, the credentials must always be specified
explicitly with the SWIFT_API_USER and SWIFT_API_KEY environment variables, and
the object URL must be in `swift://<container>/<name>` form, with an explicit
`--endpoint` parameter.

```
SWIFT_API_USER=<user> \
SWIFT_API_KEY=<key> \
cos check -v \
  -e 'http://cloud.sdsc.edu/auth/v1.0' \
  'swift://distrib.stage.9001.__c5e/ark:/99999/fk4kw5kc1z|1|producer/6GBZeroFile.txt'
```

## For developers

`cos` is a [Go 1.11 module](https://github.com/golang/go/wiki/Modules). 

As such, it requires Go 1.11 or later, and should be cloned _outside_
`$GOPATH/src`.

### Building

From the project root:

- to build `cos`, writing the executable to the source directory, use `go build`.
- to build `cos` and install it in `$GOPATH/bin`, use `go install`.

### Running tests

To run all tests in all subpackages, from the project root, use `go test ./...`.

To run all tests in all subpackages with coverage and view a coverage report, use

```
go test -coverprofile=coverage.out ./... \
&& go tool cover -html=coverage.out
```

#### Cross-compiling

To cross-compile for Linux (Intel, 64-bit):

```
GOOS=linux GOARCH=amd64 go build -o <output file>
```

### Configuring JetBrains IDEs (GoLand or IDEA)

In **Preferences > Go > Go Modules (vgo)** (GoLand) or **Preferences >
Languages & Frameworks Go > Go Modules (vgo)** (IDEA + Go plugin) , check
“Enable Go Modules (vgo) integration“. The “Vgo Executable” field should
default to “Project SDK” (1.11.x).

## Roadmap

- ✅ fixity checking: expected vs. actual
- ✅ sanity check: can we create/retrieve/verify/delete a file?
- 🔲 weird filenames
- 🔲 scalability
  - large files
  - large numbers of files per bucket
  - large numbers of files per key prefix
- 🔲 streaming download performance
- 🔲 reliability
