
#!/bin/bash
# Script to export besu-operator image to tarball and generated combined YAML

# exit when any command fails
set -e

#VERSION=$1
VERSION=1.0.1
if [[ "x$VERSION" == "x" ]]; then
    # Use latest commit id if no version is provided
    VERSION=`git rev-parse HEAD | cut -c1-12`
fi
#IMAGE="docker.io/besu/besu-operator:${VERSION}"
#IMAGE_ID=`docker images besu/besu-operator:latest -q`
IMAGE="quay.io/freight/besu-operator:${VERSION}"
IMAGE_ID=`quay images besu/besu-operator:latest -q`

echo Tagging image ${IMAGE_ID} as besu/besu-operator:${VERSION}
docker tag ${IMAGE_ID} besu/besu-operator:${VERSION}

echo Generating release-${VERSION}/besu-operator-${VERSION}.tar.gz
mkdir -p release-${VERSION}
rm -f release-${VERSION}/*
docker image save besu/besu-operator:${VERSION} | gzip -c > release-${VERSION}/besu-operator-${VERSION}.tar.gz

echo Generating release-${VERSION}/besu-operator-crds.yaml
touch release-${VERSION}/besu-operator-crds.yaml
for i in `ls deploy/crds/|grep crd.yaml`
do
    echo "---" >> release-${VERSION}/besu-operator-crds.yaml
    cat deploy/crds/$i >> release-${VERSION}/besu-operator-crds.yaml
done

echo Generating release-${VERSION}/besu-operator-noadmin.yaml
cat deploy/service_account.yaml deploy/role.yaml deploy/role_binding.yaml > release-${VERSION}/besu-operator-noadmin.yaml
echo "---" >> release-${VERSION}/besu-operator-noadmin.yaml
yq w deploy/operator.yaml "spec.template.spec.containers[0].image" $IMAGE >> release-${VERSION}/besu-operator-noadmin.yaml

echo Generating release-${VERSION}/besu-operator-install.yaml
cat release-${VERSION}/besu-operator-crds.yaml release-${VERSION}/besu-operator-noadmin.yaml > release-${VERSION}/besu-operator-install.yaml

echo Rebuilding release-${VERSION}/besu-operator-cluster.yaml
cat release-${VERSION}/besu-operator-crds.yaml deploy/namespace.yaml > release-${VERSION}/besu-operator-cluster.yaml
echo "---" >> release-${VERSION}/besu-operator-cluster.yaml
yq w deploy/service_account.yaml metadata.namespace besu-operator >> release-${VERSION}/besu-operator-cluster.yaml
echo "---" >> release-${VERSION}/besu-operator-cluster.yaml
yq w deploy/role.yaml metadata.namespace besu-operator | yq w - kind ClusterRole >> release-${VERSION}/besu-operator-cluster.yaml
echo "---" >> release-${VERSION}/besu-operator-cluster.yaml
yq w deploy/role_binding.yaml metadata.namespace besu-operator | yq w - roleRef.kind ClusterRole >> release-${VERSION}/besu-operator-cluster.yaml
cat deploy/cluster_role.yaml deploy/cluster_role_binding.yaml >> release-${VERSION}/besu-operator-cluster.yaml
echo "---" >> release-${VERSION}/besu-operator-cluster.yaml
yq w deploy/operator.yaml metadata.namespace besu-operator | yq w - "spec.template.spec.containers[0].image" $IMAGE | yq w - "spec.template.spec.containers[0].env[0].value" "" | yq d - "spec.template.spec.containers[0].env[0].valueFrom" >> release-${VERSION}/besu-operator-cluster.yaml

ls -la release-${VERSION}/