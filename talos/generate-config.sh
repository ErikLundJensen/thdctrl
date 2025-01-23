# TODO: change node name and output file name
# TODO: error handling...

talosctl gen secrets --output gen/secrets.yaml

talosctl gen config --with-secrets gen/secrets.yaml --output-types talosconfig -o gen/talosconfig demo-1 https://138.201.200.236:6443

talosctl gen config \
        --output gen/c1.yaml                                      \
        --output-types controlplane                               \
        --with-cluster-discovery=false                            \
        --with-secrets gen/secrets.yaml                           \
        --config-patch @cluster.yaml                              \
        --config-patch @all-nodes.yaml                            \
        --kubernetes-version $KUBERNETES_VERSION                  \
        $CLUSTER_NAME                                             \
        $API_ENDPOINT

talosctl -n  138.201.200.236 -e 138.201.200.236 apply-config -f gen/c1.yaml --insecure

./gen-cilium.sh

kubectl apply -f gen/cilium.yaml 