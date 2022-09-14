#!/bin/bash
dir=$(cd -P -- "$(dirname -- "$0")" && pwd -P)
cd $dir/..

GIN_MODE=debug go run .