package main

import (
	"os"

	// Frameworks
	"github.com/djthorpe/gopi"

	// Modules
	_ "github.com/djthorpe/gopi-hw/sys/hw"
	_ "github.com/djthorpe/gopi-hw/sys/metrics"
	_ "github.com/djthorpe/gopi-rpc/sys/grpc"
	_ "github.com/djthorpe/gopi/sys/logger"

	// RPC Services
	_ "github.com/djthorpe/gopi-hw/rpc/grpc/hw"
	_ "github.com/djthorpe/gopi-hw/rpc/grpc/metrics"
)

///////////////////////////////////////////////////////////////////////////////

func main() {
	// Create the configuration
	config := gopi.NewAppConfig("rpc/service/metrics", "rpc/service/hw")

	// Set the RPCServiceRecord for server discovery
	config.Service = "hw"

	// Run the server and register all the services
	os.Exit(gopi.RPCServerTool(config))
}
