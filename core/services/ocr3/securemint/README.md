Run integration test:

```bash
docker run --name cl-postgres -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=dbname -p 5432:5432 -d postgres
make setup-testdb
```

example command: `CL_DATABASE_URL=postgresql://chainlink_dev:insecurepassword@localhost:5432/chainlink_development_test?sslmode=disable go test -timeout 15m -run ^TestIntegration_LLO_evm_premium_legacy$ github.com/smartcontractkit/chainlink/v2/core/services/ocr2/plugins/llo -v`
