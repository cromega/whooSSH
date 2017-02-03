#!/bin/bash

rerun -c -b -p "**/*.{go,html,js,css}" "killall whooSSH; go build && ./whooSSH"
