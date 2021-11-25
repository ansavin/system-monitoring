#!/usr/bin/env sh

CODE=1

CLNT_LOG_FILE=$(mktemp)
SRV_LOG_FILE=$(mktemp)

WAIT_TIME=5

make grpc-server 1>/dev/null 2>&1

sudo chown 0:0 ./grpc-server
sudo chmod u+s ./grpc-server

make grpc-client 1>/dev/null 2>&1

./grpc-server 1> $SRV_LOG_FILE 2>&1 &

SRV_PID=$!

sleep 1

./grpc-client $WAIT_TIME 1 1> $CLNT_LOG_FILE 2>&1 &

CLNT_PID=$!

sleep $WAIT_TIME #FIXME

if  grep -q 'CPU statistics:' $CLNT_LOG_FILE && \
    grep -vEq 'error|cannot|panic' $CLNT_LOG_FILE && \
    grep -vEq 'error|cannot|panic' $SRV_LOG_FILE
then 
    echo 'Test passed'
    CODE=0
else
    echo 'Test failed: no expected words "CPU statistics:" found in client`s output'
    echo 'Client output:'
    cat $CLNT_LOG_FILE
    echo 'Sevrer output:'
    cat $SRV_LOG_FILE
fi

make clean 1>/dev/null 2>&1

kill $SRV_PID $CLNT_PID
KILL_CODE=$?

rm $CLNT_LOG_FILE
rm $SRV_LOG_FILE

if [ $KILL_CODE != 0 ]
then
    exit $KILL_CODE
fi
exit $CODE