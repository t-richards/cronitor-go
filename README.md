# cronitor-go

[![Go Reference](https://pkg.go.dev/badge/github.com/t-richards/cronitor-go.svg)](https://pkg.go.dev/github.com/t-richards/cronitor-go)

A minimal Go client library to send telemetry events to [Cronitor](https://cronitor.io).

## Usage

```go
package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/t-richards/cronitor-go"
)

func main() {
	if err := nightlyJob(); err != nil {
		log.Fatal(err)
	}
}

func nightlyJob() error {
	// Create a new Cronitor client with your API key and job name.
	crn := cronitor.New("https://cronitor.link/p/<API-KEY>/nightly-job")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Optionally, specify a custom HTTP client to change timeouts, etc.
	crn.HTTPClient = &http.Client{Timeout: 3 * time.Second}

	// Mark the job as started.
	_ = crn.Run(ctx)

	// Perform your job logic here.
	err := doWork(ctx)
	if err != nil {
		// Mark the job as failed.
		_ = crn.Fail(ctx)
		return err
	}

	// Mark the job as completed.
	_ = crn.Complete(ctx)

	return nil
}
```

For simplicity, the above example does not handle API errors from Cronitor. Consider handling these errors in your application as needed. Cronitor API errors are usually the result of a service outage, so we recommend treating Cronitor reporting as a best-effort operation. For example,

```go
err := crn.Run(ctx)
if err != nil {
    // Couldn't report the job start to Cronitor.
    log.Println(err)
}
```

## Roadmap

This package is intentionally very simple and only supports the `Jobs` API with minimal parameters. If you need more advanced functionality, please consider using a different library.

PRs for minimal functionality additions _may_ be accepted on a case-by-case basis. Before submitting a PR, please open an issue to discuss your proposed changes.

## License

[MIT](LICENSE).
