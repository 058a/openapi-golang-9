~/go/bin/oapi-codegen -generate types,server -package hello ./../../api/hello.yaml > ./../../internal/infra/oapicodegen/hello/hello.go && go mod tidy

~/go/bin/oapi-codegen -generate types,server -package item ./../../api/stock/item.yaml > ./../../internal/infra/oapicodegen/stock/item/item.go && go mod tidy

~/go/bin/oapi-codegen -generate types,server -package location ./../../api/stock/location.yaml > ./../../internal/infra/oapicodegen/stock/location/location.go && go mod tidy