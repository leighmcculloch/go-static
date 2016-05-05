package static

import (
	"bytes"
	"net/http"
	"reflect"
	"testing"
)

func TestResponseWriterHeader(t *testing.T) {
	writtenBuffer := bytes.Buffer{}
	responseWriter := newResponseWriter(&writtenBuffer)

	t.Log("Setting the Content-Type to text/html.")
	header := responseWriter.Header()
	header["Content-Type"] = []string{"text/html"}
	t.Logf("Header => %#v", header)

	t.Log("Expect Content-Type to be persisted.")
	headerExpected := http.Header{"Content-Type": []string{"text/html"}}
	header = responseWriter.Header()
	t.Logf("Header => %#v", header)
	if !reflect.DeepEqual(header, headerExpected) {
		t.Errorf("Header => %#v, want %#v", header, headerExpected)
	}
}

func TestResponseWriterStatusCode(t *testing.T) {
	writtenBuffer := bytes.Buffer{}
	responseWriter := newResponseWriter(&writtenBuffer)

	t.Log("Expect default StatusCode to be 0.")
	expectedStatusCode := 0
	statusCode := responseWriter.StatusCode()
	t.Logf("StatusCode => %d", statusCode)
	if statusCode != expectedStatusCode {
		t.Logf("StatusCode => %d, want %d", statusCode, expectedStatusCode)
	}
}

func TestResponseWriterStatusCodeAfterWriteHeader(t *testing.T) {
	writtenBuffer := bytes.Buffer{}
	responseWriter := newResponseWriter(&writtenBuffer)

	t.Log("When WriteHeader has been called with a status code.")
	t.Log("Expect StatusCode to be the status code.")
	{
		expectedStatusCode := 404
		responseWriter.WriteHeader(404)
		t.Logf("WriteHeader(404)")
		statusCode := responseWriter.StatusCode()
		t.Logf("StatusCode => %d", statusCode)
		if statusCode != expectedStatusCode {
			t.Logf("StatusCode => %d, want %d", statusCode, expectedStatusCode)
		}
	}

	t.Log("When WriteHeader has been called multiple times.")
	t.Log("Expect StatusCode to be the last status code.")
	{
		expectedStatusCode := 403
		responseWriter.WriteHeader(403)
		t.Logf("WriteHeader(403)")
		statusCode := responseWriter.StatusCode()
		t.Logf("StatusCode => %d", statusCode)
		if statusCode != expectedStatusCode {
			t.Logf("StatusCode => %d, want %d", statusCode, expectedStatusCode)
		}
	}
}

func TestResponseWriterStatusCodeAfterWrite(t *testing.T) {
	writtenBuffer := bytes.Buffer{}
	responseWriter := newResponseWriter(&writtenBuffer)

	t.Log("When Write has been called.")
	t.Log("Expect StatusCode to be 200 OK.")
	expectedStatusCode := 200
	responseWriter.Write([]byte{})
	t.Logf("Write([]byte{})")
	statusCode := responseWriter.StatusCode()
	t.Logf("StatusCode => %d", statusCode)
	if statusCode != expectedStatusCode {
		t.Logf("StatusCode => %d, want %d", statusCode, expectedStatusCode)
	}
}

func TestResponseWriterStatusCodeAfterWriteHeaderAndWrite(t *testing.T) {
	writtenBuffer := bytes.Buffer{}
	responseWriter := newResponseWriter(&writtenBuffer)

	t.Log("When WriteHeader and Write have been called.")
	t.Log("Expect StatusCode to be the WriteHeader input.")
	expectedStatusCode := 404
	responseWriter.WriteHeader(404)
	t.Logf("WriteHeader(404)")
	responseWriter.Write([]byte{})
	t.Logf("Write([]byte{})")
	statusCode := responseWriter.StatusCode()
	t.Logf("StatusCode => %d", statusCode)
	if statusCode != expectedStatusCode {
		t.Logf("StatusCode => %d, want %d", statusCode, expectedStatusCode)
	}
}

func TestResponseWriterStatusCodeAfterWriteAndWriteHeader(t *testing.T) {
	writtenBuffer := bytes.Buffer{}
	responseWriter := newResponseWriter(&writtenBuffer)

	t.Log("When Write and WriteHeader have been called.")
	t.Log("Expect StatusCode to be the WriteHeader input.")
	expectedStatusCode := 404
	responseWriter.Write([]byte{})
	t.Logf("Write([]byte{})")
	responseWriter.WriteHeader(404)
	t.Logf("WriteHeader(404)")
	statusCode := responseWriter.StatusCode()
	t.Logf("StatusCode => %d", statusCode)
	if statusCode != expectedStatusCode {
		t.Logf("StatusCode => %d, want %d", statusCode, expectedStatusCode)
	}
}

func TestResponseWriterWrite(t *testing.T) {
	writtenBuffer := bytes.Buffer{}
	responseWriter := newResponseWriter(&writtenBuffer)

	t.Log("When nothing has been written to the responseWriter.")
	t.Log("Expect nothing will be written to the Writer.")
	{
		expected := []byte{}
		written := writtenBuffer.Bytes()
		t.Logf("written => %#v", written)
		if !bytes.Equal(written, expected) {
			t.Errorf("written => %#v, want %#v", written, expected)
		}
	}

	t.Log("When bytes are written to the responseWriter.")
	t.Log("The bytes will be written to the Writer.")
	{
		responseWriter.Write([]byte{0x01, 0x02, 0x03, 0x04})

		written := writtenBuffer.Bytes()
		t.Logf("written => %#v", written)

		expected := []byte{0x01, 0x02, 0x03, 0x04}
		if !bytes.Equal(written, expected) {
			t.Errorf("written => %#v, want %#v", written, expected)
		}
	}

	t.Log("When additional bytes are written to the responseWriter.")
	t.Log("The additional bytes will be written to the Writer.")
	{
		responseWriter.Write([]byte{0x05, 0x06})
		t.Logf("Write(%#v)", []byte{0x05, 0x06})
		responseWriter.Write([]byte{0x07, 0x08})
		t.Logf("Write(%#v)", []byte{0x07, 0x08})

		written := writtenBuffer.Bytes()
		t.Logf("written => %#v", written)

		expected := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
		if !bytes.Equal(written, expected) {
			t.Errorf("written => %#v, want %#v", written, expected)
		}
	}
}
