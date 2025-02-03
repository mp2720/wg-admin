.PHONY: admin-test

admin-test:
	go test ./wg-admin/utils ./wg-admin/transaction ./wg-admin/storage/data \
		-coverprofile=coverage.out
	go tool cover -func coverage.out | grep -F total:

admin-generate:
	sqlc generate -f wg-admin/db/sql/sqlc.yaml
