# Creating a calculator webservice with Go

## Project

1. `go mod init github.com/brunorene/calculator-service`
1. Create Makefile with: setup, clean, vet, lint, check-system-deps, check-format, format, build
1. `go get go.uber.org/zap`
1. Setup logger
1. Create Infof log - highlight `f` convention - f stands for formatted
1. go vet heuristics to find conceptual errors on code
1. defer
1. Error handling, wrap errors
1. New package operator
1. Create Add, Subtract, Multiply, Divide with interface Operator - Notice uppercase to be accessible on main.go
1. Function receiver - like an extension method
1. `go get github.com/stretchr/testify`
1. Generate test for Add, Subtract
1. Enable code coverage Ctrl+Shift+P coverage
1. Create runServer - talk about http server
1. Create handler - split path
1. Show https://goplay.tools/ - test split there
1. Handler talk about json, struct tags, wrapped errors, nolint
1. Dockerfile - use multi-stage, base image scratch - binary only
1. NFTs
1. `go get github.com/tsenart/vegeta/v12`
1. create nft_test.go


