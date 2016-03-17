#!/bin/bash


echo "first sleep..."
sleep 5
echo "now begin to found whether haproxy had been started up"

HA_CNTS=`ps aux | grep haproxy | grep "\-f \/etc\/haproxy\/haproxy.cfg" | grep -v "grep" | wc -l`

if (( HA_CNTS == 0 ));then
# 先以后台方式启动haproxy
echo "clear... now begin to start haproxy daemon"
haproxy -D -f /etc/haproxy/haproxy.cfg 
fi

echo "now begin to listen..."
EXIT_THD=3
CURRENT_WARN=0
while true
do
    echo "in monitor loop..."
    sleep 5
    cnts=`ps aux | grep haproxy | grep "\-f \/etc\/haproxy\/haproxy.cfg" | grep -v "grep" | wc -l`
    if (( cnts == 0 ));then
        CURRENT_WARN=`expr $CURRENT_WARN + 1`
        echo "no haproxy found..."
    else
        echo "normal..."
        CURRENT_WARN=0
    fi 
    if (( CURRENT_WARN >= EXIT_THD));then
        echo "fail too many times(3) then will be exited"
        exit 2
    fi
done
