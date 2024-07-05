#!/bin/bash
kill $(jps | grep 'ratis-examples' | grep -v 'grep' | awk '{print $1}')
rm -rf /tmp/ratis
echo "All Ratis examples have been stopped."