# local
```
go test -cover ./... -coverprofile=cover.out && go tool cover -html=cover.out -o cover.html
```

# docker
```
docker exec -it openapi-golang go clean -testcache && docker exec -it openapi-golang go test ./../../...
```