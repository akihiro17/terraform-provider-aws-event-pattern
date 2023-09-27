data "aws-event-pattern_test" "ec2_state_change" {
  event_pattern = jsonencode(
    {
      "source" : ["aws.ec2"],
      "detail-type" : ["EC2 Instance State-change Notification"],
      "detail" : {
        "state" : ["terminated"]
      }
  })

  event = jsonencode({
    "version" : "0",
    "id" : "6a7e8feb-b491-4cf7-a9f1-bf3703467718",
    "detail-type" : "EC2 Instance State-change Notification",
    "source" : "aws.ec2",
    "account" : "111122223333",
    "time" : "2017-12-22T18:43:48Z",
    "region" : "us-west-1",
    "resources" : [
      "arn:aws:ec2:us-west-1:123456789012:instance/i-1234567890abcdef0"
    ],
    "detail" : {
      "instance-id" : "i-1234567890abcdef0",
      "state" : "terminated"
    }
  })

  lifecycle {
    postcondition {
      condition     = self.match == true
      error_message = <<-EOS
      ${self.event_pattern}
        should match
      ${self.event}
      EOS
    }
  }
}
