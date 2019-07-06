set -e 

buildifier -mode=fix $(find . -iname BUILD -type f | grep -v node_modules)