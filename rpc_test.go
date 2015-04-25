package rpc

import (
	"log"
	"net/url"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestURL(t *testing.T) {
	u := BuildHttpURL("127.0.0.1", "/.json", nil)
	assert.Equal(t, "http://127.0.0.1/.json", u.String())

	q := u.Query()
	q.Set("q", "golang")
	u.RawQuery = q.Encode()
	assert.Equal(t, "http://127.0.0.1/.json?q=golang", u.String())

	u2 := BuildHttpsURL("h2object.io:9000", "/.conf", q)
	assert.Equal(t, "https://h2object.io:9000/.conf?q=golang", u2.String())	
}

func TestRPC(t *testing.T) {
	var err error
	var ret map[string]interface{}


	u := BuildHttpURL("127.0.0.1:9000", "/.json", nil)

	v := url.Values{}
	v.Set("key", "james")
	u2 := BuildHttpURL("127.0.0.1:9000", "/.json", v)

	data := map[string]interface{}{
		"name":"james",
		"age":34,
	}

	DefaultClient := NewClient(H2OAnalyser{})

	err = DefaultClient.PostJson(nil, u2, data, &ret)
	assert.Nil(t, err)
	log.Println("Post result:", ret)
	ret = map[string]interface{}{}

	u3 := BuildHttpURL("127.0.0.1:9000", "/james.json", nil)

	mod := map[string]interface{}{
		"sex":true,
	}
	err = DefaultClient.PatchJson(nil, u3, mod, &ret)	
	assert.Nil(t, err)
	log.Println("Patch result:", ret)
	ret = map[string]interface{}{}
	
	err = DefaultClient.Get(nil, u3, &ret)
	assert.Nil(t, err)
	log.Println("Get result:", ret)
	ret = map[string]interface{}{}

	err = DefaultClient.PutJson(nil, u3, mod, &ret)	
	assert.Nil(t, err)
	log.Println("Put result:", ret)
	ret = map[string]interface{}{}

	err = DefaultClient.Get(nil, u3, &ret)
	assert.Nil(t, err)
	log.Println("Get result:", ret)
	ret = map[string]interface{}{}

	err = DefaultClient.Delete(nil, u3, &ret)	
	assert.Nil(t, err)
	log.Println("Delete result:", ret)
	ret = map[string]interface{}{}

	err = DefaultClient.Get(nil, u3, &ret)
	assert.Nil(t, err)

	err = DefaultClient.PostJson(nil, u, data, &ret)
	assert.Nil(t, err)
	log.Println("Post result:", ret)
	ret = map[string]interface{}{}

	err = DefaultClient.Get(nil, u, &ret)
	assert.Nil(t, err)
	log.Println("Get result:", ret)

	// u2 := BuildHttpURL("127.0.0.1:9000", "/james.json", nil)	
	// err = DefaultClient.Get(u2, ret)
	// assert.NotNil(t, err)
	// log.Println("err:", err)
}
