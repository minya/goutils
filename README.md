# GoUtils

GoUtils is a collection of useful Go utilities for common programming tasks. This library provides simple, reusable components to help with file operations, web requests, and array manipulations.

## Installation

```
go get github.com/minya/goutils
```

## Packages

### config

Utilities for handling configuration files with JSON serialization and home directory path support.

```go
import "github.com/minya/goutils/config"

// Load configuration from a file
var cfg MyConfig
err := config.UnmarshalJson(&cfg, "~/configs/myapp.json")

// Save configuration to a file
err := config.MarshalJson(cfg, "~/configs/myapp.json")
```

### web

Utilities for web requests, including a custom cookie jar and HTTP transport configuration.

```go
import "github.com/minya/goutils/web"

// Create a custom cookie jar
jar := web.NewJar()

// Create an HTTP client with custom transport
client := &http.Client{
    Transport: web.DefaultTransport(5000), // 5000ms timeout
    Jar: jar,
}
```

## License

This project is available under the MIT License.
