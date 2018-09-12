#! /bin/sh
set -e

# NOTE
# build script for cli interface with cross-compile environment
#               for arm         : ./build.sh arm
#               for amd64       : ./build.sh
# ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

export BINARY_NAME="cupactive-server"      #fixed~~~~~

export CONTAINER_VERSION=":v1.0"

#if [ "$1" = "arm" ]; then
#        echo "****************************"
#                echo "Target Binary arch is ARM"
#                echo "****************************"
#        export GOARCH=arm GOARM=7
        #export CC="arm-linux-gnueabihf-gcc"
#        export CC="arm-none-eabi"
#        export CONTAINER_NAME="armv7/cupactive-sever"
#else
#                echo "****************************"
#                echo "Target Binary arch is amd64"
#                echo "****************************"
#        export GOARCH=amd64
#        export CGO_ENABLED=0
#        export CC="gcc"
#        export CONTAINER_NAME="cupactive-server"
#fi

echo make clean
make clean

echo make build
make build
