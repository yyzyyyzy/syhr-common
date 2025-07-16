package mqtt

import (
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/zeromicro/go-zero/core/logx"
)

type Config struct {
	Broker   string
	ClientID string
	Username string
	Password string
	Qos      byte
}

type Client struct {
	client mqtt.Client
	config Config
}

func (cfg Config) NewMqttClient() *Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(cfg.Broker)
	opts.SetClientID(cfg.ClientID)
	opts.SetUsername(cfg.Username)
	opts.SetPassword(cfg.Password)
	opts.SetAutoReconnect(true)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		logx.Errorf("MQTT connection failed: %v", token.Error())
		panic(token.Error())
	}
	logx.Info("Connected to EMQX")
	return &Client{client: client, config: cfg}
}

func (c *Client) Publish(topic string, payload []byte) error {
	token := c.client.Publish(topic, c.config.Qos, false, payload)
	token.Wait()
	return token.Error()
}

func (c *Client) Subscribe(topic string, handler mqtt.MessageHandler) error {
	token := c.client.Subscribe(topic, c.config.Qos, handler)
	token.Wait()
	return token.Error()
}

func (c *Client) Disconnect() {
	c.client.Disconnect(250)
}
