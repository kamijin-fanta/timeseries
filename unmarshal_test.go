package timeseries_test

import (
	"encoding/hex"
	"reflect"
	"testing"
	"time"

	"github.com/kamijin-fanta/timeseries"
)

func TestUnmarshal(t *testing.T) {
	testCases := []struct {
		input      string
		wantT0     uint64
		wantPoints []timeseries.Point
	}{
		{
			input:  "13ce4ca430cb400039bdf3b00100a0000000000003ffffffffe2329b000d60fffffffffffffffffc",
			wantT0: uint64(time.Date(2015, 3, 24, 2, 0, 0, 0, time.UTC).UnixNano()),
			wantPoints: []timeseries.Point{
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
		},
		{
			input:      "13ce4ca430cb4000fffffffffc0000000000000000",
			wantT0:     uint64(time.Date(2015, 3, 24, 2, 0, 0, 0, time.UTC).UnixNano()),
			wantPoints: nil,
		},
		{
			input:  "13ce4ca430cb400039bdf3b00100a0000000000003ffffffffffffffffc0",
			wantT0: uint64(time.Date(2015, 3, 24, 2, 0, 0, 0, time.UTC).UnixNano()),
			wantPoints: []timeseries.Point{
				{
					Timestamp: uint64(time.Date(2015, 3, 24, 2, 1, 2, 0, time.UTC).UnixNano()),
					Value:     12.0,
				},
			},
		},
		{
			input:  "13ce4ca430cb400039bdf3b00100a0000000000003ffffffffe2329b00378360020044cccccccccccfffffffffffffffffc0",
			wantT0: uint64(time.Date(2015, 3, 24, 2, 0, 0, 0, time.UTC).UnixNano()),
			wantPoints: []timeseries.Point{
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
		},
	}

	for _, tc := range testCases {
		input, err := hex.DecodeString(tc.input)
		if err != nil {
			t.Fatalf("failed to decode input hex string: tc.input=%s, err=%+v\n", tc.input, err)
		}

		t0, points, err := timeseries.Unmarshal(input)
		if err != nil {
			t.Fatalf("failed to unmarshal time series: tc.input=%s, err=%+v\n", tc.input, err)
		}
		if t0 != tc.wantT0 {
			t.Errorf("tc.input=%s, gotT0=%d, wantT0=%d", tc.input, t0, tc.wantT0)
		}
		if !reflect.DeepEqual(points, tc.wantPoints) {
			t.Errorf("tc.input=%s, gotPoints=%+v, wantPoints=%+v", tc.input, points, tc.wantPoints)
		}
	}
}
