# Cluster API usage

Hyperstack platform supports on-demand provisioning of managed kubernetes clusters. It works similar to cloud providers, you need to specify target kubernetes version, node type and basic parameters, everything else is handled by the platform.

## Endpoint breakdown

Note: all methods have standard wrapping:
* status / boolean / If call is successful
* message / string / Error message
* FIELD_NAME / object / Response (depends on call)

List of supported APIs:

* List supported k8s cluster versions:
  
  GET */core/clusters/versions*
  
  Response:
  * versions / list(string) - List of supported versions

* List clusters:
  
  GET */core/clusters*
  
  Response:
  * status / boolean - If success
  * message / string - Status message
  * cluster / list(cluster) - List of clusters
  
  Cluster field structure:
  * id / number - Cluster ID
  * name / string - Cluster name
  * environment name / string - Environment name
  * kubernetes_version / string - Cluster kubernetes version
  * kube_config / list(string) - kubeconfig for connection
  * status_reason / string - Error message?
  * node_count / number - Amount of VMs
  * node_flavor / object(flavor) - VM flavor fields
  * node_addresses / list(string) - node IP addresses
  * keypair_name / string - Name of keypair
  * enable_public_ip / boolean - If public IPs are enabled
  * created_at / datetime - Cluster creation time

* Create cluster:

  POST */core/clusters*

  Request:
    * environment_name / string - Name of Hyperstack environment
    * image_name / string - VM image name
    * keypair_name/ string - SSH keypair for VMs
    * kubernetes_version / string - Target kubernetes version
    * name / string  - Cluster name
    * node_count / number - Amount of VMs
    * node_flavor_name / string - Flavor name for VMs
    * region / string  - Target region
    * enable_public_ip / boolean - If public IP should be enabled

  Response:
    * id / number - Cluster ID
    * name / string - Cluster name
    * environment name / string - Environment name
    * kubernetes_version / string - Cluster kubernetes version
    * kube_config / list(string) - kubeconfig for connection
    * status_reason / string - Error message?
    * node_count / number - Amount of VMs
    * node_flavor / object(flavor) - VM flavor fields
    * node_addresses / list(string) - node IP addresses
    * keypair_name / string - Name of keypair
    * enable_public_ip / boolean - If public IPs are enabled
    * created_at / datetime - Cluster creation time

* Get cluster by ID

  GET */core/clusters/:id*

  Request:
    * id / number - Cluster ID

  Response:
    * id / number - Cluster ID
    * name / string - Cluster name
    * environment name / string - Environment name
    * kubernetes_version / string - Cluster kubernetes version
    * kube_config / list(string) - kubeconfig for connection
    * status_reason / string - Error message?
    * node_count / number - Amount of VMs
    * node_flavor / object(flavor) - VM flavor fields
    * node_addresses / list(string) - node IP addresses
    * keypair_name / string - Name of keypair
    * enable_public_ip / boolean - If public IPs are enabled
    * created_at / datetime - Cluster creation time

* Fetch cluster events by ID

  GET */core/clusters/:cluster_id/events*

  Request:
    * cluster_id / number - Cluster ID

  Response (list of events in *cluster_events*):
    * id / number - Event ID
    * cluster_id / number - Cluster ID
    * user_id / number - User ID
    * org_id / number - Org ID
    * time / dateTime - Event timestamp
    * type / string - Event type
    * reason / string - Event reason
    * object / string - Additional details
    * message / string - Event message

* Delete cluster by ID

  DELETE */core/clusters/:id*

  Request:
    * id / number - Cluster ID

  Response:
    * empty

## Examples

To provision a new k8s cluster:

* Start with fetching supported cluster versions:
  ````bash
  curl -X GET \
    "${HYPERSTACK_API_ADDRESS}/v1/core/clusters/versions" \
    -H "accept: application/json" \
    -H "api_key: ${HYPERSTACK_API_KEY}"
  ````
  
* Pick a version you want to use:
  ````bash
  K8S_VERSION=...
  ````
* Create a new environment *ENV_NAME*
* Create a new keypair *KEYPAIR_NAME*
* Pick the following parameters:
  * Image name for worker nodes
  * Cluster name
  * Amount of worker nodes
  * Worker flavor
  * Target region
  * Whenever you need public IPs
* Create a new cluster (example):
  ````bash
  curl -X POST \
    "${HYPERSTACK_API_ADDRESS}/v1/core/clusters" \
    -H "accept: application/json" \
    -H "api_key: ${HYPERSTACK_API_KEY}"
    --data '{
      "environment_name": "'"${ENV_NAME}"'",
      "image_name": "Ubuntu Server 22.04 LTS R535 CUDA 12.2",
      "keypair_name": "'"${KEYPAIR_NAME}"'",
      "kubernetes_version": "'"${K8S_VERSION}"'",
      "name": "test_cluster",
      "node_count": 3,
      "node_flavor_name": "n3-A100x1",
      "region": "poc-CANADA-1",
      "enable_public_ip": true
      }'
  ````
* Node cluster ID from the response *CLUSTER_ID*
* Test that cluster is present in the list and ready:
  ````bash
  curl -X GET \
    "${HYPERSTACK_API_ADDRESS}/v1/core/clusters/${CLUSTER_ID}" \
    -H "accept: application/json" \
    -H "api_key: ${HYPERSTACK_API_KEY}"
  ````
