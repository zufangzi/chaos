#!/bin/sh
ovs-vsctl del-br ovs-br0
ifconfig eno1 0 up
ovs-vsctl add-br ovs-br0
ovs-vsctl add-port ovs-br0 veno1 tag=27 -- set interface veno1 type=internal
ovs-vsctl add-port ovs-br0 eno1 trunk=27,33
ip link set ovs-br0 up
ifconfig veno1 10.32.27.80/24
route add default gw 10.32.27.1 veno1
