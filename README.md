# cos

A tool for testing and validating cloud object storage.

## Invocation

Invocation is in the form

```
cos <command> [flags] [URL]
```

where `<command>` is one of:

- [`check`](https://github.com/dmolesUC3/cos#cos-check): 
  compute and (optionally) verify the digest of an object
- [`crvd`](https://github.com/dmolesUC3/cos#cos-crvd): 
  create, retrieve, verify, and delete an object
- [`keys`](https://github.com/dmolesUC3/cos#cos-keys): 
  test the keys supported by an object storage endpoint
- `help`: 
  list these commands, or get help for a subcommand

and `[URL]` can be the URL of an object or of a bucket/container, depending
on the context. The protocol (`s3://` or `swift://`) of the URL is used to
determine the cloud storage API to use.

### Flags

All `cos` commands support the following flags:

| Flag | Short form | Description |
| :-- | :-- | :-- |
| `--endpoint` | `-e` | HTTP(S) endpoint URL (required) |
| `--region` | `-r` | AWS region (optional) |
| `--verbose` | `-v` | Verbose output |
| `--help` | `-h` | Print help and exit |

For Amazon S3 buckets, the region can usually be determined from the
endpoint URL. If not, and if the `--region` flag is not provided, it
defaults to `us-west-2`.

For OpenStack Swift containers, the `--region` flag is ignored.

Additional command-specific flags are listed below.

> #### TODO: document authentication for both Swift and S3

## Commands

### `cos check`

The `check` command computes and (optionally) verifies the digest of an
object. The object is streamed in 5-MiB chunks, each chunk being added to
the digest computation and then discarded, thus making it possible to 
verify objects of arbitary size, not limited by local storage space.

In addition to the flags listed above, the `check` command supports the following:

| Flag | Short form | Description |
| :-- | :-- | :-- |
| `--algorithm` | `-a` | Digest algorithm (md5 or sha256; defaults to sha256) |
| `--expected` | `-x` | Expected digest value |

By default, `check` outputs the digest to standard output, and exits:

```
$ cos check --endpoint https://s3.us-west-2.amazonaws.com/ s3://www.dmoles.net/images/fa/archive.svg/
c99ad299fa53d5d9688909164cf25b386b33bea8d4247310d80f615be29978f5
```

If given an expected value that does not match, prints a message to standard
error, and exits with a nonzero (unsuccessful) exit code.

```
$ cos check --endpoint https://s3.us-west-2.amazonaws.com/ s3://www.dmoles.net/images/fa/archive.svg/ \
  -x 5f87992eb516f08d0137424d8aeb33b683b52fc4619098869d5d35af992da99c
digest mismatch: 
expected: 5f87992eb516f08d0137424d8aeb33b683b52fc4619098869d5d35af992da99c
actual: c99ad299fa53d5d9688909164cf25b386b33bea8d4247310d80f615be29978f5
```

### `cos crvd`

> #### TODO: document this

### `cos keys`

> #### TODO: document this

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

### Running tests

To run all tests in all subpackages, from the project root, use `go test ./...`.

To run all tests in all subpackages with coverage and view a coverage report, use

```
go test -coverprofile=coverage.out ./... \
&& go tool cover -html=coverage.out
```

### Configuring JetBrains IDEs (GoLand or IDEA)

In **Preferences > Go > Go Modules (vgo)** (GoLand) or **Preferences >
Languages & Frameworks Go > Go Modules (vgo)** (IDEA + Go plugin) , check
“Enable Go Modules (vgo) integration“. The “Vgo Executable” field should
default to “Project SDK” (1.11.x).

## Roadmap

- ✅ fixity checking: expected vs. actual
- ✅ sanity check: can we create/retrieve/verify/delete a file?
- ✅ weird filenames
- 🔲 scalability
  - large files
  - large numbers of files per bucket
  - large numbers of files per key prefix
- 🔲 streaming download performance
- 🔲 reliability 
