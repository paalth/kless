#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

cd cmd/frontend/common

./build.sh

cd ../../..
