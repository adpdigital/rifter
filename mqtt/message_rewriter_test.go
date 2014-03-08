package mqtt

import (
	"testing"

	"github.com/huin/mqtt"
	. "launchpad.net/gocheck"
)

func Test2(t *testing.T) { TestingT(t) }

type MessagRewriterSuite struct {
	msgRewriter *MsgRewriter
}

var _ = Suite(&MessagRewriterSuite{})

func (s *MessagRewriterSuite) SetUpTest(c *C) {
	s.msgRewriter = &MsgRewriter{
		IngressRewriter: &TopicPrefixRewriter{
			Prefix:      "$block",
			Replacement: "$cloud/123",
		},
		EgressRewriter: &TopicPrefixRewriter{
			Prefix:      "$cloud/123",
			Replacement: "$block",
		},
	}
}

func (s *MessagRewriterSuite) TestIngressMsgPublish(c *C) {

	// client publish a message to a topic
	pub := createPublish("$block/456/789")
	expectedPub := createPublish("$cloud/123/456/789")

	modPub := s.msgRewriter.RewriteIngress(pub)
	c.Assert(modPub, DeepEquals, expectedPub)

}

func (s *MessagRewriterSuite) TestIngressMsgSubscribe(c *C) {

	// client subscribe a message to a topic
	sub := createSubscribe("$block/456/789")
	expectedSub := createSubscribe("$cloud/123/456/789")

	modSub := s.msgRewriter.RewriteIngress(sub)
	c.Assert(modSub, DeepEquals, expectedSub)

}

func (s *MessagRewriterSuite) TestIngressMsgUnsubscribe(c *C) {
	// client unsubscribe to a topic
	unsub := createUnsubscribe("$block/456/789")
	expectedUnsub := createUnsubscribe("$cloud/123/456/789")

	modUnsub := s.msgRewriter.RewriteIngress(unsub)
	c.Assert(modUnsub, DeepEquals, expectedUnsub)
}

func (s *MessagRewriterSuite) TestEgressMsgPublish(c *C) {
	// client publish a message to a topic
	pub := createPublish("$cloud/123/456/789")
	expectedPub := createPublish("$block/456/789")

	modPub := s.msgRewriter.RewriteEgress(pub)
	c.Assert(modPub, DeepEquals, expectedPub)

}

func createPublish(topic string) mqtt.Message {
	return &mqtt.Publish{
		Header: mqtt.Header{
			DupFlag:  false,
			QosLevel: mqtt.QosAtMostOnce,
			Retain:   false,
		},
		TopicName: topic,
		Payload:   mqtt.BytesPayload{1, 2, 3},
	}
}

func createSubscribe(topic string) mqtt.Message {
	return &mqtt.Subscribe{
		Header: mqtt.Header{
			DupFlag:  false,
			QosLevel: mqtt.QosAtLeastOnce,
		},
		MessageId: 0x4321,
		Topics: []mqtt.TopicQos{
			{topic, mqtt.QosExactlyOnce},
		},
	}
}

func createUnsubscribe(topic string) mqtt.Message {
	return &mqtt.Unsubscribe{
		Header: mqtt.Header{
			DupFlag:  false,
			QosLevel: mqtt.QosAtLeastOnce,
		},
		MessageId: 0x4321,
		Topics:    []string{topic},
	}
}
