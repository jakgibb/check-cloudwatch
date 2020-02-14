check-cloudwatch
=========

check-cloudwatch is a simple CLI program, utilising the AWS SDK, to check the status of a specified CloudWatch alarm.  
The program is designed to be used by Nagios and so will exit with the following exit codes, and echo a description of the alert.
* RC 0: OK
* RC 2: CRITICAL
* RC 3: UNKNOWN

### Setup

The `cloudwatch:DescribeAlarms` policy is required to read CloudWatch alarms

Access keys must be stored within an file, as per the format below.
The typical location to use is `~/.aws/credentials` or `%USERPROFILE%/.aws/credentials`   

````
[nagios-cloudwatch]
aws_access_key_id = ACCESSKEY
aws_secret_access_key = SECRETKEY
````

### Usage
Create a Nagios command to run the script (`go run`) or invoke the executable (if built with `go build`)   
````/usr/local/nagios/libexec/check_cloudwatch --credentials $ARG1$ --profile $ARG2$ --region $ARG3$ --alarm $ARG4$ ````

### Example
````
/usr/local/nagios/libexec/check_cloudwatch --credentials /home/nagios/.aws/credentials --profile nagios-cloudwatch --region eu-west-1 --alarm "CPU Usage"
CRITICAL: Threshold Crossed: 1 out of the last 1 datapoints [2.0 (10/02/20 14:34:00)] was greater than or equal to the threshold (0.01) (minimum 1 datapoint for OK -> ALARM transition).
exit status 2
````