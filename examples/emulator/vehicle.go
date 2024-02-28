package emulator

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/teslamotors/fleet-telemetry/messages/tesla"
	"github.com/teslamotors/fleet-telemetry/protos"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Producer handles dispatching data received from the vehicle
type Vehicle interface {
	EnableBackgroundTransmission(context context.Context)
	SendMessage(message []byte) error
	Vin() string
}

type MockVehicle struct {
	vin             string
	connection      *websocket.Conn
	logger          *logrus.Logger
	streamingConfig *StreamingConfig
	vehicleOnline   time.Duration
	counter         int64
}

func NewMockVehicle(vin, clienCertPath, clientKeyPath, caClientPath, streamingConfigPath string) (Vehicle, error) {
	streamingConfig, err := NewStreamingConfig(streamingConfigPath)
	if err != nil {
		return nil, err
	}
	tlsConfig, err := getTLSConfig(clienCertPath, clientKeyPath, streamingConfig.CA)
	if err != nil {
		return nil, err
	}
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339Nano})
	logger.SetLevel(logrus.DebugLevel)

	u := url.URL{Scheme: "wss", Host: "app:4443", Path: "/"}

	tlsDialer := websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: 45 * time.Second,
		TLSClientConfig:  tlsConfig,
	}
	connection, _, err := tlsDialer.Dial(u.String(), http.Header{})
	if err != nil {
		return nil, err
	}
	return &MockVehicle{
		vin:             vin,
		connection:      connection,
		logger:          logger,
		streamingConfig: streamingConfig,
		vehicleOnline:   30 * time.Second,
		counter:         0,
	}, nil

}

func (m *MockVehicle) Vin() string {
	return m.vin
}

func (m *MockVehicle) EnableBackgroundTransmission(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	timer := time.NewTimer(m.vehicleOnline)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				currentSecond := time.Now().Unix()
				var data []*protos.Datum

				for name, setting := range m.streamingConfig.Fields {
					frequency := setting.IntervalSeconds
					if currentSecond%int64(frequency) == 0 {
						val := GetFieldValue(name)
						if val != nil {
							data = append(data, val)
						}
					}
				}
				if len(data) > 0 {
					err := m.transmitTelemetryData(data)
					if err != nil {
						m.logger.Errorf("error transmitting data %v", err)
					} else {
						m.logger.Println("Data transmitted")
					}
				}

			case <-timer.C:
				m.logger.Println("Timer expired")
				return
			case <-ctx.Done():
				m.logger.Println("Context cancelled")
				return
			}
		}
	}()

	// Wait for either timer to expire or context to be cancelled
	select {
	case <-timer.C:
		m.logger.Println("Timer expired")
	case <-ctx.Done():
		m.logger.Println("Context cancelled")
	}
}

func (m *MockVehicle) EnableBackgroundTransmission1() {
	ticker := time.NewTicker(1 * time.Second)
	timer := time.NewTimer(m.vehicleOnline)
	go func() {
		for {
			select {
			case <-ticker.C:
				currentSecond := time.Now().Unix()
				var data []*protos.Datum

				for name, setting := range m.streamingConfig.Fields {
					frequency := setting.IntervalSeconds
					if currentSecond%int64(frequency/1000) == 0 {
						val := GetFieldValue(name)
						if val != nil {
							data = append(data, val)
						}
					}
				}
				err := m.transmitTelemetryData(data)
				if err != nil {
					m.logger.Errorf("error transmitting data %v", err)
				} else {
					m.logger.Infoln("Data transmitted")
				}

			case <-timer.C:
				ticker.Stop() // Stop the ticker when the timer expires
				return
			}
		}
	}()

	// Keep the main goroutine alive
	select {}

}

func (m *MockVehicle) transmitTelemetryData(data []*protos.Datum) error {
	msg, err := proto.Marshal(&protos.Payload{
		Data:      data,
		CreatedAt: timestamppb.Now(),
	})
	if err != nil {
		return err
	}
	m.counter++
	serverMessage := tesla.FlatbuffersStreamToBytes([]byte(fmt.Sprintf("vehicle_device.%s", m.vin)), []byte("V"), []byte(string(m.counter)), msg, 1, []byte(string(m.counter)), []byte("vehicle_device"), []byte(m.vin), uint64(time.Now().UnixMilli()))
	return m.SendMessage(serverMessage)
}

func (m *MockVehicle) SendMessage(message []byte) error {
	return m.connection.WriteMessage(websocket.BinaryMessage, message)
}

// GetTLSConfig returns a TLSConfig object from cert, key and optional client chain files.
func getTLSConfig(clienCertPath, clientKeyPath, ca string) (*tls.Config, error) {
	var cert tls.Certificate
	certFilePath, err := filepath.Abs(clienCertPath)
	if err != nil {
		return nil, err
	}
	keyFilePath, err := filepath.Abs(clientKeyPath)
	if err != nil {
		return nil, err
	}

	cert, err = tls.LoadX509KeyPair(certFilePath, keyFilePath)
	if err != nil {
		return nil, fmt.Errorf("can't properly load cert pair (%s, %s): %s", certFilePath, keyFilePath, err.Error())
	}

	clientCertPool := x509.NewCertPool()
	if ca != "" {
		_ = clientCertPool.AppendCertsFromPEM([]byte(ca))
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            clientCertPool,
		InsecureSkipVerify: true,
	}

	return tlsConfig, nil
}
