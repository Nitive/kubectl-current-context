#!/bin/bash
set -euo pipefail

root="$(realpath "$(dirname "${BASH_SOURCE[0]}")")"
executable=$1

function error() {
  printf "\e[31;1m%s\e[0m\n" "$*" >&2;
}

function assertEqual() {
  if [ "$1" != "$2" ]; then
    error "Actual value '$1' is not equal expected '$2'"
    exit 1
  fi
}

KUBECONFIG="$root/kubeconfigs/one-cluster.yaml:$root/kubeconfigs/two-clusters.yaml:"

expectedContext="$(kubectl config current-context)"
expectedNamespace="$(kubectl config view --minify --output 'jsonpath={..namespace}')"

echo "Expected context: $expectedContext"
echo "Expected namespace: $expectedNamespace"

assertEqual "$($executable -o context)" "$expectedContext"
assertEqual "$($executable -o namespace)" "$expectedNamespace"
assertEqual "$($executable -o slug)" "$expectedContext/$expectedNamespace"
assertEqual "$($executable)" "$expectedContext/$expectedNamespace"
assertEqual "$($executable -o json)" "{\"context\":\"$expectedContext\",\"namespace\":\"$expectedNamespace\"}"

echo "All passed!"
