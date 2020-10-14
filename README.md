

# Besu-Kubernetes (k8s)

The following repo has example reference implementations of private networks using k8s. This is intended to get developers and ops people familiar with how to run a private ethereum network in k8s and understand the concepts involved.

It provides examples using multiple tools such as kubectl, helm, helmfile etc. Please select the one that meets your deployment requirements.

## Local Development:
The reference examples in this repo can be used locally, to get familiar with the deployment setup. You will require:
- [Minikube](https://kubernetes.io/docs/setup/learning-environment/minikube/) This is the local equivalent of a K8S cluster
- [Helm](https://helm.sh/docs/)
- [Helmfile](https://github.com/roboll/helmfile)
- [Helm Diff plugin](https://github.com/databus23/helm-diff)


Minikube defaults to 2 CPU's and 2GB of memory, unless configured otherwise. We recommend you starting with at least 8GB, depending on the amount of nodes you are spinning up - the recommended requirements for each besu node are 4GB
```bash
minikube start --memory 16384
# or with RBAC
minikube start --memory 16384 --extra-config=apiserver.Authorization.Mode=RBAC
```

Verify kubectl and minikube are working with
```bash
$ kubectl version
Client Version: version.Info{Major:"1", Minor:"15", GitVersion:"v1.15.1", GitCommit:"4485c6f18cee9a5d3c3b4e523bd27972b1b53892", GitTreeState:"clean", BuildDate:"2019-07-18T09:18:22Z", GoVersion:"go1.12.5", Compiler:"gc", Platform:"linux/amd64"}
Server Version: version.Info{Major:"1", Minor:"15", GitVersion:"v1.15.0", GitCommit:"e8462b5b5dc2584fdcd18e6bcfe9f1e4d970a529", GitTreeState:"clean", BuildDate:"2019-06-19T16:32:14Z", GoVersion:"go1.12.5", Compiler:"gc", Platform:"linux/amd64"}
```

Install helm & helm-diff:
Please note that the documentation and steps listed use *helm3*. The API has been updated so please take that into account if using an older version
```bash
$ helm plugin install https://github.com/databus23/helm-diff --version master
```

Pick the deployment tool that suits your environment and then change directory and follow the Readme.md files there.


Namespaces:
Currently we do not deploy anything in the 'default' namespace. Anything related to Besu gets spun up in a 'besu' namespace, the monitoring pieces get spun up in a 'monitoring' namespace.
Namespaces are part of the setup and do not need to be created via kubectl prior to deploying. To change the namespaces:
- In Kubectl, you need to edit every file in the deployment
- In Helm, edit the namespace value in the values.yaml 

It is recommended you follow this approach as well in your production setups and where possible use Service Accounts to secure deployments & statefulsets. We make use of these extensively.


Monitoring:
The example setups all have out custom Grafana [dashboard](https://grafana.com/grafana/dashboards/10273) to make monitoring of the nodes and network easier.

We use the default ports 3000 for grafana and 9090 for prometheus deployments, and the respective Nodeport services use ports 30030 and 30090. Credentials are `admin:password` Open the 'Besu Dashboard' to see the status of the nodes on your network. If you do not see the dashboard, click on Dashboards -> Manage and select the dashboard from there

Please configure the kubernetes scraper and grafana security to suit your requirements, Grafana supports [multiple options](https://grafana.com/docs/grafana/latest/auth/overview/) that can be configured using env vars

Ingress Controllers:

If you require the use of ingress controllers for the RPC calls or the monitoring dashboards, we have provided [examples](./helm/ingress/README.md) with rules that are configured to do so.
 
Please use these as a reference and develop solutions to match your network topology and requirements.


## Production Network Guidelines:
| ⚠️ **Note**: After you have familiarised yourself with the examples in this repo, it is recommended that you design your network based on your needs, taking the following guidelines into account |
| --- |


#### Pod Resources:
The templates in this repository have been set to run locally on Minikube to get the user familiar with the setup. Hence the resources are set to:

```bash
besu:
  pvcSizeLimit: "1Gi"
  pvcStorageClass: "standard"
  memRequest: "1024Mi"
  memLimit: "2048Mi"
  cpuRequest: "100m"
  cpuLimit: "500m"

orion:
  pvcSizeLimit: "1Gi"
  pvcStorageClass: "standard"
  memRequest: "512Mi"
  memLimit: "1024Mi"
  cpuRequest: "100m"
  cpuLimit: "500m"
```

When designing your setup to run in `staging` or `production` environments, please ensure you **grant at least 4GB of memory to Besu pods and 2GB of memory to Orion pods.** Also ensure you **select the appropriate storage class and size for your nodes.**

Ensure that if you are using a cloud provider you have enough spread across AZ's to minimize risks - refer to our [HA](https://besu.hyperledger.org/en/latest/HowTo/Configure/Configure-HA/High-Availability/) and [Load Balancing] (https://besu.hyperledger.org/en/latest/HowTo/Configure/Configure-HA/Sample-Configuration/) documentation

When deploying a private network, eg: IBFT you need to ensure that the bootnodes are accessible to all nodes on the network. Although the minimum number needed is 1, we recommend you use more than 1 spread across AZ's. In addition we also recommend you spread validators across AZ's and have a sufficient number available in the event of an AZ going down.

You need to ensure that the genesis file is accessible to all nodes joining the network.



#### Network Topology and High Availability requirements:
Ensure that if you are using a cloud provider you have enough spread across AZ's to minimize risks - refer to our [HA](https://besu.hyperledger.org/en/latest/HowTo/Configure/Configure-HA/High-Availability/) and [Load Balancing] (https://besu.hyperledger.org/en/latest/HowTo/Configure/Configure-HA/Sample-Configuration/) documentation

When deploying a private network, eg: IBFT you need to ensure that the bootnodes are accessible to all nodes on the network. Although the minimum number needed is 1, we recommend you use more than 1 spread across AZ's. In addition we also recommend you spread validators across AZ's and have a sufficient number available in the event of an AZ going down.

You need to ensure that the genesis file is accessible to all nodes joining the network.

Hyperledger Besu supports [NAT mechanisms](https://besu.hyperledger.org/en/stable/Reference/CLI/CLI-Syntax/#nat-method) and the default is set to automatically handle NAT environments. If you experience issues with NAT and logs have messages that have the NATService throwing exceptions connecting to external IPs, please add this option in your Besu deployments `--nat-method = NONE` 

#### Data Volumes:                     
Ensure that you provide enough capacity for data storage for all nodes that are going to be on the cluster. Select the appropriate [type](https://kubernetes.io/docs/concepts/storage/volumes/) of persitent volume based on your cloud provider.

#### Nodes:
Consider the use of statefulsets instead of deployments for nodes. The term 'node' refers to bootnode, validator and network nodes.

Configuration of nodes can be done either via a single item inside a config map, as Environment Variables or as command line options. Please refer to the [Configuration](https://besu.hyperledger.org/en/latest/HowTo/Configure/Using-Configuration-File/) section of our documentation

#### RBAC:
We encourage the use of RBAC's for access to the private key of each node, ie. only a specific pod/statefulset is allowed to access a specific secret. If you need to specify a Kube config file to each pod please use the `KUBE_CONFIG_PATH` variable

#### Monitoring
As always please ensure you have sufficient monitoring and alerting setup.

Besu publishes metrics to [Prometheus](https://prometheus.io/) and metrics can be configured using the [kubernetes scraper config](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#kubernetes_sd_config).

Besu also has a custom Grafana [dashboard](https://grafana.com/grafana/dashboards/10273) to make monitoring of the nodes easier.

For ease of use, the kubectl & helm examples included have both installed and included as part of the setup. Please configure the kubernetes scraper and grafana security to suit your requirements, grafana supports multiple options that can be configured using env vars

We also have [example ingress controller config rules](./helm/ingress/), should you want to use a setup that makes use of ingress controllers. 

#### Logging
Besu's logs can be [configured](https://besu.hyperledger.org/en/latest/HowTo/Troubleshoot/Logging/#advanced-custom-logging) to suit your environment. For example, if you would like to log to file and then have parsed via logstash into an ELK cluster, please follow out documentation.

### New nodes joining the network:
The general rule is that any new nodes joining the network need to have the following accessible:
- genesis.json of the network
- Bootnodes need to be accessible on the network
- Bootnodes enode's (public key and IP) should be passed in at boot
- If you’re using permissioning on your network, specifically authorise the new nodes

If the initial setup was on Kubernetes, you have the following scenarios:

#### 1. New node also being provisioned on the K8S cluster:
In this case anything that applies to how current nodes are provisioned should be applicable and the only thing that need be done is increase the number of replicas

#### 2. New node being provisioned elsewhere
Ensure that the host being provisioned can find and connect to the bootnode's. You may need to use `traceroute`, `telnet` or the like to ensure you have connectivity. Once connectivity has been verified, you need to pass the enode of the bootnodes and the genesis file to the node. This can be done in many ways, for example query the k8s cluster via APIs prior to joining if your environment allows for that. Alternatively put this data somewhere accessible to new nodes that may join in future as well, and pass the values in at runtime.

Ensure that the host being provisioned can also connect to the other nodes that you have on the k8s cluster, otherwise it will be unable to connect to any peers (bar the bootnodes). The most reliable way to do this is via a VPN so it has access to the bootnodes as well as any nodes on the k8s cluster. You can alternatively use ingresses on the nodes (ideally more than just bootnodes) you wish to expose, where TCP & UDP on port 30303 need to be open for discovery.

Additionally if you’re using permissioning on your network you will also have to specifically authorise the new nodes


# Architecture

OLM is composed of two Operators: the OLM Operator and the Catalog Operator.

Each of these Operators is responsible for managing the CRDs that are the basis for the OLM framework:

| Resource                 | Short Name | Owner   | Description                                                                                |
|--------------------------|------------|---------|--------------------------------------------------------------------------------------------|
| ClusterServiceVersion | csv        | OLM     | application metadata: name, version, icon, required resources, installation, etc...        |
| InstallPlan           | ip         | Catalog | calculated list of resources to be created in order to automatically install/upgrade a CSV |
| CatalogSource         | catsrc         | Catalog | a repository of CSVs, CRDs, and packages that define an application                        |
| Subscription          | sub        | Catalog | used to keep CSVs up to date by tracking a channel in a package                            |
| OperatorGroup         | og     | OLM     | used to group multiple namespaces and prepare for use by an operator                     |

Each of these Operators are also responsible for creating resources:

| Operator | Creatable Resources        |
|----------|----------------------------|
| OLM      | Deployment                 |
| OLM      | Service Account            |
| OLM      | (Cluster)Roles             |
| OLM      | (Cluster)RoleBindings      |
| Catalog  | Custom Resource Definition |
| Catalog  | ClusterServiceVersion      |

## What is a ClusterServiceVersion

ClusterServiceVersion combines metadata and runtime information about a service that allows OLM to manage it.

ClusterServiceVersion:

- Metadata (name, description, version, links, labels, icon, etc)
- Install strategy
  - Type: Deployment
    - Set of service accounts / required permissions
    - Set of deployments

- CRDs
  - Type
  - Owned - managed by this service
  - Required - must exist in the cluster for this service to run
  - Resources - a list of k8s resources that the Operator interacts with
  - Descriptors - annotate CRD spec and status fields to provide semantic information

## OLM Operator

The OLM Operator is responsible for installing applications defined by a `ClusterServiceVersion` resources once the required resources specified in the CSV are present in the cluster (the job of the Cluster Operator). This may be as simple as setting up a single `Deployment` resulting in an operator's pod becoming available. The OLM Operator is not concerned with the creation of the underlying resources. If this is not done manually, the Catalog Operator can help provide resolution of these needs.

This separation of concerns enables users incremental buy-in of the OLM framework components. Users can choose to manually create these resources, or define an InstallPlan for the Catalog Operator or allow the Catalog Operator to develop and implement the InstallPlan. An operator creator does not need to learn about the full operator package system before seeing a working operator.

While the OLM Operator is often configured to watch all namespaces, it can also be operated alongside other OLM Operators so long as they all manage separate namespaces defined by `OperatorGroups`.

### ClusterServiceVersion Control Loop

```
           +------------------------------------------------------+
           |                                                      |
           |                                      +--> Succeeded -+
           v                                      |               |
None --> Pending --> InstallReady --> Installing -|               |
           ^                                       +--> Failed <--+
           |                                              |
           +----------------------------------------------+
\                                                                 /
 +---------------------------------------------------------------+
    |
    v
Replacing --> Deleting
```

| Phase      | Description                                                                                                                                                                                                                           |
|------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| None       | initial phase, once seen by the Operator, it is immediately transitioned to `Pending`                                                                                                                                                 |
| Pending    | requirements in the CSV are not met, once they are this phase transitions to `Installing`                                                                                                                                             |
| InstallReady | all requirements in the CSV are present, the Operator will begin executing the install strategy                                                                                                                                     |
| Installing | the install strategy is being executed and resources are being created, but not all components are reporting as ready                                                                                                                 |
| Succeeded  | the execution of the Install Strategy was successful; if requirements disappear, or an APIService cert needs to be rotated this may transition back to `Pending`; if an installed component disappears this may transition to `Failed`|
| Failed     | upon failed execution of the Install Strategy, or an installed component disappears the CSV transitions to this phase; if the component can be recreated by OLM, this may transition to `Pending`                                     |
| Replacing  | a newer CSV that replaces this one has been discovered in the cluster. This status means the CSV is marked for GC                                                                                                                     |
| Deleting   | the GC loop has determined this CSV is safe to delete from the cluster. It will disappear soon.                                                                                                                                       |
> Note: In order to transition, a CSV must first be an active member of an OperatorGroup

## Catalog Operator

The Catalog Operator is responsible for monitoring `Subscriptions`, `CatalogSources` and the catalogs themselves. When it finds a new or changed `Subscription`, it builds out the subscribed Operator. When it finds a new or changed CatalogSource it builds out the required catalog, if appropriate, and begins regular monitoring of the package in the catalog. The packages in the catalog will include various `ClusterServiceVersions`, `CustomResourceDefinitions` and a channel list for each package. A catalog has packages. A package has channels and CSVs. A Channels identifies a specific CSV. The CSVs identify specific CRDs.

A user wanting a specific operator creates a `Subscription` which identifies a catalog, operator and channel within that operator. The Catalog Operator then receives that information and queries the catalog for the latest version of the channel requested. Then it looks up the appropriate `ClusterServiceVersion` identified by the channel and turns that into an `InstallPlan`. When updates are found in the catalog for the channel, a similar process occurs resulting in a new `InstallPlan`. (Users can also create an InstallPlan resource directly, containing the names of the desired ClusterServiceVersions and an approval strategy.) 

When the Catalog Operator find a new `InstallPlan`, even though it likely created it, it will create an "execution plan" and embed that into the `InstallPlan` to create all of the required resources. Once approved, whether manually or automatically, the Catalog Operator will implement its portion of the the execution plan, satisfying the underlying expectations of the OLM Operator. 

The OLM Operator will then pickup the installation and carry it through to completion of everything required in the identified CSV.

### InstallPlan Control Loop

```
None --> Planning +------>------->------> Installing --> Complete
                  |                       ^
                  v                       |
                  +--> RequiresApproval --+
```

| Phase            | Description                                                                                    |
|------------------|------------------------------------------------------------------------------------------------|
| None             | initial phase, once seen by the Operator, it is immediately transitioned to `Planning`         |
| Planning         | dependencies between resources are being resolved, to be stored in the InstallPlan `Status` |
| RequiresApproval | occurs when using manual approval, will not transition phase until `approved` field is true    |
| Installing       | resolved resources in the InstallPlan `Status` block are being created                      |
| Complete         | all resolved resources in the `Status` block exist                                             |

### Subscription Control Loop

```
None --> UpgradeAvailable --> UpgradePending --> AtLatestKnown -+
         ^                                   |                  |
         |                                   v                  v
         +----------<---------------<--------+---------<--------+
```

| Phase            | Description                                                                                                   |
|------------------|---------------------------------------------------------------------------------------------------------------|
| None             | initial phase, once seen by the Operator, it is immediately transitioned to `UpgradeAvailable`                |
| UpgradeAvailable | catalog contains a CSV which replaces the `status.installedCSV`, but no `InstallPlan` has been created yet |
| UpgradePending   | `InstallPlan` has been created (referenced in `status.installplan`) to install a new CSV                   |
| AtLatestKnown    | `status.installedCSV` matches the latest available CSV in catalog                                             |

## Catalog (Registry) Design

The Catalog Registry stores CSVs and CRDs for creation in a cluster, and stores metadata about packages and channels.

A package manifest is an entry in the catalog registry that associates a package identity with sets of ClusterServiceVersions. Within a package, channels point to a particular CSV. Because CSVs explicitly reference the CSV that they replace, a package manifest provides the catalog Operator all of the information that is required to update a CSV to the latest version in a channel (stepping through each intermediate version).

```
Package {name}
  |
  +-- Channel {name} --> CSV {version} (--> CSV {version - 1} --> ...)
  |
  +-- Channel {name} --> CSV {version}
  |
  +-- Channel {name} --> CSV {version}
```
