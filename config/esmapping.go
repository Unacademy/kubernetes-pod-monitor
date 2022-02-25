package config

const Mapping = `
{
	"mappings": {
		"%s": {
			"properties": {
				"namespace": {
					"type":"keyword"
				},
				"pod_name": {
					"type":"keyword"
				},
				"container_name": {
					"type":"keyword"
				},
				"cluster_name": {
					"type":"keyword"
				},
				"created_at": {
					"type":   "date", 
					"format": "strict_date_optional_time||epoch_millis"
				},
				"restart_count": {
					"type": "long"
				},
				"logs": {
					"type": "text",
      				"index": false
				},
				"termination_state": {
					"type": "text"
				}
			}
		}
	}
}`

const Mappingv7 = `
{
	"mappings": {
		"properties": {
			"namespace": {
				"type":"keyword"
			},
			"pod_name": {
				"type":"keyword"
			},
			"container_name": {
				"type":"keyword"
			},
			"cluster_name": {
				"type":"keyword"
			},
			"created_at": {
				"type":   "date", 
				"format": "strict_date_optional_time||epoch_millis"
			},
			"restart_count": {
				"type": "long"
			},
			"logs": {
				"type": "text",
				"index": false
			},
			"termination_state": {
				"type": "text"
			}
		}
	}
}`
