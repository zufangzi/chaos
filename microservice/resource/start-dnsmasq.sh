#!/bin/sh
echo '
nameserver 127.0.0.1
nameserver 10.14.5.57
nameserver 10.14.5.58
nameserver 172.16.3.124
'> /etc/resolv.conf

cp /etc/resolv.conf /etc/resolv.dnsmasq.conf
cp /etc/hosts /etc/dnsmasq.hosts

# startup dnsmasq
dnsmasq --no-daemon --log-queries

