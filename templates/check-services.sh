#! /bin/bash

param=$(netstat -tnl | grep :443)
if [[ -z "${param// }" ]]; then
  systemctl restart hitch
fi

param=$(netstat -tnl | grep :8080)
if [[ -z "${param// }" ]]; then
  systemctl restart nginx
fi

param=$(netstat -tnl | grep :9090)
if [[ -z "${param// }" ]]; then
  systemctl restart status-server
fi

param=$(systemctl | grep fail | grep varnish)
if [[ -z "${param// }" ]]; then
  systemctl restart varnishlog
  systemctl restart varnishncsa
fi
