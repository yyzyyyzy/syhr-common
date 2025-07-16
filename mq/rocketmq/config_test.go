package rocketmq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConf_Validate(t *testing.T) {
	c := &ProducerConf{
		NsResolver:                 []string{"127.0.0.1:9876"},
		GroupName:                  "",
		Namespace:                  "",
		InstanceName:               "",
		MsgTimeOut:                 0,
		DefaultTopicQueueNums:      0,
		CreateTopicKey:             "",
		CompressMsgBodyOverHowMuch: 0,
		CompressLevel:              0,
		Retry:                      0,
	}

	err := c.Validate()
	assert.Nil(t, err)
}
