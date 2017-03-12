#!/bin/bash

file="main"

if [ -f $file ] ; then
    rm $file
fi

CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .
