package static

import (
	"bytes"
	"reflect"
	"testing"
)

func TestResponseBufferHeader(t *testing.T) {
	responseBuffer := newResponseBuffer()

	// Header returns a map of header keys to values, which can be mutated.
	{
		header := responseBuffer.Header()
		t.Logf("header => %#v", header)
		header["Content-Type"] = []string{"text/html"}
		t.Logf("header => %#v", header)
	}

	// Mutations to the Header returned are persisted within responseBuffer
	{
		header := responseBuffer.Header()
		t.Logf("header => %#v", header)
		content_type := header["Content-Type"]
		content_type_expected := []string{"text/html"}
		if !reflect.DeepEqual(content_type, content_type_expected) {
			t.Errorf("header[\"Content-Type\"] => %#v, want %#v", content_type, content_type_expected)
		}
	}
}

func TestResponseBufferWrite(t *testing.T) {
	responseBuffer := newResponseBuffer()

	// Write adds to the internal buffer
	{
		before := responseBuffer.Bytes()
		t.Logf("buffer => %#v", before)

		responseBuffer.Write([]byte{0x01, 0x02, 0x03, 0x04})
		t.Logf("Write(%#v)", []byte{0x01, 0x02, 0x03, 0x04})

		after := responseBuffer.Bytes()
		t.Logf("buffer => %#v", after)

		expected := []byte{0x01, 0x02, 0x03, 0x04}
		if !bytes.Equal(after, expected) {
			t.Errorf("buffer => %#v, want %#v", after, expected)
		}
	}

	// Multiple writes build on the internal buffer
	{
		before := responseBuffer.Bytes()
		t.Logf("buffer => %#v", before)

		responseBuffer.Write([]byte{0x05, 0x06})
		t.Logf("Write(%#v)", []byte{0x05, 0x06})
		responseBuffer.Write([]byte{0x07, 0x08})
		t.Logf("Write(%#v)", []byte{0x07, 0x08})

		after := responseBuffer.Bytes()
		t.Logf("buffer => %#v", after)

		expected := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08}
		if !bytes.Equal(after, expected) {
			t.Errorf("buffer => %#v, want %#v", after, expected)
		}
	}
}

func TestResponseBufferWriteHeader(t *testing.T) {
	responseBuffer := newResponseBuffer()

	// WriteHeader does nothing
	{
		responseBuffer.WriteHeader(200)
		t.Logf("WriteHeader(200)")
	}

	// WriteHeader can be called multiple times with different values
	{
		responseBuffer.WriteHeader(200)
		t.Logf("WriteHeader(200)")
		responseBuffer.WriteHeader(404)
		t.Logf("WriteHeader(404)")
		responseBuffer.WriteHeader(500)
		t.Logf("WriteHeader(500)")
	}

	// WriteHeader has no impact on the buffer
	{
		responseBuffer.Write([]byte("Hello World!"))
		before := responseBuffer.Bytes()
		t.Logf("buffer => %#v", before)

		responseBuffer.WriteHeader(200)
		t.Logf("WriteHeader(200)")
		after := responseBuffer.Bytes()
		t.Logf("buffer => %#v", after)

		expected := []byte("Hello World!")
		if !bytes.Equal(after, expected) {
			t.Errorf("buffer => %#v, want %#v", after, expected)
		}
	}
}

func TestResponseBufferWriteTo(t *testing.T) {
	responseBuffer := newResponseBuffer()

	// No bytes written to the buffer, results in no bytes being written to
	{
		buffer := responseBuffer.Bytes()
		t.Logf("buffer => %#v", buffer)

		writer := bytes.Buffer{}
		responseBuffer.WriteTo(&writer)
		t.Logf("WriteTo(writer)")
		written := writer.Bytes()
		t.Logf("written => %#v", written)

		expected := []byte{}
		if !bytes.Equal(written, expected) {
			t.Errorf("written => %#v, want %#v", written, expected)
		}
	}

	// Bytes written to the buffer, can be written out to another
	{
		responseBuffer.Write([]byte{0x01, 0x02, 0x03, 0x04})
		buffer := responseBuffer.Bytes()
		t.Logf("buffer => %#v", buffer)

		writer := bytes.Buffer{}
		responseBuffer.WriteTo(&writer)
		t.Logf("WriteTo(writer)")
		written := writer.Bytes()
		t.Logf("written => %#v", written)

		expected := []byte{0x01, 0x02, 0x03, 0x04}
		if !bytes.Equal(written, expected) {
			t.Errorf("written => %#v, want %#v", written, expected)
		}
	}
}
