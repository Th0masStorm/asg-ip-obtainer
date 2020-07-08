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

func handler(/*request events.APIGatewayProxyRequest*/) (/*events.APIGatewayProxyResponse, error*/ []string, error) {
	ips := make([]string, 1, 10)
	// Init AWS session
	sess := session.Must(session.NewSession())
	// Init ASG client
	asgService := autoscaling.New(sess)
	// Init EC2 client
	ec2Service := ec2.New(sess)

	asgInput := autoscaling.DescribeAutoScalingGroupsInput{}
	asgOutput, err := asgService.DescribeAutoScalingGroups(&asgInput)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	for _, asg := range asgOutput.AutoScalingGroups {
		for _, instance := range asg.Instances {

			//instance.InstanceId
			input := &ec2.DescribeInstancesInput{
				InstanceIds: []*string{
					aws.String(fmt.Sprintf("%v", instance.InstanceId)),
				},
			}

			instanceInfo, err := ec2Service.DescribeInstances(input)
			if err != nil {
				return events.APIGatewayProxyResponse{}, err
			}
			for _, reserv := range instanceInfo.Reservations {
				for _, instInReserv := range reserv.Instances {
					ips = append(ips, fmt.Sprintf("%v", instInReserv.PrivateIpAddress))
				}
			}
			// instance -> output struct -> reserv slice -> reserv struct -> slice instances -> instance -> privateip
		}

	}

	//return events.APIGatewayProxyResponse{
	//	Body:       fmt.Sprintf("%v", ips),
	//	StatusCode: 200,
	//}, nil
	return ips, nil
}

func main() {
	//lambda.Start(handler)
	stuff, err := handler()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(stuff)
}
