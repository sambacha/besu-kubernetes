
#!/bin/bash
# deploy/olm-catalog /deploy/olm-certified

# this is a WIP 
# exit when any command fails
set -e

VERSION=`grep "Version.*=.*\".*\"" version/version.go | sed "s,.*Version.*=.*\"\(.*\)\".*,\1,"`
OLD_VERSIONS="v1alpha3 v1alpha2"
DOCKER_IO_PATH="docker.io/as2"
REDHAT_REGISTRY_PATH="registry.connect.redhat.com/as2"
OPERATOR_IMAGE="$DOCKER_IO_PATH/as2-operator:${VERSION}"
OLM_CATALOG=deploy/olm-catalog
OLM_CERTIFIED=deploy/olm-certified
YAML_SCRIPT_FILE=.yq_script.yaml

# create yq template to append older CRD versions
rm -f $YAML_SCRIPT_FILE
for v in $OLD_VERSIONS; do
  cat << EOF >>$YAML_SCRIPT_FILE
- command: update
  path: spec.versions[+]
  value:
    name: $v
    served: true
    storage: false
EOF
done

# append older versions to CRD files
for crd in deploy/crds/*_crd.yaml; do
  yq w -i -s $YAML_SCRIPT_FILE $crd
done

RESOURCES="
  - kind: StatefulSets
    version: apps/v1
  - kind: Deployments
    version: apps/v1
  - kind: Pods
    version: v1
  - kind: Services
    version: v1
  - kind: ConfigMaps
    version: v1
  - kind: Secrets
    version: v1
"

cat << EOF >$YAML_SCRIPT_FILE
- command: update
  path: spec.install.spec.deployments[0].spec.template.spec.containers[0].image
  value: $OPERATOR_IMAGE
- command: update
  path: spec.install.spec.permissions[0].serviceAccountName
  value: as2-operator
- command: update
  path: spec.customresourcedefinitions.owned[0].resources
  value: $RESOURCES
- command: update
  path: spec.customresourcedefinitions.owned[1].resources
  value: $RESOURCES
- command: update
  path: spec.customresourcedefinitions.owned[2].resources
  value: $RESOURCES
- command: update
  path: spec.customresourcedefinitions.owned[3].resources
  value: $RESOURCES
- command: update
  path: spec.customresourcedefinitions.owned[4].resources
  value: $RESOURCES
- command: update
  path: spec.customresourcedefinitions.owned[0].displayName
  value: IndexerCluster
- command: update
  path: spec.customresourcedefinitions.owned[1].displayName
  value: LicenseMaster
- command: update
  path: spec.customresourcedefinitions.owned[2].displayName
  value: SearchHeadCluster
- command: update
  path: spec.customresourcedefinitions.owned[3].displayName
  value: Spark
- command: update
  path: spec.customresourcedefinitions.owned[4].displayName
  value: Standalone
- command: update
  path: metadata.annotations.alm-examples
  value: |-
    [{
      "apiVersion": "rpc.as2.network/v1beta1",
      "kind": "IndexerCluster",
      "metadata": {
        "name": "example",
        "finalizers": [ "rpc.as2.network/delete-pvc" ]
      },
      "spec": {
        "replicas": 1
      }
    },
    {
      "apiVersion": ".as2.network/v1beta1",
      "kind": "Standalone",
      "metadata": {
        "name": "example",
        "finalizers": [ "rpc.as2.network/delete-pvc" ]
      },
      "spec": {}
    }]
EOF
