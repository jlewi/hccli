package pkg

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func Test_Deserialize(t *testing.T) {
	type testCase struct {
		name     string
		raw      string
		expected Query
	}

	startedAt, err := time.Parse(time.RFC3339, "2024-03-21T22:21:19.179830+00:00")
	if err != nil {
		t.Fatalf("Failed to parse time; %v", err)
	}

	completedAt, err := time.Parse(time.RFC3339, "2024-03-21T22:21:20.581871+00:00")
	if err != nil {
		t.Fatalf("Failed to parse time; %v", err)
	}
	cases := []testCase{
		{
			name: "simple",
			raw:  `{"input":{"nlq":"EMISSING slowest traces","cols":"['sli.latency', 'duration_ms', 'net.transport', 'http.method', 'error', 'http.target', 'http.route', 'rpc.method', 'ip', 'http.request_content_length', 'rpc.service', 'apdex', 'name', 'message.type', 'http.host', 'service.name', 'rpc.system', 'http.scheme', 'sli.platform-time', 'type', 'http.flavor', 'span.kind', 'dc.platform-time', 'library.version', 'status_code', 'net.host.port', 'net.host.ip', 'app.request_id', 'bucket_duration_ms', 'library.name', 'sli_product', 'message.uncompressed_size', 'rpc.grpc.status_code', 'net.peer.port', 'log10_duration_ms', 'http.status_code', 'status_message', 'http.user_agent', 'net.host.name', 'span.num_links', 'message.id', 'parent_name', 'app.cart_total', 'num_products', 'product_availability', 'revenue_at_risk', 'trace.trace_id', 'trace.span_id', 'ingest_timestamp', 'http.server_name', 'trace.parent_id']"},"output":"{'breakdowns': ['http.route'], 'calculations': [{'column': 'duration_ms', 'op': 'HEATMAP'}, {'column': 'duration_ms', 'op': 'MAX'}], 'filters': [{'column': 'trace.parent_id', 'op': 'does-not-exist'}, {'column': 'duration_ms', 'op': '>', 'value': 'threshold_value'}], 'orders': [{'column': 'duration_ms', 'op': 'MAX', 'order': 'descending'}], 'time_range': 7200}","id":null,"version":null,"created_at":null,"started_at":"2024-03-21T22:21:19.179830+00:00","completed_at":"2024-03-21T22:21:20.581871+00:00","logs":"","error":null,"status":"succeeded","metrics":{"predict_time":1.402041},"output_file_prefix":null,"webhook":null,"webhook_events_filter":["start","output","logs","completed"]}`,
			expected: Query{
				Input: &QueryInput{
					NLQ:  "EMISSING slowest traces",
					COLS: "['sli.latency', 'duration_ms', 'net.transport', 'http.method', 'error', 'http.target', 'http.route', 'rpc.method', 'ip', 'http.request_content_length', 'rpc.service', 'apdex', 'name', 'message.type', 'http.host', 'service.name', 'rpc.system', 'http.scheme', 'sli.platform-time', 'type', 'http.flavor', 'span.kind', 'dc.platform-time', 'library.version', 'status_code', 'net.host.port', 'net.host.ip', 'app.request_id', 'bucket_duration_ms', 'library.name', 'sli_product', 'message.uncompressed_size', 'rpc.grpc.status_code', 'net.peer.port', 'log10_duration_ms', 'http.status_code', 'status_message', 'http.user_agent', 'net.host.name', 'span.num_links', 'message.id', 'parent_name', 'app.cart_total', 'num_products', 'product_availability', 'revenue_at_risk', 'trace.trace_id', 'trace.span_id', 'ingest_timestamp', 'http.server_name', 'trace.parent_id']",
				},
				Output:              PtrToString("{'breakdowns': ['http.route'], 'calculations': [{'column': 'duration_ms', 'op': 'HEATMAP'}, {'column': 'duration_ms', 'op': 'MAX'}], 'filters': [{'column': 'trace.parent_id', 'op': 'does-not-exist'}, {'column': 'duration_ms', 'op': '>', 'value': 'threshold_value'}], 'orders': [{'column': 'duration_ms', 'op': 'MAX', 'order': 'descending'}], 'time_range': 7200}"),
				Status:              PtrToString("succeeded"),
				Metrics:             Metrics{PredictTime: 1.402041},
				StartedAt:           &startedAt,
				CompletedAt:         &completedAt,
				WebhookEventsFilter: []string{"start", "output", "logs", "completed"},
				Logs:                PtrToString(""),
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := &Query{}
			if err := json.Unmarshal([]byte(c.raw), actual); err != nil {
				t.Fatalf("Failed to deserialize; %v", err)
			}

			if d := cmp.Diff(c.expected, *actual); d != "" {
				t.Fatalf("Deserialized object is not equal to expected;diff:\n%v", d)
			}
		})
	}
}
