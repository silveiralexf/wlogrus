# wlogrus

A simple wrapper for logrus package to facilitate formating logs to console output or forwarding structured logs into a central stash.

## Install

Get the latest version with `go get` as shown below:

```shell
go get github.com/silveiralexf/wlogrus
```


## Usage

For using this wrapper, just follow the sample code below:

```go
package main

import (
	"github.com/silveiralexf/wlogrus"
)

func main() {
	informational()
	somethingSmelling()
	nothingToSeeHere()
}

func informational() {
    wlogrus.Info("Example1","info message are plain and simple")
    wlogrus.Debug("Example1","debug msg only shown when env var WLOGRUS_DEBUG=true",wlogrus.CallerInfo())
}

func somethingSmelling() {
	wlogrus.Warn("Example2","something is about to go wrong")
	wlogrus.Error("Example2","something failed, but it might recover/proceed",wlogrus.CallerInfo())
}

func nothingToSeeHere() {
	wlogrus.Fatal("Example3","halts execution with an error",wlogrus.CallerInfo())
}
```

Output of the example provided should result at the following:

***Note**: for enabling debug messages to be displayed, export environment variable `WLOGRUS_DEBUG=true`*

```shell
$ export WLOGRUS_DEBUG=true
$ go run examples/examples.go
2020-12-11 19:24:03 [INFO] [Example1] info message are plain and simple
2020-12-11 19:24:03 [DEBUG] [Example1] debug messages are only displayed when env var WLOGRUS_DEBUG=true
2020-12-11 19:24:03 [WARNING] [Example2] warn messages inform something is about to go wrong
2020-12-11 19:24:03 [ERROR] [Example2] error messages signals that something failed, but it might recover/proceed [examples.go.main.somethingSmelling:33]
2020-12-11 19:24:03 [FATAL] [Example3] fatal messages will halt execution with an error [examples.go.main.nothingToSeeHere:41]
exit status 1
```

For a structured output export environment variable `WLOGRUS_JSON=true` as shown below:

```json
$ export WLOGRUS_JSON=true
$ go run examples/examples.go
{"body":"info message are plain and simple","level":"info","msg":"[Example1] info message are plain and simple","severity":"INFO","tag":"Example1","time":"2020-12-11T19:27:48-03:00"}
{"body":"warn messages inform something is about to go wrong","level":"warning","msg":"[Example2] warn messages inform something is about to go wrong","severity":"WARNING","tag":"Example2","time":"2020-12-11T19:27:48-03:00"}
{"body":"error messages signals that something failed, but it might recover/proceed","caller":"examples.go.main.somethingSmelling:33","level":"error","msg":"[Example2] error messages signals that something failed, but it might recover/proceed [examples.go.main.somethingSmelling:33]","severity":"ERROR","tag":"Example2","time":"2020-12-11T19:27:48-03:00"}
{"body":"fatal messages will halt execution with an error","caller":"examples.go.main.nothingToSeeHere:41","level":"fatal","msg":"[Example3] fatal messages will halt execution with an error [examples.go.main.nothingToSeeHere:41]","severity":"FATAL","tag":"Example3","time":"2020-12-11T19:27:48-03:00"}
exit status 1
```

