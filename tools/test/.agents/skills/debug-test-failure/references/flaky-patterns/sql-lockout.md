<symptom>
SQL lock errors.
`failed to create ...: ERROR: canceling statement due to lock timeout (SQLSTATE 55P03)`
</symptom>

<fix>
Randomize unique database keys.
Avoid:
```go
// Collision: multiple iterations of this test use the same ID
owner := Keccak256([]byte(t.Name()))[:20]
name := t.Name()
```

Prefer:
```go
// Isolation: every iteration/process gets a unique row
owner := testutils.NewAddress().Bytes()
name := testutils.RandomizeName(t.Name())
```

</fix>
