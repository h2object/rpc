package rpc

import (
	"log"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestMultiForm(t *testing.T) {
	multipart := NewMultipartForm()
	multipart.AddField("key1","valu1")
	multipart.AddField("key2","valu2")
	multipart.AddAttachment("file","test.png","./testdata/h2object-2p.png")

	ct, err := multipart.ContentType()
	assert.Nil(t, err)
	log.Println("multipart form content-type:", ct)

	sz, err := multipart.Size()
	assert.Nil(t, err)
	log.Println("multipart form sz:", sz)

	b, err := multipart.Bytes()
	assert.Nil(t, err)

	log.Println("multipart form content-type:", string(b))
}