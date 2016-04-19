package static

import (
	"bytes"
	"reflect"
	"testing"
)

func TestResponseWriterHeader(t *testing.T) {
	writtenBuffer := bytes.Buffer{}
	responseWriter := newResponseWriter(&writtenBuffer)

	// Header returns a map of header keys to values, which can be mutated.
	{
		header := responseWriter.Header()
		t.Logf("header => %#v", header)
		header["Content-Type"] = []string{"text/html"}
		t.Logf("header => %#v", header)
	}

	// Mutations to the Header returned are persisted within responseWriter
	{
		header := responseWriter.Header()
		t.Logf("header => %#v", header)
		content_type := header["Content-Type"]
		content_type_expected := []string{"text/html"}
		if !reflect.DeepEqual(content_type, content_type_expected) {
			t.Errorf("header[\"Content-Type\"] => %#v, want %#v", content_type, content_type_expected)
		}
	}

	// Changes to Header have no impact on the what's written
	{
		responseWriter.Write([]byte{0x01, 0x02, 0x03, 0x04})
		t.Logf("Write(%#v)", []byte{0x01, 0x02, 0x03, 0x04})
		written := writtenBuffer.Bytes()
		t.Logf("written => %#v", written)

		expected := []byte{0x01, 0x02, 0x03, 0x04}
		if !bytes.Equal(written, expected) {
			t.Errorf("written => %#v, want %#v", written, expected)
		}
	}
}

func TestResponseWriterWriteHeader(t *testing.T) {
	writtenBuffer := bytes.Buffer{}
	responseWriter := newResponseWriter(&writtenBuffer)

	// WriteHeader is safe and doesn't panic
	{
		responseWriter.WriteHeader(200)
		t.Logf("WriteHeader(200)")
	}

	// WriteHeader can be called multiple times with different values and doesn't panic
	{
		responseWriter.WriteHeader(200)
		t.Logf("WriteHeader(200)")
		responseWriter.WriteHeader(404)
		t.Logf("WriteHeader(404)")
		responseWriter.WriteHeader(500)
		t.Logf("WriteHeader(500)")
	}

	// WriteHeader has no impact on the what's written
	{
		responseWriter.Write([]byte{0x01, 0x02, 0x03, 0x04})
		t.Logf("Write(%#v)", []byte{0x01, 0x02, 0x03, 0x04})

		responseWriter.WriteHeader(200)
		t.Logf("WriteHeader(200)")
		written := writtenBuffer.Bytes()
		t.Logf("written => %#v", written)

		expected := []byte{0x01, 0x02, 0x03, 0x04}
		if !bytes.Equal(written, expected) {
			t.Errorf("written => %#v, want %#v", written, expected)
		}
	}
}

func TestResponseWriterWrite(t *testing.T) {
	writtenBuffer := bytes.Buffer{}
	responseWriter := newResponseWriter(&writtenBuffer)

	// No bytes written to the buffer, results in no bytes being written to
	{
		written := writtenBuffer.Bytes()
		t.Logf("written => %#v", written)

		expected := []byte{}
		if !bytes.Equal(written, expected) {
			t.Errorf("written => %#v, want %#v", written, expected)
		}
	}

	// Bytes written to the responseWriter, are written out to the writer
	{
		responseWriter.Write([]byte{0x01, 0x02, 0x03, 0x04})

		written := writtenBuffer.Bytes()
		t.Logf("written => %#v", written)

		expected := []byte{0x01, 0x02, 0x03, 0x04}
		if !bytes.Equal(written, expected) {
			t.Errorf("written => %#v, want %#v", written, expected)
		}
	}

	// Bytes written to the responseWriter in parts, are written out to the writer
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
