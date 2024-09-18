package aws

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

// Caller defines an interface to get the AWS caller name
type Caller interface {
	GetCallerName() (string, error)
}

// RealAWSCaller is the actual implementation of AWSCaller
type RealAWSCaller struct{}

func NewAWSCaller() Caller {
	return &RealAWSCaller{}
}

// GetCallerName returns the name of the AWS caller
func (c *RealAWSCaller) GetCallerName() (string, error) {
	sess := session.Must(session.NewSession())
	svc := sts.New(sess)

	result, err := svc.GetCallerIdentity(nil)
	if err != nil {
		return "", fmt.Errorf("error getting AWS caller identity: %v", err)
	}

	arnParts := strings.Split(*result.Arn, "/")
	if len(arnParts) > 1 {
		return arnParts[1], nil
	}
	return "", nil
}
