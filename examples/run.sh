#!/bin/sh

export BITFLYER_API_KEY='YOUR_API_KEY'
export BITFLYER_API_SECRET='YOUR_API_SECRET'

script=example.go
if [ x"${1}" != x"" ]; then
    script="${1}"
fi

go run "${script}"
