package main

import (
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(Provider())
}
