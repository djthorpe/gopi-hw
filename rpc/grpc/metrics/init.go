/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved
	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package metrics

import (
	// Frameworks
	"github.com/djthorpe/gopi"
)

////////////////////////////////////////////////////////////////////////////////
// INIT

func init() {
	gopi.RegisterModule(gopi.Module{
		Name:     "rpc/service/metrics",
		Type:     gopi.MODULE_TYPE_SERVICE,
		Requires: []string{"rpc/server", "metrics"},
		New: func(app *gopi.AppInstance) (gopi.Driver, error) {
			return gopi.Open(Service{
				Server:  app.ModuleInstance("rpc/server").(gopi.RPCServer),
				Metrics: app.ModuleInstance("metrics").(gopi.Metrics),
			}, app.Logger)
		},
	})
}
