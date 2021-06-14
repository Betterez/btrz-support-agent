module github.com/betterez/btrz-support-agent/agent

go 1.15

replace github.com/betterez/btrz-support-agent/utils => ../utils

require (
	github.com/aws/aws-sdk-go v1.38.60
	github.com/betterez/btrz-support-agent/utils v0.0.0-00010101000000-000000000000
	github.com/bsphere/le_go v0.0.0-20200109081728-fc06dab2caa8
)
