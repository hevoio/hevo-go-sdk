# Contributing

## Adding an operation

Each API operation is one file in its resource's package, named `<resource>_<operation>.go`
(`pipeline_get.go`, `pipelines_list.go`), with its unit test beside it and a program in
`examples/<operation>/main.go`.

An operation has up to three types:

- **Service** (`GetPipelineService`): exported. Embeds `httputils.HTTPService`; every field is
  an unexported pointer, set by a chained setter, so an unset field is left out of the request.
  `Do(ctx)` checks the path parameters are set, sends the request and returns the Response.
- **Request** (`pipelineUpdateRequest`): unexported, in the package's `common_types.go`. Pointer
  fields with `omitempty`. Only operations with a body have one.
- **Response** (`PipelineDetailsResponse`): exported, in `common_types.go`. Value fields, except
  where the API returns `null` and the difference matters.

A nested part of a body, such as a Pipeline's schedule, is a builder of its own (`Schedule`)
with a `Request()` method, a constructor in `hevo.go` (`hevo.NewSchedule()`), and a value-typed
`…Response` for reads.

Add the client method to `client.go`, named after the operation's `operationId` in the API spec.

## Rules

- Standard library only.
- No business-rule validation: the API owns that. Check only what is needed to build the URL.
- Names follow the API: Pipeline, Source, Destination, object, job, connector.
- Run `make lint test` before opening a PR. Name branches `feature/<description>` or
  `fix/<description>`, and write PR titles as `type(scope): description`, for example
  `fix(pipelines): send the cursor on the second page`. PRs are squash-merged, so the title becomes
  the commit message.
