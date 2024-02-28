package main

import (
	"context"
	"log"

	"github.com/teslamotors/fleet-telemetry/examples/emulator"
)

var (
	clientCert = "/Users/apatro/go_workspace/src/github.com/fleet-telemetry/test/integration/test-certs/vehicle_device.device-1.cert"
	clientKey  = "/Users/apatro/go_workspace/src/github.com/fleet-telemetry/test/integration/test-certs/vehicle_device.device-1.key"
	caClient   = "/Users/apatro/go_workspace/src/github.com/fleet-telemetry/test/integration/test-certs/vehicle_device.CA.cert"
)

func main() {

	mockVehicle, err := emulator.NewMockVehicle("device-1", clientCert, clientKey, caClient, "/Users/apatro/go_workspace/src/github.com/fleet-telemetry/examples/streaming_config.json")
	if err != nil {
		log.Fatalf("Cannot instantiate vehicle %v", err)
	}
	mockVehicle.EnableBackgroundTransmission(context.Background())
}
