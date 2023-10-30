# Python SpiceDB Validation Library

This is forked from Authzed and maintains the same Apache 2.0 license as `zed` and SpiceDB.

## Building

```shell
go build -buildmode=c-shared -o src/spicedb_validation/dll/spicedb_validation.so main.go
```

```shell
pip wheel -w whl .
```

Check the package content:

```shell
unzip -l whl/spicedb_validation-0.0.1-py3-none-any.whl
```

## Attribution

Original code Copyright Authzed, Inc. 2023

New code Copyright Lee Trout 2023
