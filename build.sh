#!/usr/bin/env bash

# Get the program name and version
MOCKERFILE="mock.go"
PROGRAMNAME=`grep "programName " ${MOCKERFILE} | cut -d'=' -f2|sed 's/"//g'`
PROGRAMVERSION=`grep "programVersion" ${MOCKERFILE} | cut -d'=' -f2|sed 's/"//g'|sed -e 's/^[[:space:]]*//'`
PLATFORM=("windows/amd64" "darwin/amd64" "linux/amd64")

# Loop through the platform and build a package
for platform in "${PLATFORM[@]}"
do

    # Extract the os and architecture
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}

    echo "Building package for platform ${GOOS} architecture ${GOARCH}"

    # Build the package name
    output_name=${PROGRAMNAME}'-'${GOOS}'-'${GOARCH}'-'${PROGRAMVERSION}
    if [[ ${GOOS} = "windows" ]]; then
        output_name+='.exe'
    fi

    # Build the package using go build
    env GOOS=${GOOS} GOARCH=${GOARCH} go build -o ${output_name}
    if [[ $? -ne 0 ]]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done