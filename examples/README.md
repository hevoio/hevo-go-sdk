# Examples

One runnable program per operation, at `examples/<operation>/main.go`. Each reads the team's API
credentials and region from the environment:

```sh
export HEVO_ACCESS_KEY=... HEVO_SECRET_KEY=... HEVO_REGION=us
go run ./examples/listPipelines
```

The ids, names and settings in the examples are made up. Change them to your own before running
an example: those that create, update or delete act on your Hevo account.
