/*
 * @Date: 2021-12-23 14:10:57
 * @LastEditors: sunshouxun
 * @LastEditTime: 2021-12-26 15:36:26
 * @FilePath: /ohmytiup/components/cloud-cluster/main.go
 */
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
	privateKey, publicKey, err := utils.GenerateSSHKeys()

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
			Region:     aws.String("us-east-2"),
			MaxRetries: aws.Int(3),
		}))
	// Create EC2 service client
	svc := ec2.New(sess)
	bcs := base64.StdEncoding.EncodeToString([]byte(initScript))
	// Specify the details of the instance that you want to create.
	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		ImageId:          aws.String("ami-002068ed284fb165b"),
		InstanceType:     aws.String("t2.micro"),
		MinCount:         aws.Int64(1),
		MaxCount:         aws.Int64(1),
		UserData:         aws.String(bcs),
		SubnetId:         aws.String("subnet-3a482e76"),
		SecurityGroupIds: []*string{aws.String("sg-0f11e4ceb2270d1bf")},
	})

	if err != nil {
		fmt.Println("Could not create instance", err)
		return
	}
	// svc.DescribeInstances(&ec2.DescribeInstancesInput{})

	//out, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
	//	InstanceIds: []*string{aws.String("10.0.0.6")},
	//})

	fmt.Println("Created instance", *runResult.Instances[0].InstanceId)
	fmt.Printf("privateKey, %s", privateKey)
}
