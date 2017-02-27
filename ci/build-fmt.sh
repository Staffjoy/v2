set -e 

buildifier -showlog -mode=fix $(find . -iname BUILD -type f | grep -v node_modules)