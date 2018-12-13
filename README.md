# coscheck

A tool for checking cloud object storage.

## Roadmap

- fixity checking: expected vs. actual
  - [SHA-256](https://golang.org/pkg/crypto/sha256/)
  - [MD5](https://golang.org/pkg/crypto/md5/)
- streaming download performance
  - time download to nowhere 
  - time download to file
  - include fixity check
