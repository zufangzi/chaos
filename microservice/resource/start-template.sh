#!/bin/bash
myIp=`hostname -i`
consulServer=${myIp%.*}".1:8500"
consul-template -consul=$consulServer -template "/etc/consul-templates/nginx.conf:/etc/nginx/nginx.conf:nginx -s reload"
