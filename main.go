//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs --provider-name=NexGenCloud/terraform-provider-hyperstack --providers-schema=artifacts/provider-spec.json

package main

import (
	"context"
	"flag"
	"log"

	"github.com/NexGenCloud/terraform-provider-hyperstack/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	version         string = ""
	providerAddress string = ""
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider in debug mode")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: providerAddress,
		Debug:   debug,
	}

	err := providerserver.Serve(
		context.Background(),
		provider.New(version),
		opts,
	)

	if err != nil {
		log.Fatal(err.Error())
	}
}
