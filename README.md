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
