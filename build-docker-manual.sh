#!/bin/bash -v
set -e

# source .npmrc
npx semantic-release --dry-run --no-ci

appName=go-clean-arch

nextReleaseVersion=`cat .VERSION`
echo $nextReleaseVersion
# deployVersion=latest
deployVersion=$nextReleaseVersion

dockerImageTag=lhhoangit/${appName}:${deployVersion}
docker image rm -f ${dockerImageTag}

docker build -t ${appName}:${nextReleaseVersion} -f Dockerfile . 
docker tag ${appName}:${nextReleaseVersion} ${dockerImageTag}
docker push ${dockerImageTag}
