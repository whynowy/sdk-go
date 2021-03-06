package kafka_sarama_to_http_request_encode

import (
	"context"
	nethttp "net/http"
	"testing"

	"github.com/Shopify/sarama"

	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/cloudevents/sdk-go/pkg/binding"
	"github.com/cloudevents/sdk-go/pkg/binding/test"
	"github.com/cloudevents/sdk-go/pkg/bindings/http"
	"github.com/cloudevents/sdk-go/pkg/bindings/kafka_sarama"
)

var (
	e                                   = test.FullEvent()
	structuredConsumerMessageWithoutKey = &sarama.ConsumerMessage{
		Value: test.MustJSON(e),
		Headers: []*sarama.RecordHeader{{
			Key:   []byte(kafka_sarama.ContentType),
			Value: []byte(cloudevents.ApplicationCloudEventsJSON),
		}},
	}
	structuredConsumerMessageWithKey = &sarama.ConsumerMessage{
		Key:   []byte("aaa"),
		Value: test.MustJSON(e),
		Headers: []*sarama.RecordHeader{{
			Key:   []byte(kafka_sarama.ContentType),
			Value: []byte(cloudevents.ApplicationCloudEventsJSON),
		}},
	}
	binaryConsumerMessageWithoutKey = &sarama.ConsumerMessage{
		Value: []byte("hello world!"),
		Headers: mustToSaramaConsumerHeaders(map[string]string{
			"ce_type":            e.Type(),
			"ce_source":          e.Source(),
			"ce_id":              e.ID(),
			"ce_time":            test.Timestamp.String(),
			"ce_specversion":     "1.0",
			"ce_dataschema":      test.Schema.String(),
			"ce_datacontenttype": "text/json",
			"ce_subject":         "topic",
			"ce_exta":            "someext",
		}),
	}
	binaryConsumerMessageWithKey = &sarama.ConsumerMessage{
		Key:   []byte("akey"),
		Value: []byte("hello world!"),
		Headers: mustToSaramaConsumerHeaders(map[string]string{
			"ce_type":            e.Type(),
			"ce_source":          e.Source(),
			"ce_id":              e.ID(),
			"ce_time":            test.Timestamp.String(),
			"ce_specversion":     "1.0",
			"ce_dataschema":      test.Schema.String(),
			"ce_datacontenttype": "text/json",
			"ce_subject":         "topic",
			"ce_exta":            "someext",
		}),
	}
)

func mustToSaramaConsumerHeaders(m map[string]string) []*sarama.RecordHeader {
	res := make([]*sarama.RecordHeader, len(m))
	i := 0
	for k, v := range m {
		res[i] = &sarama.RecordHeader{Key: []byte(k), Value: []byte(v)}
		i++
	}
	return res
}

// Avoid DCE
var M binding.Message
var Req *nethttp.Request
var Err error

func BenchmarkStructuredWithKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		M, Err = kafka_sarama.NewMessage(structuredConsumerMessageWithKey)
		Req, Err = nethttp.NewRequest("POST", "http://localhost", nil)
		Err = http.EncodeHttpRequest(context.TODO(), M, Req, binding.TransformerFactories{})
	}
}

func BenchmarkStructuredWithoutKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		M, Err = kafka_sarama.NewMessage(structuredConsumerMessageWithoutKey)
		Req, Err = nethttp.NewRequest("POST", "http://localhost", nil)
		Err = http.EncodeHttpRequest(context.TODO(), M, Req, binding.TransformerFactories{})
	}
}

func BenchmarkBinaryWithKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		M, Err = kafka_sarama.NewMessage(binaryConsumerMessageWithKey)
		Req, Err = nethttp.NewRequest("POST", "http://localhost", nil)
		Err = http.EncodeHttpRequest(context.TODO(), M, Req, binding.TransformerFactories{})
	}
}

func BenchmarkBinaryWithoutKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		M, Err = kafka_sarama.NewMessage(binaryConsumerMessageWithoutKey)
		Req, Err = nethttp.NewRequest("POST", "http://localhost", nil)
		Err = http.EncodeHttpRequest(context.TODO(), M, Req, binding.TransformerFactories{})
	}
}
