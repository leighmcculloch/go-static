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
	header_expected := http.Header{"Content-Type": []string{"text/html"}}
	header = responseWriter.Header()
	t.Logf("Header => %#v", header)
	if !reflect.DeepEqual(header, header_expected) {
		t.Errorf("Header => %#v, want %#v", header, header_expected)
	}
}

func TestResponseWriterStatus(t *testing.T) {
	writtenBuffer := bytes.Buffer{}
	responseWriter := newResponseWriter(&writtenBuffer)

	t.Log("Expect default Status to be 0.")
	expectedStatus := 0
	status := responseWriter.Status()
	t.Logf("Status => %d", status)
	if status != expectedStatus {
		t.Logf("Status => %d, want %d", status, expectedStatus)
	}
}

func TestResponseWriterStatusAfterWriteHeader(t *testing.T) {
	writtenBuffer := bytes.Buffer{}
	responseWriter := newResponseWriter(&writtenBuffer)

	t.Log("When WriteHeader has been called with a status code.")
	t.Log("Expect Status to be the status code.")
	{
		expectedStatus := 404
		responseWriter.WriteHeader(404)
		t.Logf("WriteHeader(404)")
		status := responseWriter.Status()
		t.Logf("Status => %d", status)
		if status != expectedStatus {
			t.Logf("Status => %d, want %d", status, expectedStatus)
		}
	}

	t.Log("When WriteHeader has been called multiple times.")
	t.Log("Expect Status to be the last status code.")
	{
		expectedStatus := 403
		responseWriter.WriteHeader(403)
		t.Logf("WriteHeader(403)")
		status := responseWriter.Status()
		t.Logf("Status => %d", status)
		if status != expectedStatus {
			t.Logf("Status => %d, want %d", status, expectedStatus)
		}
	}
}

func TestResponseWriterStatusAfterWrite(t *testing.T) {
	writtenBuffer := bytes.Buffer{}
	responseWriter := newResponseWriter(&writtenBuffer)

	t.Log("When Write has been called.")
	t.Log("Expect Status to be 200 OK.")
	expectedStatus := 200
	responseWriter.Write([]byte{})
	t.Logf("Write([]byte{})")
	status := responseWriter.Status()
	t.Logf("Status => %d", status)
	if status != expectedStatus {
		t.Logf("Status => %d, want %d", status, expectedStatus)
	}
}

func TestResponseWriterStatusAfterWriteHeaderAndWrite(t *testing.T) {
	writtenBuffer := bytes.Buffer{}
	responseWriter := newResponseWriter(&writtenBuffer)

	t.Log("When WriteHeader and Write have been called.")
	t.Log("Expect Status to be the WriteHeader input.")
	expectedStatus := 404
	responseWriter.WriteHeader(404)
	t.Logf("WriteHeader(404)")
	responseWriter.Write([]byte{})
	t.Logf("Write([]byte{})")
	status := responseWriter.Status()
	t.Logf("Status => %d", status)
	if status != expectedStatus {
		t.Logf("Status => %d, want %d", status, expectedStatus)
	}
}

func TestResponseWriterStatusAfterWriteAndWriteHeader(t *testing.T) {
	writtenBuffer := bytes.Buffer{}
	responseWriter := newResponseWriter(&writtenBuffer)

	t.Log("When Write and WriteHeader have been called.")
	t.Log("Expect Status to be the WriteHeader input.")
	expectedStatus := 404
	responseWriter.Write([]byte{})
	t.Logf("Write([]byte{})")
	responseWriter.WriteHeader(404)
	t.Logf("WriteHeader(404)")
	status := responseWriter.Status()
	t.Logf("Status => %d", status)
	if status != expectedStatus {
		t.Logf("Status => %d, want %d", status, expectedStatus)
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
