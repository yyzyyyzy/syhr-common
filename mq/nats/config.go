package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"time"
)

type Conf struct {
	// Host url example: nats://localhost:4222
	Hosts         []string `json:","`
	ReconnectWait int      `json:",optional,default=5"`
	MaxReconnect  int      `json:",optional,default=5"`
	UserCred      string   `json:",optional"`
	NkeyFile      string   `json:",optional"`
	TlsClientCert string   `json:",optional"`
	TlsClientKey  string   `json:",optional"`
	TlsCACert     string   `json:",optional"`
	UserJwt       string   `json:",optional"`
}

// NewConnect returns a nats connection
func (c Conf) NewConnect() (*nats.Conn, error) {
	option := []nats.Option{
		nats.ReconnectWait(time.Duration(c.ReconnectWait) * time.Second),
		nats.MaxReconnects(c.MaxReconnect),
		nats.RetryOnFailedConnect(true),
	}

	if c.UserCred != "" {
		option = append(option, nats.UserCredentials(c.UserCred))
	}

	if c.TlsClientCert != "" && c.TlsClientKey != "" {
		option = append(option, nats.ClientCert(c.TlsClientCert, c.TlsClientKey))
	}

	if c.NkeyFile != "" && c.UserJwt == "" {
		nKeyOption, err := nats.NkeyOptionFromSeed(c.NkeyFile)
		logx.Must(err)
		option = append(option, nKeyOption)
	}

	if c.NkeyFile != "" && c.UserJwt != "" {
		option = append(option, nats.UserCredentials(c.UserJwt, c.NkeyFile))
	}

	connect, err := nats.Connect(strings.Join(c.Hosts, ","), option...)

	if err != nil {
		logx.Error("failed to connect nat's server", logx.Field("detail", err), logx.Field("config",
			fmt.Sprintf("hosts: %s, userCred: %s", strings.Join(c.Hosts, ","), c.UserCred)))
		return nil, err
	}

	return connect, nil
}

// NewJetStream returns jet stream client instance
func (c Conf) NewJetStream() (jetstream.JetStream, error) {
	conn, err := c.NewConnect()
	if err != nil {
		return nil, err
	}
	js, err := jetstream.New(conn)
	if err != nil {
		logx.Error("failed to connect jet stream server", logx.Field("detail", err))
		return nil, err
	}
	return js, nil
}
