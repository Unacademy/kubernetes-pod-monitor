## Kubernetes Pod Monitor

Kubernetes Pod Monitor actively tracks your K8S pods and alerts container restarts alongwith its crash logs thereby decreasing mean time to detect (MTTD). The features include:

- Alerting using slack integration
- Capturing critical crash logs and storing in Elasticsearch
- Historical pod crashes
- Storing container state that gives transparency on pod lifetime and status before termination
- Kibana Visualization for filtering through crashes
- Ability to configure slack channel based on namespace
- Ability to ignore certain namespaces


![Elasticsearch Dashboard](getting-started/dashboard.jpeg)

## Requirements

<a name="requirements"></a>
The following table lists the minimum requirements for running Kubernetes Pod Monitor.

Tool | Minimum version | Minimum configuration
--------- | ----------- | -------
Kubernetes | 1.13 | 100 MB RAM
MySQL | 5.7 | `-`
Elasticsearch | 6.5 | 4 GB RAM

To send alerts via Slack integration, access tokens can be generated here: https://api.slack.com/authentication/token-types

## Getting Started

You can deploy Kubernetes Pod Monitor on any Kubernetes 1.13+ cluster in a matter of minutes, if not seconds.
### Using Helm chart (recommended)
  - [Apply MySQL migrations](#mysql-migrations)
  - [Install using the Helm chart](helm-chart/kubernetes-pod-monitor/README.md)
  - Import [Kibana dashboard](getting-started/es_saved_objects.json) into Elasticsearch by following https://www.elastic.co/guide/en/kibana/current/managing-saved-objects.html

### Using docker compose
  - Add kuberentes configuration file to `config` directory and update `CLUSTER_NAME` env variable in docker-compose
  - Start docker compose using:
  
    ```sh
    docker-compose up
    ```


## MySQL Migrations

You can run the following queries to create the required database and tables:

```sql
CREATE DATABASE kubernetes_pod_monitor
```

```sql
CREATE TABLE `k8s_crash_monitor` (
`clustername` char(64) NOT NULL,
`namespace` char(64) NOT NULL,
`podname` char(255) NOT NULL,
`containername` char(255) NOT NULL,
`restartcount` int(11) DEFAULT NULL,
`retries` int(11) DEFAULT NULL,
`edited_at` int(11) DEFAULT NULL,
PRIMARY KEY (`clustername`,`namespace`,`podname`,`containername`)
);
```

```sql
CREATE TABLE `k8s_pod_crash` (
`id` int(11) NOT NULL AUTO_INCREMENT,
`clustername` varchar(120) NOT NULL,
`namespace` varchar(120) NOT NULL,
`containername` varchar(120) NOT NULL,
`restartcount` int(11) NOT NULL DEFAULT '0',
`date` datetime(6) DEFAULT NULL,
PRIMARY KEY (`id`)
);
```

```sql
CREATE TABLE `k8s_pod_crash_notify` (
`clustername` varchar(255) NOT NULL,
`namespace` varchar(255) NOT NULL,
`slack_channel` varchar(255) NOT NULL,
PRIMARY KEY (`clustername`,`namespace`)
);
```

```sql
CREATE TABLE `k8s_crash_ignore_notify` (
`clustername` varchar(255) NOT NULL,
`namespace` varchar(255) NOT NULL,
`containername` varchar(255) NOT NULL,
PRIMARY KEY (`clustername`,`namespace`,`containername`)
);
```

## Configuring notifications

You can easily configure slack notifications, by using the [notification management utility](scripts/notification_management_utility.py). 

The following lists the minimum requirements for running this utility:
- Python v3.6 or higher
- PyMSQL package to manage MySQL tables: https://pypi.org/project/PyMySQL/
  ```sh
  pip3 install PyMySQL
  ```
- Tabulate package to render tables: https://pypi.org/project/tabulate/
  ```sh
  pip3 install tabulate
  ```

Run the utility and follow the onscreen steps:

```sh
python3 scripts/notification_management_utility.py
```

## Sample Elasticsearch document
An indexed document in Elasticsearch consists of following fields:
  - `namespace`: Namespace of the crashed pod
  - `pod_name`: Name of the pod that crashed
  - `container_name`: Container name which restarted. Helpful incase of multiple containers in a pod
  - `created_at`: Timestamp in milliseconds
  - `cluster_name`: Name of the cluster
  - `logs`: Logs of the container before restarting
  - `restart_count`: Number of times the pod restarted
  - `termination_state`: State of the container with reason, message, started at timestamp and finished at timestamp

```json
{
  "_index": "k8s-crash-monitor-2022.03.11",
  "_type": "_doc",
  "_id": "Zn3DeH8BpsFVE9gY0heI",
  "_version": 1,
  "_score": null,
  "_source": {
    "namespace": "prometheus",
    "pod_name": "prometheus-server-68bf5b8675-bxpq6",
    "container_name": "prometheus-server",
    "created_at": 1646998573563,
    "cluster_name": "dev-001",
    "logs": "level=error ts=2022-03-11T11:35:53.889Z caller=main.go:723 err=\"opening storage failed: zero-pad torn page: write /data/wal/00000269: no space left on device\"\n",
    "restart_count": 183,
    "termination_state": "&ContainerStateTerminated{ExitCode:1,Signal:0,Reason:Error,Message:,StartedAt:2022-03-11 11:35:53 +0000 UTC,FinishedAt:2022-03-11 11:35:53 +0000 UTC,ContainerID:docker://3cc68f0bdff60e4ac3ab494235225af22bfa3efa97ab5ea55464fcb510dbb0f6,}"
  },
  "fields": {
    "created_at": [
      "2022-03-11T11:36:13.563Z"
    ]
  },
  "sort": [
    1646998573563
  ]
}
```


## Software stack

Golang application. 
Kubernetes.
Elasticsearch.
MySQL.

## Contributors
<table>
  <tr>
    <td align="center"><a href="https://www.linkedin.com/in/shivam-gupta-dtu/"><img src="https://avatars1.githubusercontent.com/u/22556869?s=460&u=bd28a7d3ffa18bf409071ae6c9eae80692d0143e&v=4" width="100px;" alt="https://github.com/Shivam9268"/><br /><sub><b>Shivam Gupta</b></sub></a><br /></td>
    </tr>
</table>
