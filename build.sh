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

IP=$(dig +short myip.opendns.com @resolver1.opendns.com)
USERNAME=$(git config user.name)
TEXT="\"User *$USERNAME* (IP: $IP) builded and pushed *$GIT* (version *$1*) to docker repository!\""
CHANNEL='"#hatch"'
USERNAME='"build-script"'
ICON_EMOJI='":rocket:"'
URL='https://hooks.slack.com/services/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'

curl -X POST -H 'Content-type: application/json' --data "{\"text\": $TEXT, \"channel\": $CHANNEL, \"username\": $USERNAME, \"icon_emoji\": $ICON_EMOJI}" $URL
