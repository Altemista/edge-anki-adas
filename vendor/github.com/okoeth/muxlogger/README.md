# Logging middleware for HTTP
Generic logger for HTTP middleware 

The logger can be used as follows:

First import the base package (consider using `glide` for vendoring).
```
import (
	...
	"github.com/okoeth/muxlogger"
	...
)
```

Add the logger to your HTTP middleware by wrapping the handler functions (here: named `handlerfunc`):
```
mux.HandleFunc(pat.Get("/v1/an/url"), Logger(handlerfunc))
```
