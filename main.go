package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/jakgibb/check-cloudwatch/response"
	"gopkg.in/ini.v1"
)

func main() {
	// Policy required to read alarm information: cloudwatch:DescribeAlarms
	// Credentials file must contain access_key and secret_key (see: https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials)
	// Profile specifies which profile within the credentials file to use
	credFile := flag.String("credentials", "", "Path to the file containing the AWS access keys")
	profile := flag.String("profile", "", "The profile within the credentials file to use")
	region := flag.String("region", "", "Region to perform the check in")
	alarm := flag.String("alarm", "", "Name of the CloudWatch alarm to check")

	flag.Parse()

	if *credFile == "" || *profile == "" || *region == "" || *alarm == "" {
		response.Unknown("Required options not specified: [credentials profile region alarm]").Exit()
	}

	// Load the credentials from the credentials file
	// The SDK can load credentials automatically (using NewSessionWithOptions) from the ~/.aws/credentials file, however,
	// due to the way Nagios handles arguments, this prevents this method from working, hence loading the credentials manually
	credCfg, err := ini.Load(*credFile)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	// Retrieve the aws_access_key_id and aws_secret_access_key from the credentials file
	akid := credCfg.Section(*profile).Key("aws_access_key_id").String()
	skey := credCfg.Section(*profile).Key("aws_secret_access_key").String()

	// DescribeAlarmsInput struct provides filtering of CloudWatch alarms by name
	// Alarm names are specified as a slice of string pointers (Nagios will only ever pass a single alarm name)
	var inp cloudwatch.DescribeAlarmsInput
	var alarmNames []*string
	alarmNames = append(alarmNames, alarm)
	inp.AlarmNames = alarmNames

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(*region),
		Credentials: credentials.NewStaticCredentials(akid, skey, ""),
	})
	svc := cloudwatch.New(sess)

	resp, err := svc.DescribeAlarms(&inp)
	if err != nil {
		response.Unknown("error retrieving alarm data: " + err.Error()).Exit()
	}

	if len(resp.MetricAlarms) < 1 {
		response.Unknown("alarm not found - check region and alarm name").Exit()
	}

	var output string
	var state string

	for _, alarm := range resp.MetricAlarms {
		output = *alarm.StateReason
		state = *alarm.StateValue
	}

	switch state {
	case "OK":
		response.Ok(output).Exit()
	case "ALARM":
		response.Critical(output).Exit()
	default:
		response.Unknown(output).Exit()
	}

}
