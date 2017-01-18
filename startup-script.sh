#! /bin/bash

# Mount configs
sudo gcsfuse -o allow_other --uid=$(id -u www-data) --gid=$(id -g www-data) onesie-configs /opt/onesie-configs/
echo "mounted!"
ls -al /opt/onesie-configs/

# Fix varnish log
sudo chown -R varnish:varnish /var/lib/varnish
sudo systemctl restart varnish
sudo systemctl restart varnishlog
sudo systemctl status varnish*
