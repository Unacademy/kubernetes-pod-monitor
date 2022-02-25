package sessions

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/olivere/elastic"
	"github.com/sha1sum/aws_signing_client"
	"github.com/spf13/viper"
)

var (
	esClient *elastic.Client
)

func NewElasticSearchLocalClient(esURL string) (*elastic.Client, error) {
	return elastic.NewClient(
		elastic.SetURL(esURL),
		elastic.SetSniff(false),
	)
}

func NewElasticSearchAwsClient(esURL string, awsRegion string) (*elastic.Client, error) {

	sess := session.Must(session.NewSession(aws.NewConfig().WithRegion(awsRegion)))
	svc := sts.New(sess)
	creds := GetAWSChainCredentialsV1(svc, sess)
	signer := v4.NewSigner(creds)

	awsClient, err := aws_signing_client.New(signer, nil, "es", awsRegion)
	if err != nil {
		return nil, err
	}

	return elastic.NewSimpleClient(
		elastic.SetURL(esURL),
		elastic.SetScheme(viper.GetString("elasticsearch.scheme")),
		elastic.SetSniff(false),
		elastic.SetHttpClient(awsClient),
	)
}

// InitElasticsearchClient initializes the Elastic Search Client
// The function panics if there is any error making connection
// with the elastic search cluster.
func InitElasticsearchClient() {
	if esClient != nil {
		return
	}

	url := fmt.Sprintf("%s:%s", viper.GetString("elasticsearch.url"),
		viper.GetString("elasticsearch.port"))
	env := viper.GetString("DEPLOY_ENV")
	if env == "local" {
		client, err := NewElasticSearchLocalClient(url)
		if err != nil {
			panic(err)
		}
		esClient = client
	} else {
		awsRegion := viper.GetString("aws.region")
		client, err := NewElasticSearchAwsClient(url, awsRegion)
		if err != nil {
			panic(err)
		}
		esClient = client
	}
}

// GetElasticSearchClient returns the instance of Elastic Search Client that have
// already been initialized through InitElasticsearchClient
func GetElasticSearchClient() *elastic.Client {
	InitElasticsearchClient()
	return esClient
}
