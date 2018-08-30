package internal

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"
)

const unixRefTime int64 = 1136239445

func TestUnixTimeStamp(t *testing.T) {
	refTime := time.Unix(unixRefTime, 0)

	unixRef := UnixTimestamp(refTime)

	unixByte, err := json.Marshal(unixRef)
	if err != nil {
		t.Fatalf("Failed to Marshal UnixTimestamp: %v", err)
	}

	unixStr := string(unixByte)

	if unixStr != strconv.FormatInt(unixRefTime, 10) {
		t.Fatalf("%v does not equal %v", unixStr, strconv.FormatInt(unixRefTime, 10))
	}

}
