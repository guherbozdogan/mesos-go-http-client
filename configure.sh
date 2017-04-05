#!/bin/bash


#downloads mesos proto folder 
function downloadGitFolder {
cd mesos&& git init && git config core.sparsecheckout true &&  echo 'include/mesos/v1' > .git/info/sparse-checkout && git remote add -f origin git://github.com/apache/mesos.git &&  git pull origin master

}

if [ -d mesos  ]; then
	downloadGitFolder
else
	mkdir mesos
	downloadGitFolder
fi
