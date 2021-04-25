package yu

import (
	"encoding/json"
	"time"
)

// timestamp format layout
const (
	// TimeLayout1 2006-01-02 15:04:05
	TimeLayout1 = "2006-01-02 15:04:05"
	// TimeLayout2 2006-01-02
	TimeLayout2 = "2006-01-02"
	// TimeLayout3 15:04:05
	TimeLayout3 = "15:04:05"
	// TimeLayout4 20060102150405
	TimeLayout4 = "20060102150405"
	// TimeLayout5 20060102
	TimeLayout5 = "20060102"
	// TimeLayout6 150405
	TimeLayout6 = "150405"
	// TimeLayout7 2006.01.02
	TimeLayout7 = "2006.01.02"
	// TimeLayout8 2006/01/02
	TimeLayout8 = "2006/01/02"
)

// TStamp custom type for int64
type TStamp int64

// NewTimeStamp new TStamp
func NewTimeStamp(ts int64) TStamp {
	return TStamp(ts)
}

// NewStrStamp new TStamp by time string
func NewStrStamp(ts string) TStamp {
	if ts == "" {
		return 0
	}
	d, err := time.ParseInLocation(TimeLayout1, ts, time.Local)
	if err != nil {
		return 0
	}
	return NewTimeStamp(d.Unix())
}

// NewStampTime new TStamp by time.Time
func NewStampTime(te time.Time) TStamp {
	return NewTimeStamp(te.Unix())
}

// String implement Stringer
func (t TStamp) String() string {
	if t <= 0 {
		return ""
	}
	s := t.Format(TimeLayout1)
	return s
}

// GoString implement GoStringer
func (t TStamp) GoString() string {
	return t.String()
}

// MarshalJSON implement Marshaler
func (t TStamp) MarshalJSON() (buf []byte, err error) {
	if t <= 0 {
		buf, err = json.Marshal("")
		return
	}
	buf, err = json.Marshal(t.Format(TimeLayout1))
	return
}

// UnmarshalJSON implement Unmarshaler
func (t *TStamp) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	d, err := time.ParseInLocation(TimeLayout1, s, time.Local)
	if err != nil {
		*t = 0
		return nil
	}
	*t = NewTimeStamp(d.Unix())
	return nil
}

// Int64 get int64 value
func (t TStamp) Int64() int64 {
	return int64(t)
}

// Time get time.Time
func (t TStamp) Time() time.Time {
	return time.Unix(int64(t), 0)
}

// Format get format time
func (t TStamp) Format(layout string) string {
	return t.Time().Format(layout)
}
