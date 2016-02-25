#!/bin/sh

echo "
daemon off;
user              root;
worker_processes  1;

error_log  /var/log/nginx/error.log;

pid        /var/run/nginx.pid;

events {
    worker_connections  1024;
}

http {
    include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    log_format  main  '\$remote_addr - \$remote_user [\$time_local] "\$request" '
                     '\$status \$body_bytes_sent "\$http_referer" '
                      '"\$http_user_agent" "\$http_x_forwarded_for"';

    access_log  /var/log/nginx/access.log  main;

    sendfile        on;
    #tcp_nopush     on;

    #keepalive_timeout  0;
    keepalive_timeout  65;

    #gzip  on;

    include /etc/nginx/conf.d/*.conf;" > nginx.conf

# TODO: if dependency file not exists, then reload all service from consul
for NAME in `cat $1`;do
 
echo "
 upstream $NAME{    
        {{range service \"$NAME\"}}  
        server {{.Address}}:{{.Port}} weight=10; 
         {{else}}server 127.0.0.1:65535;{{end}} 
    } 

    server { 
      listen 80; 
      server_name $NAME.service.consul; 
      root /home/work/$NAME; 
      index index.html; 
      location / {  
        proxy_pass        http://$NAME; 
        proxy_set_header  X-Real-IP  \$remote_addr; 
        proxy_redirect          off;  
        proxy_set_header        Host \$host;
        proxy_set_header        X-Real-IP \$remote_addr; 
        proxy_set_header        X-Forwarded-For \$proxy_add_x_forwarded_for; 
        client_max_body_size    100m; 
        client_body_buffer_size 128k;
        proxy_connect_timeout   300; 
        proxy_send_timeout      300; 
        proxy_read_timeout      300;
        proxy_buffer_size       4k; 
        proxy_buffers           4 32k;
        proxy_busy_buffers_size 64k; 
        proxy_temp_file_write_size 64k; 
      } 
    }" >> nginx.conf

done

echo "}" >> nginx.conf
