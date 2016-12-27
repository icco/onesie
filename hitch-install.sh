#! /bin/bash

git clone https://github.com/varnish/hitch.git
cd hitch
./bootstrap
./configure
ls -al *
make
sudo make install
