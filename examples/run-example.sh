#!/bin/sh

export API_KEY='YOUR_API_KEY'
export API_SECRET='YOUR_SECRET_KEY'

script=example.go
if [ x"${1}" != x"" ]; then
    script="${1}"
fi

go run "${script}"
