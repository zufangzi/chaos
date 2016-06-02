#!/bin/sh
source /etc/profile
mkdir -p /var/log/docker
CNAME=$1
IP=$2
GW=10.32.33.1

echo "IP is: "$IP >> /var/log/docker/net.log
echo "CNAME is: "$CNAME >> /var/log/docker/net.log

if [ -z $CNAME ];then
      echo "Error: Please Enter ContainerID...." >> /var/log/docker/net.log
      exit 1
fi

#IP=`/usr/local/etcd/etcdctl --peers http://10.32.27.82:4001 ls /coreos.com/network/WorkNet/CIP/ | grep $CNAME | awk -F- '{print $NF}'`
HIP=`ifconfig veno1 | grep -w inet | awk '{print $2}'`
pid=$(docker inspect -f '{{.State.Pid}}' $CNAME)

if [ -z $pid ];then
      echo "Error: $CNAME ContainerID Not Exist...." >> /var/log/docker/net.log
      exit 1
elif [ $pid -eq 0 ];then
      echo "Error: $CNAME ContainerID Is Stoped ...." >> /var/log/docker/net.log
      exit 1
fi
#if [ -z $IP ];then
#IP="10.32.33."`/usr/local/etcd/etcdctl --peers http://10.32.27.82:4001 ls /coreos.com/network/WorkNet/IP/ | awk -F. '{print $NF}' | sort -n | awk 'NR==1 {print $0}'`
#fi
#/usr/local/etcd/etcdctl --peers http://10.32.27.82:4001 rmdir /coreos.com/network/WorkNet/IP/$IP > /dev/null 2>&1
#/usr/local/etcd/etcdctl --peers http://10.32.27.82:4001 mkdir /coreos.com/network/WorkNet/CIP/${CNAME}-${IP} > /dev/null 2>&1

echo "PID is: "$pid >> /var/log/docker/net.log
ovs-vsctl del-port ovs-br0 vh"$CNAME" > /dev/null 2>&1
#docker start $CNAME
echo "BEGIN to build the bridge..."
mkdir -p /var/run/netns
ln -s /proc/$pid/ns/net /var/run/netns/$CNAME
ip netns exec $CNAME ip link del eth0
ip link add name vh"$CNAME" mtu 1500 type veth peer name vc"$CNAME" mtu 1500
ovs-vsctl add-port ovs-br0 vh"$CNAME" tag=33 > /dev/null 2>&1
ip link set vh"$CNAME" up 
ip link set vc"$CNAME" netns $CNAME
ip netns exec $CNAME ip link set dev vc"$CNAME" name eth0
ip netns exec $CNAME ip addr add "$IP"/24 dev eth0
ip netns exec $CNAME ip link set eth0 up
ip netns exec $CNAME route add default gw $GW dev eth0
mkdir -p /data/container/$CNAME
echo "$HIP" > /data/container/$CNAME/kickoff_sign_file
if [ $? -ne 0 ];then
      echo "$CNAME Container Start Fail ..." >> /var/log/docker/net.log
      exit 1
fi
COUNT=`ovs-vsctl show | grep ${CNAME} | wc -l`
if [ $COUNT -ne 2 ];then
      echo "$CNAME OVS Create Net Fail ..." >> /var/log/docker/net.log
      exit 1
fi 
#COUNT2=`/usr/local/etcd/etcdctl --peers http://10.32.27.82:4001 ls /coreos.com/network/WorkNet/CIP/ | grep $CNAME | wc -l`
#if [ $COUNT2 -ne 1 ];then
#   echo "Error: $CNAME NET Regist Fail ..." >> /var/log/docker/net.log
#   exit 1
#fi
echo "$CNAME Start Successfully ..." >> /var/log/docker/net.log
