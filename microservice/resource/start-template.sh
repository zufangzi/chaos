#!/bin/bash
myIp=`hostname -i`
consulServer=${myIp%.*}".1:8500"
consul-template -consul=$consulServer -template "/home/work/nginx.conf:/etc/nginx/nginx.conf:nginx -s reload"
