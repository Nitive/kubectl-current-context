# kubectl-current-context

Get current kubectl context and namespace **super fast**. Useful for shell prompts.

## Options

Get context/namespace
```
$ kubectl-current-context -o slug
my-cluster/default
```

Custom separator with `-s` flag
```
$ kubectl-current-context -s='⎈ '
my-cluster⎈default
```

Just context
```
$ kubectl-current-context -o context
my-cluster
```

Just namespace
```
$ kubectl-current-context -o namespace
default
```

JSON
```
$ kubectl-current-context -o json
{"context":"my-cluster","namespace":"default"}
```

## Benchmarks

Native kubectl commands to get current context and namespace take ~136ms to run

```
$ kubectl config current-context
$ kubectl config view --minify --output "jsonpath={..namespace}"
```

```
Benchmark #1: kubectl config current-context && kubectl config view --minify --output "jsonpath={..namespace}"
  Time (mean ± σ):     135.7 ms ±   6.9 ms    [User: 139.9 ms, System: 30.5 ms]
  Range (min … max):   123.3 ms … 164.4 ms    100 runs
```

This implementation takes only ~7ms — 19 times less

```
$ kubectl-current-context -o json
```

```
Benchmark #1: ./kubectl-current-context -o json
  Time (mean ± σ):       7.2 ms ±   3.0 ms    [User: 6.3 ms, System: 6.0 ms]
  Range (min … max):     4.9 ms …  34.7 ms    100 runs
```

# Installation

Download [latest release](https://github.com/Nitive/kubectl-current-context/releases/latest) binary and put it in your /usr/local/bin directory
