Palo Alto Networks cloudngfw
============================

[![GoDoc](https://godoc.org/github.com/PaloAltoNetworks/cloud-ngfw-aws-go?status.svg)](https://godoc.org/github.com/PaloAltoNetworks/cloud-ngfw-aws-go)

Package cloudngfw is a golang SDK for interacting with the Cloud NGFW AWS API.

Setup
=====

This uses the AWS golang SDK under the hood, so it assumes you have
credentials stored in the [standard
spots](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html).


Example Script
==============

```go
package main

import (
    "context"
    "log"

    "github.com/paloaltonetworks/cloud-ngfw-aws-go"
)

func main() {
    var err error

    c := &awsngfw.Client{
        Host: "api.endpoint.com",
        Region: "us-east-1",
        LfaArn: "arn:aws:iam::123456789:role/CloudNgfwFirewallAdmin",
        LraArn: "arn:aws:iam::123456789:role/CloudNgfwRulestackAdmin",
    }
    if err = c.Setup(); err != nil {
        log.Fatal(err)
    }

    if err = c.RefreshJwts(context.TODO()); err != nil {
		log.Fatal(err)
	}

    log.Printf("Firewall JWT: %s", c.FirewallJwt)
    log.Printf("Rulestack JWT: %s", c.RulestackJwt)
}
```
