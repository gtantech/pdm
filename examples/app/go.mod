module github.com/gtantech/pdm/examples/app

go 1.27

replace github.com/gtantech/pdm => ../../

require (
	github.com/golang-migrate/migrate/v4 v4.20.1
	github.com/gtantech/pdm v0.0.0-00010101000000-000000000000
	modernc.org/sqlite v1.58.0
)

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gtantech/go-container v1.0.1 // indirect
	github.com/gtantech/toposort/v2 v2.2.0 // indirect
	github.com/mattn/go-isatty v0.0.24 // indirect
	github.com/ncruces/go-strftime v1.0.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	golang.org/x/sys v0.47.0 // indirect
	modernc.org/libc v1.75.6 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.12.1 // indirect
)
