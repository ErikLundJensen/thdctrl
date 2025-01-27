# TODO: change node name and output file name
# TODO: error handling...

# TODO: Only during bootstrapping of cluster:
#talosctl gen secrets --output gen/secrets.yaml
#talosctl gen config --with-secrets gen/secrets.yaml --output-types talosconfig -o gen/talosconfig demo-1 https://138.201.200.236:6443

cp all-nodes.yaml gen/all-nodes.yaml

yq -i '.machine.install.disk = "/dev/nvme0n1" | .machine.network.hostname = "node1"' gen/all-nodes.yaml

talosctl gen config \
        --output gen/c1.yaml                                      \
        --output-types controlplane                               \
        --with-cluster-discovery=false                            \
        --with-secrets gen/secrets.yaml                           \
        --config-patch @cluster.yaml                              \
        --config-patch @gen/all-nodes.yaml                        \
        --kubernetes-version $KUBERNETES_VERSION                  \
        $CLUSTER_NAME                                             \
        $API_ENDPOINT

# TODO: loop through all nodes
yq -i '.machine.install.disk = "/dev/sda" | .machine.network.hostname = "node2"' gen/all-nodes.yaml

talosctl gen config \
        --output gen/c2.yaml                                      \
        --output-types controlplane                               \
        --with-cluster-discovery=false                            \
        --with-secrets gen/secrets.yaml                           \
        --config-patch @cluster.yaml                              \
        --config-patch @gen/all-nodes.yaml                        \
        --kubernetes-version $KUBERNETES_VERSION                  \
        $CLUSTER_NAME                                             \
        $API_ENDPOINT

#talosctl -n  138.201.200.236 -e 138.201.200.236 apply-config -f gen/c1.yaml --insecure
#talosctl -n  136.243.103.75 -e 136.243.103.75 apply-config -f gen/c2.yaml --insecure

#TODO: only at first node
#./gen-cilium.sh
#kubectl apply -f gen/cilium.yaml 