#!/bin/sh
source /etc/profile
mkdir -p /var/log/docker
CNAME=$1

if [ -z $CNAME ];then
      echo "Error: Please Enter ContainerID...." >> /var/log/docker/net.log
      exit 1
fi

IP="10.32.33."`/usr/local/etcd/etcdctl --peers http://10.32.27.82:4001 ls /coreos.com/network/WorkNet/IP/ | awk -F. '{print $NF}' | sort -n | awk 'NR==1 {print $0}'`
HIP=`ifconfig veno1 | grep -w inet | awk '{print $2}'`

ovs-vsctl del-port ovs-br0 vh"$CNAME" > /dev/null 2>&1
rm -rf /var/run/netns/$CNAME
echo "$CNAME Stop Successfully ..." >> /var/log/docker/net.log
