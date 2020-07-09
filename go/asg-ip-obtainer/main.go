/*
Almighty ASG IP obtainer!
I do the following:
1. I connect to AWS API
2. I read all the ASGs
3. I get all the instances in these ASGs!
4. I get the private IP address from all the instances from all the ASGs!
*/
package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	ips := make([]string, 1, 10)
	// Init AWS session
	sess := session.Must(session.NewSession())
	// Init ASG client
	asgService := autoscaling.New(sess)
	// Init EC2 client
	ec2Service := ec2.New(sess)

	// Give me all the autoscaling groups!
	// I don't filter ASGs by names, so I pass empty Input structure
	asgOutput, err := asgService.DescribeAutoScalingGroups(&autoscaling.DescribeAutoScalingGroupsInput{})
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// For each ASG in the output
	for _, asg := range asgOutput.AutoScalingGroups {
		// For each instance in an ASG
		for _, instance := range asg.Instances {
			// Give me the instance description!
			instanceInfo, err := ec2Service.DescribeInstances(&ec2.DescribeInstancesInput{
				InstanceIds: []*string{
					// Here I pass the instanceId, which is stored in the Instance structure.
					// InstanceIds field expects a slice of string
					// So i use aws.String function to convert the string to the pointer
					// to this string
					// Ye, Go is amazing
					aws.String(*instance.InstanceId),
				},
			})
			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}
			// Instance description doesn't return instance struct,
			// but instance description output struct
			// This struct contains a slice of Resevations (huh???)
			// Which i need to dig in
			// So I can get the Instance struct
			for _, reserv := range instanceInfo.Reservations {
				for _, instInReserv := range reserv.Instances {
					// Once I found what I need
					// I append the ips slice with an IP address
					ips = append(ips, *aws.String(*instInReserv.PrivateIpAddress))
				}
			}
		}
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("%v", ips),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
