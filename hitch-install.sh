#! /bin/bash

cd /tmp
git clone https://anonscm.debian.org/git/collab-maint/hitch.git
cd hitch
./bootstrap
./configure
ls -al
make
sudo make install
