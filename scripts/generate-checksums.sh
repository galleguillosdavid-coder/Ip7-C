#!/bin/bash

# Generate SHA256SUMS for release artifacts

echo "Generating SHA256SUMS..."

find . -name "ipv7*" -o -name "FluxVPN*" | xargs sha256sum > SHA256SUMS.txt

echo "SHA256SUMS.txt generated."