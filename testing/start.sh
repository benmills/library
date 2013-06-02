base_port=5000
node_count=$1

if [ -z $node_count ]; then
  node_count=3
fi

(sleep 2; $PWD/link.sh $base_port $node_count && echo "Linked!") &
foreman start --port $base_port -c library=$node_count
