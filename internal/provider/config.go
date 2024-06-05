package provider

import (
	"fmt"
)

var (
	Version         = "0.0.0"
	CommitHash      = "n/a"
	CommitTimestamp = "n/a"
	BuildTimestamp  = "n/a"
)

var (
	ProviderName = ""
	EnvPrefix    = ""
	ApiAddress   = ""
)

func BuildVersion() string {
	return fmt.Sprintf("%s-%s (%s)", Version, CommitHash, BuildTimestamp)
}
