package timeseries_test

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/kamijin-fanta/timeseries"
)

func TestMarshal(t *testing.T) {
	testCases := []struct {
		t0     uint64
		points []timeseries.Point
		want   string
	}{
		{
			t0: uint64(time.Date(2015, 3, 24, 2, 0, 0, 0, time.UTC).UnixNano()),
			points: []timeseries.Point{
				{
					Timestamp: uint64(time.Date(2015, 3, 24, 2, 1, 2, 0, time.UTC).UnixNano()),
					Value:     12.0,
				},
				{
					Timestamp: uint64(time.Date(2015, 3, 24, 2, 2, 2, 0, time.UTC).UnixNano()),
					Value:     12.0,
				},
				{
					Timestamp: uint64(time.Date(2015, 3, 24, 2, 3, 2, 0, time.UTC).UnixNano()),
					Value:     24.0,
				},
			},
			want: "13ce4ca430cb400039bdf3b00100a0000000000003ffffffffe2329b000d60fffffffffffffffffc",
		},
		{
			t0:     uint64(time.Date(2015, 3, 24, 2, 0, 0, 0, time.UTC).UnixNano()),
			points: []timeseries.Point{},
			want:   "13ce4ca430cb4000fffffffffc0000000000000000",
		},
		{
			t0: uint64(time.Date(2015, 3, 24, 2, 0, 0, 0, time.UTC).UnixNano()),
			points: []timeseries.Point{
				{
					Timestamp: uint64(time.Date(2015, 3, 24, 2, 1, 2, 0, time.UTC).UnixNano()),
					Value:     12.0,
				},
			},
			want: "13ce4ca430cb400039bdf3b00100a0000000000003ffffffffffffffffc0",
		},
		{
			t0: uint64(time.Date(2015, 3, 24, 2, 0, 0, 0, time.UTC).UnixNano()),
			points: []timeseries.Point{
				{
					Timestamp: uint64(time.Date(2015, 3, 24, 2, 1, 2, 0, time.UTC).UnixNano()),
					Value:     12.0,
				},
				{
					Timestamp: uint64(time.Date(2015, 3, 24, 2, 2, 2, 0, time.UTC).UnixNano()),
					Value:     12.5,
				},
				{
					Timestamp: uint64(time.Date(2015, 3, 24, 2, 3, 2, 0, time.UTC).UnixNano()),
					Value:     -24.2,
				},
			},
			want: "13ce4ca430cb400039bdf3b00100a0000000000003ffffffffe2329b00378360020044cccccccccccfffffffffffffffffc0",
		},
	}

	for _, tc := range testCases {
		buf, err := timeseries.Marshal(tc.t0, tc.points)
		if err != nil {
			t.Fatalf("failed to marshal points: t0=%d, points=%+v, err=%+v\n", tc.t0, tc.points, err)
		}
		got := hex.EncodeToString(buf)

		if got != tc.want {
			t.Errorf("t0=%d, points=%+v, got=%s, want=%s", tc.t0, tc.points, got, tc.want)
		}
	}
}
