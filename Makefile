.PHONY: admin-test

admin-generate:
	go generate ./wg-admin/app/api/v1
	go generate ./wg-admin/db

admin-test: admin-generate
	go test ./wg-admin/utils ./wg-admin/transaction ./wg-admin/storage/data \
		-coverprofile=coverage.out
	go tool cover -func coverage.out | grep -F total:

admin: admin-generate
	go build -o bin/wg-admin ./wg-admin
