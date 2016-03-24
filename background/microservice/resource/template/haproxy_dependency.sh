#!/bin/sh

function fill_tcp_config(){
	for DATA in `cat $1`;do
	ARR=(${DATA//:/ })
	NAME=${ARR[0]}

	len=${#ARR[@]}
	if (( len == 1 ));then
	    echo "found http server dependency: $NAME"
	    continue
	fi

	PORT=${ARR[1]}     
	# 不能缩进。否则格式就乱了
	cat >> haproxy.cfg << EOF	

listen $NAME 127.0.0.1:$PORT
    mode tcp
    option tcplog
    balance    roundrobin
    {{range service "$NAME"}}
    server $NAME{{.Address}} {{.Address}}:{{.Port}};
    {{else}}server empty 127.0.0.1:65535;{{end}}
EOF
	done
}


function fill_http_config(){
	cat >> haproxy.cfg << EOF
	
frontend HTTP_SLB
    bind *:80
    log     global
    mode    http
    option  httplog
    acl valid_method method GET HEAD POST PUT DELETE OPTIONS
    http-request deny if !valid_method
EOF

	for DATA in `cat $1`;do
		ARR=(${DATA//:/ })  
		NAME=${ARR[0]}

		len=${#ARR[@]}
		if (( len > 1 ));then
		    echo "found tcp server dependency: $NAME"
		    continue
		fi
		cat >> haproxy.cfg << EOF
    acl $NAME hdr_beg(host) -i $NAME.service.consul
    use_backend $NAME.service.consul if $NAME
EOF

	done

	for DATA in `cat $1`;do
		ARR=(${DATA//:/ })  
		NAME=${ARR[0]}

		len=${#ARR[@]}
		if (( len > 1 ));then
		    continue
		fi
		cat >> haproxy.cfg << EOF

backend $NAME.service.consul
    mode http
    balance roundrobin
    option  redispatch
    option  httpclose
    option  forwardfor
    {{range service "$NAME"}}
    server $NAME{{.Address}} {{.Address}}:{{.Port}};
    {{else}}server empty 127.0.0.1:65535;{{end}}
EOF
	done
}

cat > haproxy.cfg << EOF
global
    maxconn 51200
    log 127.0.0.1 local6 info
    uid 0
    gid 0
    nbproc 1
    pidfile /home/work/haproxy.pid
    stats socket /var/run/haproxy.stat mode 666  	

# TODO：timeout threshold should be change by case.
defaults
    log    global
    option    dontlognull
    retries    3
    option redispatch
    maxconn 2000
    timeout connect 1d
    timeout client 1d
    timeout server 1d

EOF

# TODO: if dependency file not exists, then reload all service from consul
tcp_num=`cat $1 | grep ":" | wc -l`
total_num=`cat $1 | wc -l`
if (( tcp_num != total_num ));then
  fill_http_config $1
fi


# process tcp server load balance
if (( tcp_num > 0 ));then
  fill_tcp_config $1
fi


