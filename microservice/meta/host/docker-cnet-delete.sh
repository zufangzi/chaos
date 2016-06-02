#!/bin/sh
source /etc/profile
mkdir -p /var/log/docker
CNAME=$1
GW=10.32.33.1

if [ -z $CNAME ];then
      echo "Please Enter ContainerID...." >> /var/log/docker/net.log
      exit 1
fi

HIP=`ifconfig veno1 | grep -w inet | awk '{print $2}'`
CIP=`/usr/local/etcd/etcdctl --peers http://10.32.27.82:4001 ls /coreos.com/network/WorkNet/CIP/ | grep $CNAME | awk -F- '{print $NF}'`

#if [ -z $CIP ];then
#   echo "Error: $CNAME ContainerID Not Injected To ETCD ...." >> /var/log/docker/net.log
#   exit 1
#fi

/usr/bin/ovs-vsctl del-port ovs-br0 vh"$CNAME" > /dev/null 2>&1
COUNT=`ovs-vsctl show | grep ${CNAME} | wc -l`

if [ $COUNT -eq 2 ];then
      echo "Error: $CNAME OVS Deleted Fail ..." >> /var/log/docker/net.log
      exit 1
fi

rm -rf /var/run/netns/$CNAME

/usr/local/etcd/etcdctl --peers http://10.32.27.82:4001 rmdir /coreos.com/network/WorkNet/CIP/${CNAME}-${CIP} > /dev/null 2>&1
/usr/local/etcd/etcdctl --peers http://10.32.27.82:4001 mkdir /coreos.com/network/WorkNet/IP/$CIP > /dev/null 2>&1

echo "$CNAME Deleted Successfully ..." >> /var/log/docker/net.log
