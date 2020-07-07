import json
from boto3 import Session


def lambda_handler(event, context):
    # Resulted list
    ips = []

    # Init Boto3 session
    sess = Session()
    asg_client = sess.client("autoscaling")
    ec2res = sess.resource("ec2")

    # Give me all the AutoScaling groups!
    for asg in asg_client.describe_auto_scaling_groups()['AutoScalingGroups']:
        # Give me all the instances in the AutoScaling group!
        for instance in asg['Instances']:
            # Is you healthy?
            if instance['HealthStatus'] == "Healthy":
                # Give me ze private address
                instance = ec2res.Instance(instance['InstanceId'])
                ips.append(instance.private_ip_address)


    return {
        "statusCode": 200,
        "body": json.dumps({
            "ips": ips,
        }),
    }
