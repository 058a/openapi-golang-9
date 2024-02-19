## create
```
docker run -v ./:/migrations migrate/migrate -path=/migrations/ -database postgresql://user:password@openapi-db:5432/openapi?sslmode=disable create -ext sql -dir ./migrations -seq stock_item && sudo chown -R k8suser:k8suser db
```

## up
```
docker run --rm -v ./:/migrations --network deployments_default migrate/migrate -path=/migrations -database postgresql://user:password@openapi-db:5432/openapi?sslmode=disable up
```

## down
```
docker run --rm -v ./:/migrations --network deployments_default migrate/migrate -path=/migrations -database postgresql://user:password@openapi-db:5432/openapi?sslmode=disable down
```