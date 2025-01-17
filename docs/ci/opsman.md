# Test operations package

We use the functionality provided by the [test operations package] in order to
manage components used during tests, especially e2e.

The main functionality used currently in the tests is related to managing
containers, the package exposes the function `StartComponent` which takes a
container name (as defined in the [docker compose file]) and a variadic parameter
with a set of condition functions to check when the container can be considered
as ready. So we can call it without any condition like:
```go
operations.StartComponent("my-container")
```
or adding readiness conditions as:
```go
operations.StartComponent("my-container", func() (done bool, err error){
  // run some checks
  return true, nil
}, func()(done bool, err error){
  // run some other checks
  return true, nil
})
```

[test operations package]: ../../test/operations/
[docker compose file]: ../../docker-compose.yml
