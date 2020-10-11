# Beacon

go test ./... -coverprofile fmtcoverage.html fmt
go tool cover -html=fmtcoverage.html -o coverage.html