# Run a test

1. Start a postgres container.
```
docker run -d --rm --name chainlink-postgres -e POSTGRES_PASSWORD=postgres -e POSTGRES_HOST_AUTH_METHOD=trust -v $HOME/chainlink-pg-data/:/var/lib/postgresql/data -p 5432:5432 postgres:14 postgres -N 500 -B 1024MB
```

2. From this repo root directory.
```
export CL_DATABASE_URL="postgresql://postgres:@localhost:5432/chainlink_test?sslmode=disable"
make testdb
```

3. From this directory, run a test.
```
go test -count 1 -v -run "Test_CCIPBatching" .
```
