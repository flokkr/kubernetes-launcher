# Kubernetes launcher

Kubernetes Launcher is a simple config downloader and process supervisor. From high level it's similar to [consul-template](github.com/hashicorp/consul-template):

1. It can download configmaps to a specific location.

2. It can start an application a child process.

3. In case of configmap change, the new files will be redownloaded and the process will be restarted.

## Configuration

 * `KUBERNETES_CONFIGMAP`: configmap to download (data should contain the file->content mapping)
 * `KUBERNETES_CONFIGMAP_NAMESPACE`: namespace of the configmap
 * `CONF_DIR`: destination directory

Supervised process should be defined as an argument. Such as:

```
./kubernetes-launcher hdfs namenode
```
