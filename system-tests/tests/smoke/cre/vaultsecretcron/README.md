# vaultsecretcron

Cron-triggered vault secret reader for the sharded shared-vault tests
(`shard_shared_vault_test.go`). On every scheduled run it fetches one secret
(`config.SecretNamespace` / `config.SecretKey`) and emits its plaintext value in a
user log — `Vault secret fetched: <value>` — so a test can tell, from that single
log line, which shard executed the run and that the fetched secret actually
decrypted there.

Deliberately minimal; it is the cron-triggered counterpart of the neighboring
[`vaultsecret`](../vaultsecret) workflow, which is an HTTP-triggered *verifier* with
per-invocation phases/checks for exercising vault CRUD behavior (owner mismatch,
identifier validation, batch limits, delete liveness). This one exists to answer one
question per run — "did this shard read this secret" — which is all the sharding and
failover scenarios need, and it needs no gateway trigger plumbing.
