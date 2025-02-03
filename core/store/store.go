package store

import _ "embed"

//go:embed fixtures/fixtures.sql
var fixturesSQL string

func FixturesSQL() string { return fixturesSQL }
