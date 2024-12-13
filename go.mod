module github.com/demayl/go-http-server

go 1.23.4

//require internal/server/server v1.0.0
replace internal/server => ./internal/server/

replace internal/server/args => ./internal/server/args

require (
	internal/server v0.0.0-00010101000000-000000000000
	internal/server/args v0.0.0-00010101000000-000000000000
)
