#! /bin/bash

for i in $(ls /opt/onesie-configs/certs/); do
  echo $i
  cat /opt/onesie-configs/certs/$i/{privkey,fullchain}.pem /opt/onesie-configs/dhparam.pem > /opt/onesie-configs/hitch/$i.pem
  ls -l /opt/onesie-configs/hitch/$i.pem
done
