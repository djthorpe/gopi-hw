/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package hw

import (
	"context"
	"fmt"
	"os"
	"sync"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi-rpc/sys/grpc"

	// Protocol buffers
	pb "github.com/djthorpe/gopi-hw/rpc/protobuf/hw"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Service struct {
	Server   gopi.RPCServer
	Hardware gopi.Hardware
}

type service struct {
	log      gopi.Logger
	hardware gopi.Hardware
	sync.Mutex
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open the server
func (config Service) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<grpc.service.hw>Open{ server=%v hardware=%v }", config.Server, config.Hardware)

	// Check for bad input parameters
	if config.Server == nil || config.Hardware == nil {
		return nil, gopi.ErrBadParameter
	}

	this := new(service)
	this.log = log
	this.hardware = config.Hardware

	// Register service with GRPC server
	pb.RegisterHardwareServer(config.Server.(grpc.GRPCServer).GRPCServer(), this)

	// Success
	return this, nil
}

func (this *service) Close() error {
	this.log.Debug("<grpc.service.hw>Close{ hardware=%v }", this.hardware)

	// Release resources
	this.hardware = nil

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *service) String() string {
	return fmt.Sprintf("grpc.service.hw{ hardware=%v }", this.hardware)
}

////////////////////////////////////////////////////////////////////////////////
// PING

func (this *service) Ping(ctx context.Context, _ *pb.EmptyRequest) (*pb.EmptyReply, error) {
	this.log.Debug2("<grpc.service.hw>Ping{ }")
	return &pb.EmptyReply{}, nil
}

////////////////////////////////////////////////////////////////////////////////
// INFO

func (this *service) Info(ctx context.Context, _ *pb.EmptyRequest) (*pb.InfoReply, error) {
	this.log.Debug2("<grpc.service.hw>Info{ }")
	this.Lock()
	defer this.Unlock()

	// Get the hostname
	if hostname, err := os.Hostname(); err != nil {
		return nil, err
	} else {
		// Return the hardware information
		return &pb.InfoReply{
			Name:             this.hardware.Name(),
			SerialNumber:     this.hardware.SerialNumber(),
			NumberOfDisplays: uint32(this.hardware.NumberOfDisplays()),
			Hostname:         hostname,
		}, nil
	}
}
