# asg-ip-obtainer

An Lambda implementation in Python and Go, demonstrating a difference between Boto3 and Go SDK

## What are we doing here

This repo contains 2 serverless applications (AWS API Gateway + AWS Lambda) serving a single purpose - get a list of IP addresses of the EC2 instances, which reside in an AutoScaling group.
The development of the serverless app in Go was a part of my stream on Twitch (LINK_PLACEHOLDER)

## How to get it running

### Prerequisites

You will need to have the following:
1. [AWS CLI](https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html)
2. [AWS SAM](https://github.com/awslabs/aws-sam-cli)
3. AWS Account

### Creating the ASGs
So you don't spend time creating ASGs, just provision the stack using the template in the repository:
```bash
aws cloudformation deploy --stack-name asg --template-file asg.yaml
```
You can override the parameter by using `--parameter-overrides` argument
> Mind the AWS costs, when creating stacks! Don't forget to clean resources afterwards!