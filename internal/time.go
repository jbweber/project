package internal

import (
	"fmt"
	"strconv"
	"time"
)

type UnixTimestamp time.Time

func (u UnixTimestamp) MarshalJSON() ([]byte, error) {
	ts := time.Time(u).Unix()
	return []byte(fmt.Sprint(ts)), nil
}

func (u *UnixTimestamp) UnmarshalJSON(data []byte) error {
	v, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return err
	}

	*u = UnixTimestamp(time.Unix(v, 0))

	return nil
}
