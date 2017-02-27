#!/bin/bash
set -e
set -u
set -x

# Run with `make test`

#
# NPM
#

declare -a npmservices=("myaccount")
for npmservice in "${npmservices[@]}"
do
    cd $npmservice
    npm install
    lintcount=$(./node_modules/.bin/eslint src/ | wc -l)
    if [ "$lintcount" -gt 0  ]; then
        echo "eslint found files that need formatting - please fix!"
        exit 1
    fi
    cd ..
done

#
# GOLANG
#

go vet $(glide novendor)
go test -race -cover $(glide novendor)
go install -race -v $(glide novendor)

gocount=$(git ls-files | grep '.go$' | grep -v 'pb.go$' | grep -v 'bindata.go$' | xargs gofmt -e -l -s | wc -l)
if [ "$gocount" -gt 0 ]; then
	echo "Some Go files are not formatted. Check your formatting!"
	exit 1
fi

buildcount=$(buildifier -mode=check $(find . -iname BUILD -type f -not -path "./vendor/*") | wc -l)
if [ "$buildcount" -gt 0 ]; then
	echo "Some BUILD files are not formatted. Run make build-fmt"
	exit 1
fi

# Go through folders, and if they have go files then test
for pkg in $(go list ./... | grep -v /vendor/) ; do    
    # check for packages with auto-generated files
    relativeFolder=$(echo $pkg | sed -e "s/v2.staffjoy.com\///")
    if [ $(ls -1 $relativeFolder -- *.go 2>/dev/null | grep .pb.go | wc -l) -eq 0 ]; then
        if [ $(ls -1 $relativeFolder -- *.go 2>/dev/null | grep bindata.go | wc -l) -eq 0 ]; then
            golint -set_exit_status $pkg
        fi
    fi
done


# Test some dockerization
declare -a services=("www" "faraday" "account/api" "account/server" "myaccount" "whoami" "company/api" "company/server" "sms/server" "bot/server" "app"  "ical")
for service in "${services[@]}"
do
    bazel run //$service:docker
done

echo "Congratulations, brave warrior. Your tests have passed."
