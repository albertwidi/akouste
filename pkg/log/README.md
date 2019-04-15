# TDK Log

TDK Log is based on [Zerolog](https://github.com/rs/zerolog) log package

## Log level

Log level is supported in the logger, available log level is:

- Debug
- Info
- Warning
- Error
- Fatal

Log is disabled if `LogLevel` < `CurrentLogLevel`, for example `Debug` log is disabled when current level is `Info`

Example of `SetLevel`:

```go
import "github.com/tokopedia/tdk/log"

func main() {
    log.SetLevel(log.InfoLevel)
    log.Infow("this is a log", "key1", "val1")
}
```

## Log to file

### For >= Info Level

All logs are written to `stderr`, but we can also write the log to file by using:

```go
import "github.com/tokopedia/tdk/log"

func main() {
    err := log.SetConfig(&log.Config{LogFile: "logfile.log"})
    if err != nil {
        panic(err)
    }
    log.Infow("this is a log", "key1", "val1")
}
```

## Key-value context in log

To add more context to log, `key-value` fields is provided. For example:

```go
import "github.com/tokopedia/tdk/log"

func main() {
    log.Infow("this is a log", "key1", "val1")
}
```

## Integration with TDK Error package

TDK error package has a features called `errors.Fields`. This fields can be used to add more context into the error, and then we can print the fields when needed. TDK log will automatically print the fields if `error = tdkerrors.Error` by using `log.Errors`. For example:

```go
import "github.com/tokopedia/tdk/log"
import "github.com/tokopedia/tdk/errors"

func() {
    err := errors.E("this is an error", errors.Fields{"field1":"value1"})
    log.Errors(err)
}

// result is
// message=this is an error field1=value1
```