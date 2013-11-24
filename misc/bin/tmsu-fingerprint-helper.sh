#!/bin/bash

while read file
do
	FILE_PATH="$(basename $(readlink "${file}"))"
	echo $FILE_PATH|sed 's/^.\+--\([^.]\+\)\(\..\+\)\?$/\1/'
done
