package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const testAccTestDataSourceConfig = `
data "aws-event-pattern_test" "ec2_state_change" {
  event_pattern = jsonencode(%s)
	event = jsonencode(%s)
}`

func TestAccTestDataSource(t *testing.T) {
	tests := []struct {
		name         string
		eventPattern string
		event        string
		match        bool
	}{
		{
			"match",
			`{
				"source" : ["aws.ec2"],
				"detail-type" : ["EC2 Instance State-change Notification"],
				"detail" : {
					"state" : ["terminated"]
				}
			}`,
			`{
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
			}`,
			true,
		},
		{
			"not match",
			`{
				"source" : ["aws.ec2"],
				"detail-type" : ["EC2 Instance State-change Notification"],
				"detail" : {
					"state" : ["terminated"]
				}
			}`,
			`{
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
					"state" : "running"
				}
			}`,
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { testAccPreCheck(t) },
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					// Read testing
					{
						Config: fmt.Sprintf(testAccTestDataSourceConfig, tt.eventPattern, tt.event),
						Check: resource.ComposeAggregateTestCheckFunc(
							resource.TestCheckResourceAttr("data.aws-event-pattern_test.ec2_state_change", "match", strconv.FormatBool(tt.match)),
						),
					},
				},
			})
		})
	}
}
