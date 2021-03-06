#!/bin/bash

##note/guherbozdogan: 
#add  dependency checking of git version (> 1.7) and 
#add dependency checking of installation of github.com/gogo/protobuf projects


#downloads mesos proto folder 
function downloadGitFolder {
cd mesos&& git init && git config core.sparsecheckout true &&  echo 'include/mesos/v1' > .git/info/sparse-checkout && git remote add -f origin git://github.com/apache/mesos.git &&  git pull origin master && mv include/mesos/v1 v1 && rm -rf include/mesos
}


#generates the protobuf go files 
function callProtoBuffGen {
cd mesos && find ./v1/.  -type f ! -name '*.proto' -exec rm {} \; && @find . -type f -name '*.proto' -exec protoc   --proto_path=${GOPATH}/src:${GOPATH}/src/github.com/gogo/protobuf/protobuf:. --gogo_out=.   {} \;
}

##note/guherbozdogan:
#add checks of outputs of the functions 
if [ -d mesos  ]; then
	downloadGitFolder 
	callProtoBuffGen
else
	mkdir mesos
	downloadGitFolder
	callProtoBuffGen
fi
