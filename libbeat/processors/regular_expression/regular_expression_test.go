package regular_expression

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoToken(t *testing.T) {
	_, err := New("hello")
	assert.Equal(t, errInvalidTokenizer, err)
}

func TestDissectConversion(t *testing.T) {
	//tests := []struct {
	//	Name     string
	//	Tok      string
	//	Msg      string
	//	Expected map[string]interface{}
	//	Fail     bool
	//}{
	//	{
	//		Name: "Convert 1 value",
	//		Tok:  "id=%{id|integer} msg=\"%{message}\"",
	//		Msg:  "id=7736 msg=\"Single value OK\"}",
	//		Expected: map[string]interface{}{
	//			"id":      int32(7736),
	//			"message": "Single value OK",
	//		},
	//		Fail: false,
	//	},
	//	{
	//		Name: "Convert multiple values values",
	//		Tok:  "id=%{id|integer} status=%{status|integer} duration=%{duration|float} uptime=%{uptime|long} success=%{success|boolean} msg=\"%{message}\"",
	//		Msg:  "id=7736 status=202 duration=0.975 uptime=1588975628 success=true msg=\"Request accepted\"}",
	//		Expected: map[string]interface{}{
	//			"id":       int32(7736),
	//			"status":   int32(202),
	//			"duration": float32(0.975),
	//			"uptime":   int64(1588975628),
	//			"success":  true,
	//			"message":  "Request accepted",
	//		},
	//		Fail: false,
	//	},
	//	{
	//		Name: "Convert 1 indirect field value",
	//		Tok:  "%{?k1}=%{&k1|integer} msg=\"%{message}\"",
	//		Msg:  "id=8268 msg=\"Single value indirect field\"}",
	//		Expected: map[string]interface{}{
	//			"id":      int32(8268),
	//			"message": "Single value indirect field",
	//		},
	//		Fail: false,
	//	},
	//	{
	//		Name: "Greedy padding skip test ->",
	//		Tok:  "id=%{id->|integer} padding_removed=%{padding_removed->|boolean} length=%{length->|long} msg=\"%{message}\"",
	//		Msg:  "id=1945     padding_removed=true    length=123456789    msg=\"Testing for padding\"}",
	//		Expected: map[string]interface{}{
	//			"id":              int32(1945),
	//			"padding_removed": true,
	//			"length":          int64(123456789),
	//			"message":         "Testing for padding",
	//		},
	//		Fail: false,
	//	},
	//}

}

func TestEmptyString(t *testing.T) {
	//d, err := New("(.*?)(?<=a)b")  //不合规的正则表达式
	d, err := New("^(\\d{4}-\\d{2}-\\d{2} \\d{2}:\\d{2}:\\d{2}\\.\\d{3})\\s(\\w+)\\s\\[(.*?)\\]\\s(?:(.*?))?\\s+(\\{.*\\})")
	_, err = d.Regular("2024-10-15 16:21:23.874 DEBUG [http-nio-8601-exec-2] c.c.u.c.c.OGremlinConnection {\"Message\":\"sql(g.V().hasLabel(\"operatingSystem\").has(\"id\", P.within(assetIdList)).project(\"id\",\"systemType\",\"hostName\",\"isVirtual\",\"ipList\",\"empeeAcct\").by(\"id\").by(\"systemType\").by(\"hostName\").by(\"isVirtual\").by(.out().hasLabel(\"ip\").values(\"ip\").fold().map{it.get().join(\",\")}).by(.out().hasLabel(\"account\").values(\"name\").fold().map{it.get().join(\",\")})) parameters(assetIdList=[9e619f97-fb9b-47af-9433-0e9e0618ff77, 9689600b-8d6a-40c9-a823-ab315ea684bb, d96e2cb4-d52c-4abc-a26e-265acb4b05bd],)\",\"User\":\"root\",\"IP\":\"134.64.110.143,134.95.172.12:36296\",\"EmpeeAcct\":\"admin\",\"Name\":\"cmdb\"}")
	assert.Equal(t, errEmpty, err)
}

// JSON tags are used to create a common test file for the `logstash-filter-dissect` and the
// beat implementation.
type dissectTest struct {
	Name     string `json:"name"`
	Tok      string `json:"tok"`
	Msg      string `json:"msg"`
	Expected Map    `json:"expected"`
	Skip     bool   `json:"skip"`
	Fail     bool   `json:"fail"`
}

var tests []dissectTest

func init() {
	content, err := ioutil.ReadFile("testdata/dissect_tests.json")
	if err != nil {
		fmt.Printf("could not read the content of 'dissect_tests', error: %s", err)
		os.Exit(1)
	}

	if err := json.Unmarshal(content, &tests); err != nil {
		fmt.Printf("could not parse the content of 'dissect_tests', error: %s", err)
		os.Exit(1)
	}
}

func TestDissect(t *testing.T) {
	if len(tests) == 0 {
		t.Error("No test cases were loaded")
	}

	for _, test := range tests {
		if test.Skip {
			continue
		}
		t.Run(test.Name, func(t *testing.T) {
			d, err := New(test.Tok)
			if !assert.NoError(t, err) {
				return
			}

			if test.Fail {
				_, err := d.Regular(test.Msg)
				assert.Error(t, err)
				return
			}

			r, err := d.Regular(test.Msg)
			if !assert.NoError(t, err) {
				return
			}

			assert.Equal(t, test.Expected, r)
		})
	}
}

var results []string
var o [][]string

func BenchmarkDissect(b *testing.B) {
	for _, test := range tests {
		if test.Skip {
			continue
		}
		b.Run(test.Name, func(b *testing.B) {
			tok := test.Tok
			msg := test.Msg
			d, err := New(tok)
			if !assert.NoError(b, err) {
				return
			}
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				r, err := d.Regular(msg)
				if test.Fail {
					assert.Error(b, err)
					return
				}
				assert.NoError(b, err)
				results = r
			}
		})
	}

	// Add a few regular expression matches against the same string the test suite,
	// this give us a baseline to compare to, note that we only test a raw match against the string.
	b.Run("Regular expression", func(b *testing.B) {
		re := regexp.MustCompile("/var/log/([a-z]+)/log/([a-z]+)/apache/([a-b]+)")
		by := "/var/log/docker/more/apache/super"
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			o = re.FindAllStringSubmatch(by, -1)
		}
	})

	b.Run("Larger regular expression", func(b *testing.B) {
		re := regexp.MustCompile("^(\\d{2})-(\\w{3})-(\\d{4})\\s([0-9:.]+)\\s(\\w+)\\s\\[([a-zA-Z0-9-]+)\\]\\s([a-zA-Z0-9.]+)\\s(.+)")

		by := `18-Apr-2018 06:53:20.411 INFO [http-nio-8080-exec-1] org.apache.coyote.http11.Http11Processor.service Error parsing HTTP request header
 Note: further occurrences of HTTP header parsing errors will be logged at DEBUG level.
 java.lang.IllegalArgumentException: Invalid character found in method name. HTTP method names must be tokens
    at org.apache.coyote.http11.Http11InputBuffer.parseRequestLine(Http11InputBuffer.java:426)
    at org.apache.coyote.http11.Http11Processor.service(Http11Processor.java:687)
    at org.apache.coyote.AbstractProcessorLight.process(AbstractProcessorLight.java:66)
    at org.apache.coyote.AbstractProtocol$ConnectionHandler.process(AbstractProtocol.java:790)
    at org.apache.tomcat.util.net.NioEndpoint$SocketProcessor.doRun(NioEndpoint.java:1459)
    at org.apache.tomcat.util.net.SocketProcessorBase.run(SocketProcessorBase.java:49)
    at java.util.concurrent.ThreadPoolExecutor.runWorker(ThreadPoolExecutor.java:1149)
    at java.util.concurrent.ThreadPoolExecutor$Worker.run(ThreadPoolExecutor.java:624)
    at org.apache.tomcat.util.threads.TaskThread$WrappingRunnable.run(TaskThread.java:61)
    at java.lang.Thread.run(Thread.java:748)`
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			o = re.FindAllStringSubmatch(by, -1)
		}
	})

	b.Run("regular expression to match end of line", func(b *testing.B) {
		re := regexp.MustCompile("MACHINE\\[(\\w+)\\]$")

		by := `18-Apr-2018 06:53:20.411 INFO [http-nio-8080-exec-1] org.apache.coyote.http11.Http11Processor.service Error parsing HTTP request header
 Note: further occurrences of HTTP header parsing errors will be logged at DEBUG level.
 java.lang.IllegalArgumentException: Invalid character found in method name. HTTP method names must be tokens
    at org.apache.coyote.http11.Http11InputBuffer.parseRequestLine(Http11InputBuffer.java:426)
    at org.apache.coyote.http11.Http11Processor.service(Http11Processor.java:687)
    at org.apache.coyote.AbstractProcessorLight.process(AbstractProcessorLight.java:66)
    at org.apache.coyote.AbstractProtocol$ConnectionHandler.process(AbstractProtocol.java:790)
    at org.apache.tomcat.util.net.NioEndpoint$SocketProcessor.doRun(NioEndpoint.java:1459)
    at org.apache.tomcat.util.net.SocketProcessorBase.run(SocketProcessorBase.java:49)
    at java.util.concurrent.ThreadPoolExecutor.runWorker(ThreadPoolExecutor.java:1149)
    at java.util.concurrent.ThreadPoolExecutor$Worker.run(ThreadPoolExecutor.java:624)
    at org.apache.tomcat.util.threads.TaskThread$WrappingRunnable.run(TaskThread.java:61)
    at java.lang.Thread.run(Thread.java:748) MACHINE[hello]`
		b.ReportAllocs()
		for n := 0; n < b.N; n++ {
			o = re.FindAllStringSubmatch(by, -1)
		}
	})
}

func dissectConversion(tok, msg string, b *testing.B) {
	d, err := New(tok)
	assert.NoError(b, err)

	_, err = d.Regular(msg)
	assert.NoError(b, err)
}

func benchmarkConversion(tok, msg string, b *testing.B) {
	for n := 0; n < b.N; n++ {
		dissectConversion(tok, msg, b)
	}
}

func BenchmarkDissectNoConversionOneValue(b *testing.B) {
	b.ReportAllocs()
	benchmarkConversion("id=%{id} msg=\"%{message}\"", "id=7736 msg=\"Single value OK\"}", b)
}

func BenchmarkDissectWithConversionOneValue(b *testing.B) {
	b.ReportAllocs()
	benchmarkConversion("id=%{id|integer} msg=\"%{message}\"", "id=7736 msg=\"Single value OK\"}", b)
}

func BenchmarkDissectNoConversionMultipleValues(b *testing.B) {
	b.ReportAllocs()
	benchmarkConversion("id=%{id} status=%{status} duration=%{duration} uptime=%{uptime} success=%{success} msg=\"%{message}\"",
		"id=7736 status=202 duration=0.975 uptime=1588975628 success=true msg=\"Request accepted\"}", b)
}

func BenchmarkDissectWithConversionMultipleValues(b *testing.B) {
	b.ReportAllocs()
	benchmarkConversion("id=%{id|integer} status=%{status|integer} duration=%{duration|float} uptime=%{uptime|long} success=%{success|boolean} msg=\"%{message}\"",
		"id=7736 status=202 duration=0.975 uptime=1588975628 success=true msg=\"Request accepted\"}", b)
}

func BenchmarkDissectComplexStackTraceDegradation(b *testing.B) {
	message := `18-Apr-2018 06:53:20.411 INFO [http-nio-8080-exec-1] org.apache.coyote.http11.Http11Processor.service Error parsing HTTP request header
 Note: further occurrences of HTTP header parsing errors will be logged at DEBUG level.
 java.lang.IllegalArgumentException: Invalid character found in method name. HTTP method names must be tokens
    at org.apache.coyote.http11.Http11InputBuffer.parseRequestLine(Http11InputBuffer.java:426)
    at org.apache.coyote.http11.Http11Processor.service(Http11Processor.java:687)
    at org.apache.coyote.AbstractProcessorLight.process(AbstractProcessorLight.java:66)
    at org.apache.coyote.AbstractProtocol$ConnectionHandler.process(AbstractProtocol.java:790)
    at org.apache.tomcat.util.net.NioEndpoint$SocketProcessor.doRun(NioEndpoint.java:1459)
    at org.apache.tomcat.util.net.SocketProcessorBase.run(SocketProcessorBase.java:49)
    at java.util.concurrent.ThreadPoolExecutor.runWorker(ThreadPoolExecutor.java:1149)
    at java.util.concurrent.ThreadPoolExecutor$Worker.run(ThreadPoolExecutor.java:624)
    at org.apache.tomcat.util.threads.TaskThread$WrappingRunnable.run(TaskThread.java:61)
    at java.lang.Thread.run(Thread.java:748) MACHINE[hello]`

	tests := []struct {
		Name string
		Tok  string
	}{
		{
			Name: "ComplexStackTrace-1",
			Tok:  "%{origin} %{message}",
		},
		{
			Name: "ComplexStackTrace-2",
			Tok:  "%{day} %{origin} %{message}",
		},
		{
			Name: "ComplexStackTrace-3",
			Tok:  "%{day}-%{month} %{origin} %{message}",
		},
		{
			Name: "ComplexStackTrace-4",
			Tok:  "%{day}-%{month}-%{year} %{origin} %{message}",
		},
		{
			Name: "ComplexStackTrace-5",
			Tok:  "%{day}-%{month}-%{year} %{hour} %{origin} %{message}",
		},
		{
			Name: "ComplexStackTrace-6",
			Tok:  "%{day}-%{month}-%{year} %{hour} %{severity} %{origin} %{message}",
		},
		{
			Name: "ComplexStackTrace-7",
			Tok:  "%{day}-%{month}-%{year} %{hour} %{severity} [%{thread_id}] %{origin} %{message}",
		},
		{
			Name: "ComplexStackTrace-8",
			Tok:  "%{day}-%{month}-%{year} %{hour} %{severity} [%{thread_id}] %{origin} %{first_line} %{message}",
		},
	}

	for _, test := range tests {
		b.Run(test.Name, func(b *testing.B) {
			tok := test.Tok
			msg := message
			d, err := New(tok)
			if !assert.NoError(b, err) {
				return
			}
			b.ReportAllocs()
			for n := 0; n < b.N; n++ {
				r, err := d.Regular(msg)
				assert.NoError(b, err)
				results = r
			}
		})
	}
}
