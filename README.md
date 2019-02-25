# cos

A tool for testing and validating cloud object storage.

- [Invocation](#invocation)
- [Authentication](#authentication)
- [Commands](#commands)
   - [cos check](#cos-check)
   - [cos crvd](#cos-crvd)
   - [cos keys](#cos-keys)
   - [cos suite](#cos-suite)
- [For developers](#for-developers)
   - [Building](#building)
   - [Running tests](#running-tests)
   - [Configuring JetBrains IDEs (GoLand or IDEA)](#configuring-jetbrains-ides-goland-or-idea)
- [Roadmap](#roadmap)

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
- [`suite`](https://github.com/dmolesUC3/cos#cos-suite): 
  run a suite of test cases investigating various possible limitations of a
  cloud storage service
- `help`: 
  list these commands, or get help for a subcommand

and `[URL]` can be the URL of an object or of a bucket/container, depending
on the context. The protocol (`s3://` or `swift://`) of the URL is used to
determine the cloud storage API to use.

## Authentication

For authentication, `cos` uses the same environment variables as the [AWS
CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html)
(for S3 and compatible storage) or [OpenStack Swift
CLI](https://docs.openstack.org/python-swiftclient/latest/cli/index.html)
(for Swift storage):

| Protocol | Variable                | Purpose                                       |
| :---     | :---                    | :---                                          |
| S3       | `AWS_ACCESS_KEY_ID`     | AWS access key                                |
|          | `AWS_SECRET_ACCESS_KEY` | Secret key associated with the AWS access key |
| Swift    | `ST_USER`               | Swift username                                |
|          | `ST_KEY`                | Swift password                                |

Credentials for S3 storage can also be specified [in various other
ways](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials)
supported by the AWS SDK for Go, such as a shared credentials file or, when
running in the Amazon EC2 environment, an IAM role.

### Flags

All `cos` commands support the following flags:

| Flag                  | Short form | Description                     |
| :---                  | :---       | :---                            |
| `--endpoint ENDPOINT` | `-e`       | HTTP(S) endpoint URL (required) |
| `--region REGION`     | `-r`       | AWS region (optional)           |
| `--verbose`           | `-v`       | Verbose output                  |
| `--help`              | `-h`       | Print help and exit             |

For Amazon S3 buckets, the region can usually be determined from the
endpoint URL. If not, and if the `--region` flag is not provided, it
defaults to `us-west-2`.

For OpenStack Swift containers, the `--region` flag is ignored.

Additional command-specific flags are listed below.

## Commands

### `cos check`

The `check` command computes and (optionally) verifies the digest of an
object. The object is streamed in five-megabyte chunks, each chunk being
added to the digest computation and then discarded, thus making it possible
to verify objects of arbitary size, not limited by local storage space.

In addition to the global flags listed above, the `check` command supports the following:

| Flag                | Short form | Description                                          |
| :---                | :---       | :---                                                 |
| `--algorithm ALG`   | `-a`       | Digest algorithm (md5 or sha256; defaults to sha256) |
| `--expected DIGEST` | `-x`       | Expected digest value                                |

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

The `crvd` command creates, retrieves, verifies, and deletes an object.
The object consists of a stream of random bytes of the specified size.

The size may be specified as an exact number of bytes, or using human-readable
quantities such as "5K" (4 KiB or 4096 bytes), "3.5M" (3.5 MiB or 3670016 bytes),
etc. The units supported are bytes (B), binary kilobytes (K, KB, KiB), 
binary megabytes (M, MB, MiB), binary gigabytes (G, GB, GiB), and binary 
terabytes (T, TB, TiB). If no unit is specified, bytes are assumed.

Random bytes are generated using the Go default random number generator, with
a default seed of 0, for repeatability. An alternative seed can be specified
with the `--random-seed` flag.

In addition to the global flags listed above, the `check` command supports the following:

| Flag                 | Short form | Description                                          |
| :---                 | :---       | :---                                                 |
| `--size SIZE`        | `-s`       | size of object to create (default 128 bytes)         |
| `--key KEY`          | `-k`       | key to create (defaults to `cos-crvd-TIMESTAMP.bin`) |
| `--random-seed SEED` |            | seed for random-number generator (default 1)         |
| `--keep`             |            | keep object after verification (default false)       |

```
$ crvd swift://distrib.stage.9001.__c5e/ -e http://cloud.sdsc.edu/auth/v1.0 
128B object created, retrieved, verified, and deleted (swift://distrib.stage.9001.__c5e/cos-crvd-1549324512.bin)
```

### `cos keys`

The `keys` command tests the keys supported by an object storage endpoint,
creating, retrieving, validating, and deleting a small object for each value
in the specified key list. 

In addition to the global flags listed above, the `keys` command supports
the following:

| Flag             | Short form | Description                                    |
| :---             | :---       | :---                                           |
| `--raw`          |            | write keys in raw (unquoted) format            |
| `--ok FILE`      | `-o`       | write successful ("OK") keys to specified file |
| `--bad FILE`     | `-b`       | write failed ("bad") keys to specified file    |
| `--list LIST`    | `-l`       | use the specified 'standard' list of keys      |
| `--file FILE`    | `-f`       | read keys to be tested from the specified file |
| `--sample COUNT` | `-s`       | sample size, or 0 for all keys                 |


By default, `keys` outputs only failed keys, to standard output, writing
each key as a [quoted Go string literal](https://golang.org/pkg/strconv/#Quote).

```
$ cos keys s3://uc3-s3mrt5001-stg/ -e 'https://s3-us-west-2.amazonaws.com/' --list misc 
"../leading-double-dot-path"
"../../leading-multiple-double-dot-path"
"trailing-double-dot-path/.."
(...etc.)
```

Use the `--raw` option to write the keys without quoting or escaping; note
that this may produce confusing results if any of the keys contain
newlines.

```
$ cos keys s3://uc3-s3mrt5001-stg/ -e 'https://s3-us-west-2.amazonaws.com/' --list misc --raw
../leading-double-dot-path
../../leading-multiple-double-dot-path
trailing-double-dot-path/..
trailing-multiple-double-dot-path/../..
(...etc.)
```

Use the `--ok` option to write successful keys to a file, and the `--bad`
option (or shell redirection) to write failed keys to a file instead of
stdout.

```
$ cos keys s3://uc3-s3mrt5001-stg/ -e 'https://s3-us-west-2.amazonaws.com/' --list misc \
  --ok out/keys-ok.txt --bad out/keys-bad.txt
```

Several "standard" lists are provided (though these aren't very systematic;
see [#10](https://github.com/dmolesUC3/cos/issues/10)). Use the `--file`
option to specify a file containing keys to test, one key per file,
separated by newlines (LF, `U+000A`, `\n`).

```
$ cos keys s3://uc3-s3mrt5001-stg/ -e 'https://s3-us-west-2.amazonaws.com/' \
  --file my-keys.txt
```

Use the `--sample` option to check only a random sample from a large key list:

```
$ cos keys s3://uc3-s3mrt5001-stg/ -e 'https://s3-us-west-2.amazonaws.com/' \
  --file my-very-long-list-of-keys.txt \
  --sample 500
```

### `cos suite`

The `suite` command a suite of test cases investigating various possible limitations of a
cloud storage service:

- maximum file size (`--size`)
- maximum number of files per key prefix (`--count`)
- Unicode key support (`--unicode`)

If none of `--size`, `--count`, etc. is specified, all test cases are run.

In addition to the global flags listed above, the `keys` command supports
the following:

| Flag                | Short form | Description                                                          |
| :---                | :---       | :---                                                                 |
| `--size`            | `-s`       | test file sizes                                                      |
| `--size-max SIZE`   |            | max file size to create (default "256G")                             |
| `--count`           | `-c`       | test file counts                                                     |
| `--count-max COUNT` |            | max number of files to create, or -1 for no limit (default 16777216) |
| `--unicode`         | `-u`       | test Unicode keys                                                    |

The maximum size may be specified as an exact number of bytes, or using
human-readable quantities such as "5K" (4 KiB or 4096 bytes), "3.5M" (3.5
MiB or 3670016 bytes), etc. The units supported are bytes (B), binary
kilobytes (K, KB, KiB), binary megabytes (M, MB, MiB), binary gigabytes (G,
GB, GiB), and binary terabytes (T, TB, TiB). If no unit is specified, bytes
are assumed.

## For developers

`cos` is a [Go 1.11 module](https://github.com/golang/go/wiki/Modules). 

As such, it requires Go 1.11 or later, and should be cloned _outside_
`$GOPATH/src`.

### Building

The `cos` project can be built and installed simply with `go build` and `go
install`, but it also supports [Mage](https://magefile.org).

To install the latest version of Mage:

1. visit their [releases page](https://github.com/magefile/mage/releases),
   download the appropriate binary, and place it in your `$PATH`, or
2. from _outside_ this project directory (`go get` behaves differently when
   run in the context of a module project), execute the following:

   ```
   go get -u -d github.com/magefile/mage \
   && cd $GOPATH/src/github.com/magefile/mage \
   && go run bootstrap.go
   ```

#### Mage tasks:

| Tasks        | Purpose                                                          |
| :---         | :---                                                             |
| `build`      | builds a binary for the current platform                         |
| `buildAll`   | builds a binary for each target platform                         |
| `buildLinux` | builds a linux-amd64 binary (the most common cross-compile case) |
| `clean`      | removes compiled binaries from the current working directory     |
| `install`    | installs in $GOPATH/bin                                          |
| `platforms`  | lists target platforms for buildAll                              |

Note that `mage build` is a thin wrapper around `go build` and supports the
same environment variables, e.g. `$GOOS` and `$GOARCH`.

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
