#!/bin/bash
myIp=`hostname -i`
#consulServer=${myIp%.*}".1:8500"
#consulServer="10.14.5.14:8500"
consulServer=`cat /home/work/data/$(hostname)/kickoff_sign_file`":8500"
#restart-haproxy-template
#{
#haproxy -D -f /etc/haproxy/haproxy.cfg -st $(ps aux | grep haproxy | grep '\-f \/etc\/haproxy\/haproxy.cfg' | grep -v 'consul-template' | grep -v 'grep' | head -n1 | awk '{print $2}')
#}
consul-template -consul=$consulServer -template "/home/work/haproxy.cfg:/etc/haproxy/haproxy.cfg:sh /home/work/restart-haproxy-template.sh"
