#!/bin/sh

for file in test_neigh*.json; do
        jq -Mc < "$file" > "$file.tmp" && mv "$file.tmp" "$file"
done
