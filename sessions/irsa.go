package sessions

import (
	"os"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

func checkIRSAAvailable() (bool, string, string) {
	irsa := true
	roleARN, exists := os.LookupEnv("AWS_ROLE_ARN")
	if !exists {
		irsa = false
	}
	tokenPath, exists := os.LookupEnv("AWS_WEB_IDENTITY_TOKEN_FILE")
	if !exists {
		irsa = false
	}
	return irsa, roleARN, tokenPath
}

func GetAWSChainCredentialsV1(svc *sts.STS, sess *session.Session) *credentials.Credentials {
	irsa, roleARN, tokenPath := checkIRSAAvailable()

	chain := []credentials.Provider{
		&credentials.EnvProvider{},
		&ec2rolecreds.EC2RoleProvider{
			Client: ec2metadata.New(sess),
		},
		&credentials.SharedCredentialsProvider{},
	}

	if irsa {
		chain = append(chain, stscreds.NewWebIdentityRoleProvider(svc, roleARN, "", tokenPath))
	}

	creds := credentials.NewChainCredentials(chain)
	creds.Get() //IMPORTANT DO NOT remove otherwise throws error: RequestCanceled: request context canceled
	//caused by: context deadline exceeded: no Elasticsearch node available

	return creds
}
