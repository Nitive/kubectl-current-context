.PHONY: *

run:
	@ go run main.go

benchmark-go:
	@ go build -ldflags '-s -w' -o kubectl-current-context main.go
	@ hyperfine --warmup 3 './kubectl-current-context -o json'

benchmark-kubectl:
	@ hyperfine --warmup 3 'kubectl config current-context && kubectl config view --minify --output "jsonpath={..namespace}"'