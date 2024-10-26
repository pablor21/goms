api_docs:
	#  go-swagger3 --module-path . --main-file-path ./app/server/apps/api/routes.go  --output ./app/server/apps/api/docs/swagger.json
	swag init -g ./app/server/apps/api/routes.go -o app/server/apps/api/docs -ot json
build:
	go build -o ../bin/goms .
dev:
	CGO_ENABLED=false air serve 