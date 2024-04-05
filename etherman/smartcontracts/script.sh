#!/bin/sh

set -e

gen() {
    local package=$1

    abigen --bin bin/${package}.bin --abi abi/${package}.abi --pkg=${package} --out=${package}/${package}.go
}

gen preetrogpolygonzkevmglobalexitroot
gen preetrogpolygonzkevmbridge
gen preetrogpolygonzkevm
gen elderberrypolygonzkevm
gen etrogpolygonzkevm
gen etrogpolygonzkevmbridge
gen pol
gen etrogpolygonzkevmglobalexitroot
gen etrogpolygonrollupmanager
gen mocketrogpolygonrollupmanager
gen mockverifier
gen proxy