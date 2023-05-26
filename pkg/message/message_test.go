package message

import (
	"fmt"
	"strings"
	"testing"

	"gotest.tools/v3/assert"
)

func TestNewEvent(t *testing.T) {
	message := New()
	assert.Equal(t, message.APIVersion, "v1")
	assert.Equal(t, message.Specversion, "1.0")
	message.ID = "01F8TNE3GVBKDMZW9YWKYXX2N1"
	message.Source = "http://london.com/the_clash"
	message.Success = true
	message.Type = "com.event.receiver.group.triggered"
	message.Name = "foo"
	message.Version = "1.0.1"
	message.Release = "20210409.1617939933633"
	message.PlatformID = "x64-oci-linux-2"
	message.Package = "docker"
	message.Type = "com.event.example.type"
	messageJSON, err := message.ToJSON()
	assert.NilError(t, err, "error is not nil")
	fmt.Printf("%s\n", messageJSON)
}

func TestDecodeEventFromJSON(t *testing.T) {
	s := `{
		"id" : "01F8TNE3GVBKDMZW9YWKYXX2N1",
		"specversion" : "1.0",
		"source" : "http://london.com/the_clash",
		"success": true,
		"api_version": "v1",
		"type": "com.event.receiver.group.triggered",
		"name": "foo",
		"version": "1.0.1",
		"release": "20210409.1617939933633",
		"platform_id": "x64-oci-linux-2",
		"package": "docker",
		"type": "com.event.example.type",
		"data": {
			"event_groups": [
				{
				"id": "01H1ACAJ1AHNNECH0X3MADRH8B",
				"name": "clash",
				"type": "com.event.group.example",
				"version": "0001-01-01",
				"enabled": true
				}
			],
			"event_receivers": null,
			"events": null
		}
}
`
	message, err := DecodeFromJSON(strings.NewReader(s))
	assert.NilError(t, err, "error is not nil")
	assert.Equal(t, message.PlatformID, "x64-oci-linux-2")
}
