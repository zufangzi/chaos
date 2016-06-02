haproxy -D -f /etc/haproxy/haproxy.cfg -st $(ps aux | grep haproxy | grep '\-f \/etc\/haproxy\/haproxy.cfg' | grep -v 'consul-template' | grep -v 'grep' | head -n1 | awk '{print $2}')
