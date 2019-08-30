# awsctl

This is a small maintanence tool for managing aws infrastructure easily with a single binary on a region or all regions at a single command

Tool is built using cobra, for getting started just run `awsctl` and see the examples available commands.

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

For any missing action please open an issue for a feature request.

### Contributing

Fork, implement, add tests, pull request, get my everlasting thanks and a respectable place here :).

### Copyright

Copyright (c) 2019 Omer Haim, [@omerhaim](http://twitter.com/omerhaim).
See [LICENSE](LICENSE) for further details.
