#! /bin/bash

cd /tmp
git clone https://anonscm.debian.org/git/collab-maint/hitch.git -b jessie-backports
cd hitch
sudo aptitude install -y libev-dev libssl-dev automake python-docutils flex bison pkg-config
./configure
ls -al *
make
