#!/bin/bash
GOOS=linux GOARCH=amd64 revel package bullhorn
scp bullhorn.tar.gz root@107.170.105.233:~/app