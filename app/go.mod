module github.com/app

go 1.15

replace github.com/db => ../db

require (
	github.com/db v0.0.0-00010101000000-000000000000 // indirect
	github.com/lib/pq v1.10.2
)
