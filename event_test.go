package events

import (
	"encoding/base64"
	"fmt"
	"testing"
)

var testScte35 = "/DBrAACHCXBcAP/wBQb/hYm0rABVAlNDVUVJ/////3//AAAUmXANPw8cdXJuOm5iY3VuaS5jb206YnJjOjQxNDcyMzM5NAkfU0lHTkFMOkdxQS16X2paWlY0QUFBQUFBQUFLQVE9PTQDAXIWYww="

func TestNewEvent(t *testing.T) {

	base64Bytes, _ := base64.StdEncoding.DecodeString(testScte35)
	base64Bytes = append([]byte{0}, base64Bytes...)

	event, err := NewEvent(base64Bytes)
	if err != nil {
		t.Log("ERROR parsing SCTE-35 data: ", err)
		return
	}
	t.Log("Parse SCTE-35 successful. Updateing LastEventPTS")

	fmt.Println("Event ID: ", event.EventID)
	if event.EventID != 4294967295 {
		t.Errorf("Incorrect Event ID\n   Expected: 4294967295  Actual: %X\n", event.EventID)
	}

	fmt.Println("Event: ", event.Signal)

}
