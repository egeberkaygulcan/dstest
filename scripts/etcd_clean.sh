 kill $(jps | grep 'etcd' | grep -v 'grep' | awk '{print $1}')
 rm -rf /tmp/etcd
 echo "Killed all etcd instances."
