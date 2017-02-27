set -e

# Removes mac shitty things
find assets/ -type f -name '.DS_Store' -delete

# Put assets into the binary
go-bindata assets/...
# Clean up data so it passes linter
gofmt -s -w bindata.go
# We need to make it work with linting
sed -i "s/package main/package errorpages/g" bindata.go
sed -i "s/Css/CSS/g" bindata.go
sed -i "s/Json/JSON/g" bindata.go
echo "THAT WAS EASY."
