package nat

import (
	"component-master/config"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
)

type NatClient struct {
	client *nats.Conn
}

func NewNatsConnectClient(conf *config.NatServerConfig) (*NatClient, error) {
	url := conf.Host + ":" + fmt.Sprint(conf.Port)
	opts := []nats.Option{
		nats.Name("NATS Connection"),
		nats.ReconnectWait(10 * time.Second), // Wait time between reconnect attempts
		nats.MaxReconnects(-1),
		nats.DisconnectHandler(onDisconnect),
		nats.ReconnectHandler(onReconnect),
	}
	nc, err := nats.Connect(url, opts...)
	if err != nil {
		slog.Error(fmt.Sprintf("failed to connect to nats, %v", err))
		return nil, err
	}

	return &NatClient{
		client: nc,
	}, nil
}

func (nc *NatClient) Close() {
	if nc.client == nil {
		return
	}
	nc.client.Close()
}

func onDisconnect(_ *nats.Conn) {
	slog.Info(fmt.Sprintf("Disconnected from NATS: %d", time.Now().Unix()))
}

func onReconnect(_ *nats.Conn) {
	slog.Info(fmt.Sprintf("Reconnected from NATS: %d", time.Now().Unix()))
}
