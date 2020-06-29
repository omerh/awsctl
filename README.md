# awsctl

This is a small maintanence tool for managing aws infrastructure easily with a single binary on a region or all regions at a single command

Tool is built using cobra, for getting started just run `awsctl` and see the example commands available.

To use the tool with dry run just run the command, the really execute add `--yes`

Optional output as json using `-o json`

WIP: Optionl slack message using `--slack` and setting environment variable `AWSCTL_SLACK_URL`

## build

```bash
# Use go modules add the following env var GO111MODULE=on
go build -ldflags "-s -w"
```

### Example commands

Get all EC2 events from all regions

```bash
awsctl get ec2 events -r all
```

List regions

```bash
awsctl list regions
```

List availablity zones in a region

```bash
awsctl list azs --region us-east-1
```

Delete all unused EBS in all regions

```bash
awsctl delete ebs --region all --yes
```

Set cloudwatch logs with no expirey to 14 days expiry

```bash
awsctl set cloudwatch --region all --retention 14 --yes
```

Manage RDS insatnces or clusters

```bash
awsctl get rds --region all  --type instance #or cluster
awsctl get rdssnapshots --region all  --type instance --name db01
awsctl delete rdssnapshots --name db01 --type instance --region all --older 14 --yes
```

ECR Opertaion for setting lifecycle policy to untagged repositories

```bash
awsctl set ecrregistrypolicy -r eu-west-2 --retention 7
```

ECR Repository configuration to scanOnPush repository for vulnerabilities

```bash
awsctl set ecrscanonpush --region eu-west-1 --scan true  --yes
awsctl set ecrscanonpush --region all --scan true  --yes
```

ACM Certificates

```bash
awsctl get certificates --region ue-east-1    # Get all expiring and expired certificates in region or all regions, for all expiring certificates it analyses why aren't the certificates being renewed automatically
awsctl delete certificates --region all --yes # Delete all unused certificates from the account
```

Cloudwatch Alarms

>Currently tested on `Errors` metric for lambda

```bash
# Single lambda
awsctl set cloudwatchalarm --resource lambda --metric errors --region eu-west-2 --arn arn:aws:lambda:eu-west-2:000000000000:function:test --threshold 3 --action arn:aws:sns:eu-west-2:000000000000:SNSToSlack --yes

# All lambdas
awsctl set cloudwatchalarm --resource lambda --metric errors --region eu-west-2 --threshold 3 --action arn:aws:sns:eu-west-2:000000000000:SNSToSlack --yes

# All lambdas in all regions
awsctl set cloudwatchalarm --resource lambda --metric errors --region all --threshold 3 --action arn:aws:sns:eu-west-2:000000000000:SNSToSlack --yes
```

Delete network interfaces

```bash
awsctl delete ni --region eu-west-2 --filter available --yes
```

For any missing action please open an issue for a feature request.

### Contributing

Fork, implement, add tests, pull request, get my everlasting thanks and a respectable place here :).

### Copyright

Copyright (c) 2019 Omer Haim, [@omerhaim](http://twitter.com/omerhaim).
See [LICENSE](LICENSE) for further details.
