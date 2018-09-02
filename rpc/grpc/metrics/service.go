/*
	Go Language Raspberry Pi Interface
	(c) Copyright David Thorpe 2016-2018
	All Rights Reserved

	Documentation http://djthorpe.github.io/gopi/
	For Licensing and Usage information, please see LICENSE.md
*/

package metrics

import (
	"context"
	"fmt"

	// Frameworks
	"github.com/djthorpe/gopi"
	"github.com/djthorpe/gopi-rpc/sys/grpc"

	// Protocol buffers
	pb "github.com/djthorpe/gopi-hw/rpc/protobuf/metrics"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type Service struct {
	Server  gopi.RPCServer
	Metrics gopi.Metrics
}

type service struct {
	log     gopi.Logger
	metrics gopi.Metrics
}

////////////////////////////////////////////////////////////////////////////////
// OPEN AND CLOSE

// Open the server
func (config Service) Open(log gopi.Logger) (gopi.Driver, error) {
	log.Debug("<grpc.service.metrics.Open>{ server=%v metrics=%v }", config.Server, config.Metrics)

	// Check for bad input parameters
	if config.Server == nil || config.Metrics == nil {
		return nil, gopi.ErrBadParameter
	}

	this := new(service)
	this.log = log
	this.metrics = config.Metrics

	// Register service with GRPC server
	pb.RegisterMetricsServer(config.Server.(grpc.GRPCServer).GRPCServer(), this)

	// Success
	return this, nil
}

func (this *service) Close() error {
	this.log.Debug("<grpc.service.metrics.Close>{ metrics=%v }", this.metrics)

	// Release resources
	this.metrics = nil

	// Success
	return nil
}

////////////////////////////////////////////////////////////////////////////////
// Stringify

func (this *service) String() string {
	return fmt.Sprintf("grpc.service.metrics{}")
}

////////////////////////////////////////////////////////////////////////////////
// Ping method

func (this *service) Ping(ctx context.Context, _ *pb.EmptyRequest) (*pb.EmptyReply, error) {
	this.log.Debug2("<grpc.service.metrics>Ping{ }")
	return &pb.EmptyReply{}, nil
}
