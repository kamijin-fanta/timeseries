package timeseries_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"time"

	"github.com/kamijin-fanta/timeseries"
)

func ExampleMarshal() {
	t0 := uint64(time.Date(2015, 3, 24, 2, 0, 0, 0, time.UTC).UnixNano())
	points := []timeseries.Point{
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
	}

	buf, err := timeseries.Marshal(t0, points)
	if err != nil {
		fmt.Printf("failed to marshal points: err=%+v\n", err)
		return
	}
	fmt.Println(hex.EncodeToString(buf))

	// Output: 13ce4ca430cb400039bdf3b00100a0000000000003ffffffffe2329b000d60fffffffffffffffffc
}

func ExampleUnmarshal() {
	input, err := hex.DecodeString("13ce4ca430cb400039bdf3b00100a0000000000003ffffffffe2329b000d60fffffffffffffffffc")
	if err != nil {
		fmt.Printf("failed to decode hex string: err=%+v\n", err)
		return
	}

	t0, points, err := timeseries.Unmarshal(input)
	if err != nil {
		fmt.Printf("failed to unmarshal time series: err=%+v\n", err)
		return
	}
	fmt.Printf("block timestamp=%v\n", time.Unix(int64(t0/1000000000), 0).UTC())
	for _, p := range points {
		fmt.Printf("timestamp=%v, value=%f\n", time.Unix(int64(p.Timestamp/1000000000), 0).UTC(), p.Value)
	}

	// Output:
	// block timestamp=2015-03-24 02:00:00 +0000 UTC
	// timestamp=2015-03-24 02:01:02 +0000 UTC, value=12.000000
	// timestamp=2015-03-24 02:02:02 +0000 UTC, value=12.000000
	// timestamp=2015-03-24 02:03:02 +0000 UTC, value=24.000000
}

func ExampleEncoder() {
	t0 := uint64(time.Date(2015, 3, 24, 2, 0, 0, 0, time.UTC).UnixNano())
	points := []timeseries.Point{
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
	}

	var b bytes.Buffer
	enc := timeseries.NewEncoder(&b)
	err := enc.EncodeHeader(t0)
	if err != nil {
		fmt.Printf("failed to encode time series header: err=%+v\n", err)
		return
	}

	for _, p := range points {
		err := enc.EncodePoint(p)
		if err != nil {
			fmt.Printf("failed to encode time series point: err=%+v\n", err)
			return
		}
	}

	err = enc.Finish()
	if err != nil {
		fmt.Printf("failed to encode time series finish marker: err=%+v\n", err)
		return
	}

	fmt.Println(hex.EncodeToString(b.Bytes()))

	// Output: 13ce4ca430cb400039bdf3b00100a0000000000003ffffffffe2329b000d60fffffffffffffffffc
}

func ExampleDecoder() {
	input, err := hex.DecodeString("13ce4ca430cb400039bdf3b00100a0000000000003ffffffffe2329b000d60fffffffffffffffffc")
	if err != nil {
		fmt.Printf("failed to decode hex string: err=%+v\n", err)
		return
	}

	b := bytes.NewBuffer(input)
	dec := timeseries.NewDecoder(b)

	t0, err := dec.DecodeHeader()
	if err != nil {
		fmt.Printf("failed to decode time series header: err=%+v\n", err)
		return
	}

	var points []timeseries.Point
	for {
		p, err := dec.DecodePoint()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("failed to decode time series point: err=%+v\n", err)
			return
		}
		points = append(points, p)
	}

	fmt.Printf("block timestamp=%v\n", time.Unix(int64(t0/1000000000), 0).UTC())
	for _, p := range points {
		fmt.Printf("timestamp=%v, value=%f\n", time.Unix(int64(p.Timestamp/1000000000), 0).UTC(), p.Value)
	}

	// Output:
	// block timestamp=2015-03-24 02:00:00 +0000 UTC
	// timestamp=2015-03-24 02:01:02 +0000 UTC, value=12.000000
	// timestamp=2015-03-24 02:02:02 +0000 UTC, value=12.000000
	// timestamp=2015-03-24 02:03:02 +0000 UTC, value=24.000000
}
