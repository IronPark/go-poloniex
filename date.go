package poloniex

import (
	"errors"
	"strconv"
	"time"
	"encoding/json"
)

type PoloniexDate struct {
	time.Time
}

func (pd *PoloniexDate) UnmarshalJSON(data []byte) error {
	i, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return errors.New("Timestamp invalid (can't parse int64)")
	}
	pd.Time = time.Unix(i, 0)
	return nil
}


const TIME_FORMAT = "2006-01-02T15:04:05"

type jTime struct {
	time.Time
}

func (jt *jTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	t, err := time.Parse(TIME_FORMAT, s)
	if err != nil {
		return err
	}
	jt.Time = t
	return nil
}

func (jt jTime) MarshalJSON() ([]byte, error) {
	return json.Marshal((*time.Time)(&jt.Time).Format(TIME_FORMAT))
}