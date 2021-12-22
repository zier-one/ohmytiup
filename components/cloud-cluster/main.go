package main

import (
	"encoding/base64"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pingcap/tiup/components/cloud-cluster/utils"
)

func main() {
	_, publicKey, err := utils.GenerateSSHKeys()

	initScript := fmt.Sprintf(`#!/bin/bash
useradd -g root -d /home/tidb tidb -m
mkdir -p /tmp/aaa
mkdir -p /root/.ssh
echo "%s" > /root/.ssh/authorized_keys
mkdir -p /home/tidb/.ssh
echo "%s" > /home/tidb/.ssh/authorized_keys
chown -R tidb /home/tidb
mkdir -p /tmp/bbb
`, string(publicKey), string(publicKey))

	sess := session.Must(session.NewSession(
		&aws.Config{
			Region:     aws.String("ap-northeast-1"),
			MaxRetries: aws.Int(3),
		}))
	// Create EC2 service client
	svc := ec2.New(sess)
	bcs := base64.StdEncoding.EncodeToString([]byte(initScript))
	// Specify the details of the instance that you want to create.
	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		ImageId:          aws.String("ami-09f4ee26f4eca3293"),
		InstanceType:     aws.String("t2.micro"),
		MinCount:         aws.Int64(1),
		MaxCount:         aws.Int64(1),
		UserData:         aws.String(bcs),
		SubnetId:         aws.String("subnet-0b34922bcd1f4edef"),
		SecurityGroupIds: []*string{aws.String("sg-00ddc5815c26c2c90")},
	})

	if err != nil {
		fmt.Println("Could not create instance", err)
		return
	}

	//out, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
	//	InstanceIds: []*string{aws.String("10.0.0.6")},
	//})

	fmt.Println("Created instance", *runResult.Instances[0].InstanceId)
}
