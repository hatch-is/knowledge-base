#!/bin/sh

set -e


if [ -z "$1" ]; then
    echo "[ERROR] Specify version as a first argument. Example: ./build.sh v0.0.1"
    exit
fi

eval $(aws ecr get-login --region us-east-1)
ECR_PATH=093525834944.dkr.ecr.us-east-1.amazonaws.com/hatch-insp/knowledge-base

GIT=git@github.com:hatch-is/knowledge-base.git

rm -rf ./_build/.cache/
git clone $GIT ./_build/.cache/knowledge-base
cd ./_build/.cache/knowledge-base
git checkout $1
cd ../
docker build -t $ECR_PATH:$1 ../
docker push $ECR_PATH:$1

