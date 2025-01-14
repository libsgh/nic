package nic

import (
	"bytes"
	"encoding/json"
	"github.com/axgle/mahonia"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Response is the wrapper for http.Response
type Response struct {
	*http.Response
	encoding string
	Text     string
	Bytes    []byte
}

func NewResponse(r *http.Response) (*Response, error) {
	resp := &Response{
		Response: r,
		encoding: "utf-8",
		Text:     "",
		Bytes:    []byte{},
	}

	err := resp.bytes()
	if err != nil {
		return nil, err
	}
	resp.text()
	return resp, nil
}

func (r *Response) text() {
	r.Text = string(r.Bytes)
}

func (r *Response) bytes() error {
	/*defer r.Body.Close()
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}*/
	defer r.Body.Close()
	buffer := make([]byte, 8192)
	//result := bytes.NewBuffer(nil)
	for {
		_, err := r.Body.Read(buffer)

		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}
	// for multiple reading
	// e.g. goquery.NewDocumentFromReader
	r.Body = ioutil.NopCloser(bytes.NewBuffer(buffer))
	r.Bytes = buffer
	return nil
}

// JSON could parse http json response
func (r Response) JSON(s interface{}) error {
	// JSON response not must be `application/json` type
	// maybe `text/plain`...etc.
	// nic will parse it regardless of the content-type
	/*
		cType := r.Header.Get("Content-Type")
		if !strings.Contains(cType, "json") {
			return ErrNotJsonResponse
		}
	*/
	err := json.Unmarshal(r.Bytes, s)
	return err
}

// SetEncode changes Response.encoding
// and it changes Response.Text every times be invoked
func (r *Response) SetEncode(e string) error {
	if r.encoding != e {
		r.encoding = strings.ToLower(e)
		decoder := mahonia.NewDecoder(e)
		if decoder == nil {
			return ErrUnrecognizedEncoding
		}
		r.Text = decoder.ConvertString(r.Text)
	}
	return nil
}

// GetEncode returns Response.encoding
func (r Response) GetEncode() string {
	return r.encoding
}

// SaveFile save bytes data to a local file
func (r Response) SaveFile(filename string) error {
	dst, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = dst.Write(r.Bytes)
	if err != nil {
		return err
	}
	return nil
}
