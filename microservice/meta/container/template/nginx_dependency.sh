#!/bin/sh

cat > nginx.conf << EOF
daemon off;
user              root;
worker_processes  1;

error_log  /home/work/nginx/logs/error.log;

pid        /home/work/nginx/logs/nginx.pid;

events {
    worker_connections  1024;
}

http {
    #include       /etc/nginx/mime.types;
    default_type  application/octet-stream;

    log_format  main  '\$remote_addr - \$remote_user [\$time_local] "\$request" '
                     '\$status \$body_bytes_sent "\$http_referer" '
                      '"\$http_user_agent" "\$http_x_forwarded_for"';

    access_log  /home/work/nginx/logs/access.log  main;

    sendfile        on;

    keepalive_timeout  65;

EOF


# TODO: if dependency file not exists, then reload all service from consul
for DATA in `cat $1`;do
ARR=(${DATA//:/ })  
NAME=${ARR[0]}

len=${#ARR[@]}
if (( len > 1 ));then
    echo "found tcp server dependency: $NAME"
    continue
fi
 
cat >> nginx.conf << EOF
    upstream $NAME{    
        {{range service "$NAME"}}  
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
    }
EOF

done

echo "}" >> nginx.conf

# process tcp server load balance

num=`cat $1 | grep ":" | wc -l`

if [ $num = "0" ];then
  exit 0
fi

cat >> nginx.conf << EOF
stream {
EOF

for DATA in `cat $1`;do
ARR=(${DATA//:/ })
NAME=${ARR[0]}

len=${#ARR[@]}
if (( len == 1 ));then
    echo "found http server dependency: $NAME"
    continue
fi

PORT=${ARR[1]}     

cat >> nginx.conf << EOF
	
        upstream $NAME {
            {{range service "$NAME"}}
              server {{.Address}}:{{.Port}} weight=10;
            {{else}}server 127.0.0.1:65535;{{end}}
        }
        server {
                listen $PORT;
		#server_name $NAME.service.consul; 
                proxy_connect_timeout 1s;
                proxy_timeout 3s;
                proxy_pass $NAME;
        }
EOF
done
echo "}" >> nginx.conf


