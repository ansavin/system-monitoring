#!/usr/bin/env sh

CODE=1

CLNT_LOG_FILE=$(mktemp)
SRV_LOG_FILE=$(mktemp)

# we should wait before we can get stats from client output 
WAIT_TIME=5

# fun for clearing app after tests

make grpc-server 1>/dev/null 2>&1

sudo chown 0:0 ./grpc-server
sudo chmod u+s ./grpc-server

make grpc-client 1>/dev/null 2>&1

./grpc-server 1> $SRV_LOG_FILE 2>&1 &

SRV_PID=$!

# simple test

sleep 1

./grpc-client 1> $CLNT_LOG_FILE 2>&1 &

CLNT_PID=$!

sleep $WAIT_TIME

# building & running server

if  grep -q 'CPU statistics:' $CLNT_LOG_FILE && \
    grep -vEq 'error|cannot|panic' $CLNT_LOG_FILE && \
    grep -vEq 'error|cannot|panic' $SRV_LOG_FILE
then 
    echo 'Test 1 passed'
    CODE=0
else
    echo 'Test 1 failed: no expected words "CPU statistics:" found in client`s output or error happens'
    echo 'Client output:'
    cat $CLNT_LOG_FILE
    echo 'Server output:'
    cat $SRV_LOG_FILE
    exit $CODE
fi

kill $CLNT_PID

# test that re-connection with other parameters works fine

# we reuse this vars to check if server normally works with
# both defaults and user-defined vars
WAIT_TIME=2
AVERAGING_TIME=2

./grpc-client -m $WAIT_TIME -a $AVERAGING_TIME 1> $CLNT_LOG_FILE 2>&1 &

CLNT_PID=$!

sleep $WAIT_TIME

if  grep -q 'CPU statistics:' $CLNT_LOG_FILE && \
    grep -vEq 'error|cannot|panic' $CLNT_LOG_FILE && \
    grep -vEq 'error|cannot|panic' $SRV_LOG_FILE
then 
    echo 'Test 2 passed'
    CODE=0
else
    echo 'Test 2 failed: no expected words "CPU statistics:" found in client`s output or error happens'
    echo 'Client output:'
    cat $CLNT_LOG_FILE
    echo 'Sevrer output:'
    cat $SRV_LOG_FILE
    exit $CODE
fi

# cleanup

make clean 1>/dev/null 2>&1

kill $SRV_PID $CLNT_PID
KILL_CODE=$?

rm $CLNT_LOG_FILE
rm $SRV_LOG_FILE

# exit

if [ $KILL_CODE != 0 ]
then
    exit $KILL_CODE
fi
exit $CODE