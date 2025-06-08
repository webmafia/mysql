module github.com/webmafia/mysql

go 1.23.0

// replace github.com/webmafia/fast => ../go-fast

require (
	github.com/cespare/xxhash/v2 v2.3.0
	github.com/go-sql-driver/mysql v1.9.2
	github.com/webmafia/fast v0.17.0
	github.com/webmafia/lru v1.0.0
)

require filippo.io/edwards25519 v1.1.0 // indirect
