api_docs:
	#  go-swagger3 --module-path . --main-file-path ./app/server/apps/api/routes.go  --output ./app/server/apps/api/docs/swagger.json
	swag init -g ./app/server/apps/api/routes.go -o app/server/apps/api/docs -ot json
	npx swagger-typescript-api -p ./app/server/apps/api/docs/swagger.json  -o api

gen_repos:
	go run . generate repos

build:
	go build -o ../bin/goms .

dev:
	CGO_ENABLED=false air serve 