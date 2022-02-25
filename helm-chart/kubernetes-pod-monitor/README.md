# Kubernetes Pod Monitor helm chart
Helm chart for the Kubernetes Pod Monitor project, which is created to monitor Kubernetes pod crashes.

To install via helm 3, run the following commands:

```
helm repo add kubernetes-pod-monitor https://unacademy.github.io/kubernetes-pod-monitor/
helm upgrade -i --create-namespace kubernetes-pod-monitor kubernetes-pod-monitor/kubernetes-pod-monitor --namespace kubernetes-pod-monitor
```

<br/><br/>
<a name="config-options"></a>
The following table lists commonly used configuration parameters for the Kubernetes Pod Monitor Helm chart and their default values.

Parameter | Description | Default
--------- | ----------- | -------
`serviceAccount.create` | Set this to `false` if you want to create the service account `kubecost-cost-analyzer` on your own | `true`
`tolerations` | node taints to tolerate | `[]`
`env` | environment variables | `[]`
`affinity` | pod affinity | `{}`
`resources` | resource requests and limits | `{}`
`config.deployEnv` | Set this to `local` if you want to connect using the kubeconfig file named `<CLUSTER_NAME>.yml` | `release`
`config.clusterName` | name of kubernetes cluster | `""`
`config.aws.region` | AWS region name | `""`
`config.elasticsearch.url` | Elasticsearch endpoint | `https://127.0.0.1`
`config.elasticsearch.dashboard` | Elasticsearch dashboard (if any) that is sent in slack notification | `""`
`config.elasticsearch.port` | Elasticsearch port | `443`
`config.elasticsearch.scheme` | Elasticsearch port | `https`
`config.elasticsearch.v7` | Set this to `false` if Elasticsearch version is less than 7 | `true`
`config.slack.notify` | Set this to `false` to disable slack notifications | `true`
`config.slack.channel` | Default slack channel to send notifications | `pod-crash-alerts`
`config.slack.token` | Slack token | `""`
`config.sql.host` | MySQL hostname | `127.0.0.1`
`config.sql.port` | MySQL port | `3306`
`config.sql.username` | MySQL username | `admin`
`config.sql.password` | MySQL password | `admin`
`config.sql.dbname` | MySQL database name | `kubernetes_pod_monitor`
