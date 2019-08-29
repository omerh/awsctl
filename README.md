# awsctl

This is a small maintanence tool for keeping cloudwatch expiery taht by default are set to Never to 14 days

>Most of the time its just forgeting to place the expiery.
>So just keep the dollars in yout pocket

To build

```bash
go get github.com/aws/aws-sdk-go
go get github.com/spf13/cobra
go build -ldflags "-s -w"
```
