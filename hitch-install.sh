#! /bin/bash

cd /tmp
git clone https://anonscm.debian.org/git/collab-maint/hitch.git
cd hitch
sudo aptititude install -y libev-dev libssl-dev automake python-docutils flex bison pkg-config
./bootstrap
./configure
ls -al
make
sudo make install
