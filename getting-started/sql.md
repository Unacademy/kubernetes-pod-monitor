## MySQL queries

```sql
CREATE DATABASE kubernetes_pod_monitor
```
 
```sql
CREATE TABLE `k8s_crash_ignore_notify` (
`clustername` varchar(255) NOT NULL,
`namespace` varchar(255) NOT NULL,
`containername` varchar(255) NOT NULL,
PRIMARY KEY (`clustername`,`namespace`,`containername`)
);
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
