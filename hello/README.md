# Hello

Application greets the universe.

## Local
The `hello/local` example does not take any arguments.

## Distributed
The `hello/distributed` example involves to peers: Alice and Bob.

Both peers take optionally parameters `address` and `port`.

Defaults are:
- `address` is `localhost`. 
- `port` is port number `31415`.

Execution of the peers can be done done as follows:

For Alice:
```terminal
cd $GOPATH/src/github.com/pspaces/gospace-examples
go run hello/distributed/alice/main.go
```

For Bob:
```terminal
cd $GOPATH/src/github.com/pspaces/gospace-examples
go run hello/distributed/bob/main.go
```
