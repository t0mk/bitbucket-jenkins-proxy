#!/bin/sh

rm -f bitbucket-jenkins-proxy
docker build -t t0mk/bitbucket-jenkins-proxy-builder -f Dockerfile.build ./
docker run t0mk/bitbucket-jenkins-proxy-builder
LAST=`docker ps -ql`
docker cp ${LAST}:/bitbucket-jenkins-proxy ./
docker build -t t0mk/bitbucket-jenkins-proxy ./
