#!/usr/bin/env bash

CODE=1

if [ `whoami` != 'root' ]
    then
        echo 'This script must be executed by root'
        exit $CODE
fi

make grpc-server

make grpc-client

./grpc-server &

SRV_PID=$!

sleep 1

./grpc-client > /tmp/testfile &

CLNT_PID=$!

sleep 5 #FIXME

if grep -q 'CPU statistics:' /tmp/testfile
    then 
        echo 'Test passed'
        CODE=0
    else
        echo 'Test failed: no expected words "CPU statistics:" found in client`s output'
        echo 'Client output:'
        cat /tmp/testfile
fi

kill $SRV_PID

kill $CLNT_PID

rm /tmp/testfile

exit $CODE