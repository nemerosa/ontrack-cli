#!/usr/bin/env bash

version=$1
if [ version == "" ]
then
    version=Snapshot
fi
echo version=$version

mode=$2

package_names=("yontrack")

# Add ontrack-cli for backward compatibility if mode is "release"
if [ "$mode" = "release" ]; then
    package_names+=("ontrack-cli")
fi

platforms=("windows/amd64" "darwin/amd64" "darwin/arm64" "linux/amd64" "linux/386")

for package_name in "${package_names[@]}"
do
    for platform in "${platforms[@]}"
    do
        platform_split=(${platform//\// })
        GOOS=${platform_split[0]}
        GOARCH=${platform_split[1]}
        output_name=$package_name'-'$GOOS'-'$GOARCH
        if [ $GOOS = "windows" ]; then
            output_name+='.exe'
        fi

        env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags "-X yontrack/config.Version=$version" -o $output_name $package
        if [ $? -ne 0 ]; then
            echo 'An error has occurred! Aborting the script execution...'
            exit 1
        fi
    done
done

