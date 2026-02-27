package rpc

import (
	"github.com/1111mp/gin-app/internal/rpc/grpc"
	"github.com/1111mp/gin-app/internal/rpc/rmqrpc"
	"go.uber.org/fx"
)

var Module = fx.Module(
	"rpc",

	// rabbitmq rpc
	rmqrpc.Module,
	// grpc
	grpc.Module,
)
