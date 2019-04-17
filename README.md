# Akouste

Akouste or ακούστε is a configuration manager for application in kubernetes.

Akouste is using `Consul` as KV store and `consul-template` for listening configuration changes.

## Operator

### Development

#### Running Locally

Testing configuration changes on local environment

1. Go to ./cmd/operator
2. docker-compose up -d 
3. make run

You should be able to see `consul-ui` by accessing `localhost:8500/ui`. To check if reloading happened or not, please check the `appoperator` logs.

## Downloader

Downloader downloads+unarchives given file and make sure there
are only `n` number of downloaded files (deletion starts from the oldest modified file).

Supported download src protocol:
- `local` for local filesystem
```
$ ./configdownloader \
	-bucketProto "local" \
	-bucketName "test/local-bucket" \
	...
```
- `gs` for Google Cloud Storage
```
$ ./configdownloader \
	-bucketProto "gs" \
	-bucketName "test-bucket-name" \
	...
```

### Downloader Example

#### Run `downloader`
```
$ ./configdownloader \
	-logLevel debug \
	-bucketProto "local" \
	-bucketName "test/local-bucket" \
	-downloadDIR "test/local-downloads" \
	-keepOldCount 5
```

#### Request format

Accepted `POST` form:
- `uri`: filepath relative to the `bucket`
- `unarchive`: whether to unarchive the downloaded file (`true`/`false`)

cURL example:

```
$ curl -X POST -d "uri=config-1.tar.gz&unarchive=true" localhost:9000/v1/download
```
