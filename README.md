# hevo-go

[![Go Reference](https://pkg.go.dev/badge/github.com/hevoio/hevo-go.svg)](https://pkg.go.dev/github.com/hevoio/hevo-go)
[![Testing](https://img.shields.io/github/actions/workflow/status/hevoio/hevo-go/testing.yml?branch=master)](https://github.com/hevoio/hevo-go/actions)
[![License](https://img.shields.io/github/license/hevoio/hevo-go)](./LICENSE)

The Go SDK for [Hevo](https://hevodata.com)'s API. The
[Hevo Terraform provider](https://github.com/hevoio/terraform-provider-hevo) is built on it.

## Getting started

```sh
go get github.com/hevoio/hevo-go
```

```go
client := hevo.New(accessKey, secretKey)
if err := client.Region("us"); err != nil {
	return err
}

pipeline, err := client.GetPipeline().PipelineID(42).Do(ctx)
if err != nil {
	return err // a *httputils.ResponseError; see hevo.IsNotFound, IsConflict, IsRateLimited
}

fmt.Println(pipeline.Data.Name, pipeline.Data.State)
```

Every operation is a method on the client named after it. Set its fields with chained setters
and send it with `Do(ctx)`:

```go
_, err := client.UpdatePipelineObjects().
	PipelineID(42).
	SelectionMode("EXCLUDE_ALL").
	Object(hevo.NewObjectSelection().Name("contacts").Selected(true)).
	Object(hevo.NewObjectSelection().Name("deals").Selected(true)).
	Do(ctx)
```

[`examples/`](./examples) has a runnable program for every operation.

## Operations

| Resource | Client methods |
| --- | --- |
| Pipelines | `ListPipelines`, `CreatePipeline`, `GetPipeline`, `UpdatePipeline`, `DeletePipeline`, `EnablePipeline`, `DisablePipeline`, `SyncPipelineNow`, `ResyncPipeline`, `UpdatePipelineSource`, `GetPipelinePosition` |
| Pipeline objects | `ListPipelineObjects`, `GetPipelineObject`, `UpdatePipelineObjects`, `ResyncPipelineObjects` |
| Jobs | `ListPipelineJobs`, `GetPipelineJob`, `ListPipelineJobObjects`, `CancelPipelineJob` |
| Destinations | `ListDestinations`, `CreateDestination`, `GetDestination`, `UpdateDestination`, `DeleteDestination` |
| Connectors | `ListSourceConnectors`, `GetSourceConnector`, `ListDestinationConnectors`, `GetDestinationConnector` |

## Client settings

| Setting | Default |
| --- | --- |
| `Region(name)`: one of `us`, `us2`, `eu`, `eu2`, `in`, `asia`, `au` | `us` |
| `BaseURL(url)`: any other host, such as a local stack | — |
| `CustomUserAgent(s)`: appended to `hevo-go/<version>` | — |
| `SetHTTPClient(c)` | `*http.Client` with a 60 s timeout |
| `SetHandleRateLimits(b)`: wait and retry on `429` | on |
| `SetMaxRetryAttempts(n)` | 5 |
| `SetMaxConcurrency(n)`: calls in flight per client | 4 |
| `hevo.Debug(true)`: print requests and responses to stderr | off |

**Retries.** A `429` is retried on every method after `Retry-After`, or `X-RateLimit-Reset` when
that is missing. `502`, `503`, `504` and network errors are retried with exponential backoff (1 s
to 30 s) on `GET` and `DELETE` only. A `POST` or `PATCH` that failed that way may have been
applied, so it is returned for the caller to check. A `429` that outlasts the retries carries the
quota, the bucket and the wait in its `RateLimit`.

**Asynchronous work.** Edits, actions and deletes answer `Accepted`. The Pipeline then reports a
transitional state such as `Initializing`; read it with `GetPipeline` until it settles. A `409`
means the Pipeline is mid-transition or has an edit open in the Hevo UI; retry once it settles.

## Development

```sh
make test   # unit tests against the in-memory mock in tests/mock
make lint   # golangci-lint, configured in .golangci.yml
make e2e    # end-to-end tests; needs HEVO_ACCESS_KEY, HEVO_SECRET_KEY and HEVO_BASE_URL
```

See [CONTRIBUTING.md](./CONTRIBUTING.md) for how an operation is written.

## License

[MIT](./LICENSE)
