# About
Example of transport generator (tg) for http and jRPC for golang 

# Refs to docs of transport geenrator
https://github.com/seniorGolang/tg/tree/v2

Note that this repo shows usage of v2 version of tg

# Update transport
```bash
go generate ./...
```

# Run servcie

```bash
go run ./cmd/main.go
```

# API Example

Request: 
```bash
curl -X POST "http://localhost:9000/api/v1/mathematical" -H "Content-Type: application/json" -d '[{"id":"123", "jsonrpc":"2.0", "method":"Add", "params":{"A":1, "B":2}}]'
```

Correct response: 
```bash
[{"id":"123","jsonrpc":"2.0","result":{"result":3}}]
```

Example of error response: 
```bash
[{"id":"123","jsonrpc":"2.0","error":{"code":-32601,"message":"invalid method 'mathematical.Add'"}}]
```



Log on server:
```
2022-05-23T17:45:17+03:00 ??? hello world
2022-05-23T17:45:17+03:00 INF listen on bind=:9000
2022-05-23T17:45:51+03:00 INF call add method=add request="{A:1 B:2}" response={Result:3} service=Mathematical took="19.536µs"
```
