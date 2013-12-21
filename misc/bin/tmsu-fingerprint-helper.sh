#!/bin/bash

while read file
do
	FILE_PATH="$(basename $(readlink "${file}"))"
	if [ "$?" != "0" ]
	then
		FILE_PATH=internal
	fi
	echo $FILE_PATH|sed 's/^.\+--\([^.]\+\)\(\..\+\)\?$/\1/'
done
