// This is the implementation of the following dataflow
//
// sensor1 --\
//           |
//          merger --> duplicator --> aggregator --> display
//           |             |
// sensor2 --/             \--> log

package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	. "github.com/pspaces/gospace"
)

func main() {

	streams := NewSpace("tcp://localhost:31145/streams")

	// Create the dataflow network
	var streamID int
	var streamURI string
	var stream Space

	streamID = 0
	streamURI = "tcp://localhost:" + strconv.Itoa(31146+streamID) + "/stream" + strconv.Itoa(streamID)
	stream = NewSpace(streamURI)
	streams.Put("sensor1", "merger", streamURI)

	streamID++
	streamURI = "tcp://localhost:" + strconv.Itoa(31146+streamID) + "/stream" + strconv.Itoa(streamID)
	stream = NewSpace(streamURI)
	streams.Put("sensor2", "merger", streamURI)

	streamID++
	streamURI = "tcp://localhost:" + strconv.Itoa(31146+streamID) + "/stream" + strconv.Itoa(streamID)
	stream = NewSpace(streamURI)
	streams.Put("merger", "duplicator", streamURI)

	streamID++
	streamURI = "tcp://localhost:" + strconv.Itoa(31146+streamID) + "/stream" + strconv.Itoa(streamID)
	stream = NewSpace(streamURI)
	streams.Put("duplicator", "display1", streamURI)

	streamID++
	streamURI = "tcp://localhost:" + strconv.Itoa(31146+streamID) + "/stream" + strconv.Itoa(streamID)
	stream = NewSpace(streamURI)
	streams.Put("duplicator", "aggregator", streamURI)

	streamID++
	streamURI = "tcp://localhost:" + strconv.Itoa(31146+streamID) + "/stream" + strconv.Itoa(streamID)
	stream = NewSpace(streamURI)
	streams.Put("aggregator", "display2", streamURI)

	// Launch all stream processing and routing units
	go sensor(&streams, "sensor1")
	go sensor(&streams, "sensor2")
	go merger(&streams, "merger")
	go duplicator(&streams, "duplicator")
	go aggregator(&streams, "aggregator")
	go display(&streams, "display2")
	go display(&streams, "display1")

	stream.Get("stop")

}

func sensor(streams *Space, me string) {
	var target string
	var targetURI string
	t, _ := streams.Query(me, &target, &targetURI)
	targetStream := NewRemoteSpace((t.GetFieldAt(2)).(string))
	for timestamp := 0; timestamp < 1000; timestamp++ {
		//fmt.Printf("%s: generating data...\n", me)
		time.Sleep(100 * time.Millisecond)
		targetStream.Put(int(timestamp), float32(19.0+2*rand.Float32()))
	}
}

func merger(streams *Space, me string) {
	var source string
	var target string
	var sourceURI string
	var targetURI string
	var value float32
	var timestamp int

	// Get target streams
	sources, _ := streams.QueryAll(&source, me, &sourceURI)
	sourceStreams := make([]Space, len(sources))
	for i, source := range sources {
		sourceStreams[i] = NewRemoteSpace((source.GetFieldAt(2)).(string))
	}

	t, _ := streams.Query(me, &target, &targetURI)
	targetStream := NewRemoteSpace((t.GetFieldAt(2)).(string))

	for {
		var sum float32
		var n float32
		sum = 0
		n = 0
		var t Tuple
		//fmt.Printf("%s: waiting for data.\n", me)
		for _, source := range sourceStreams {
			if n == 0 {
				t, _ = source.Get(&timestamp, &value)
			} else {
				t, _ = source.Get(timestamp, &value)
			}
			timestamp = (t.GetFieldAt(0)).(int)
			value = (t.GetFieldAt(1)).(float32)
			sum += value
			n++
		}
		//fmt.Printf("%s: all data merged.\n", me)
		targetStream.Put(timestamp, sum/n)
	}
}

func duplicator(streams *Space, me string) {
	var source string
	var target string
	var sourceURI string
	var targetURI string
	var value float32
	var timestamp int

	//fmt.Printf("%s: getting source stream.\n", me)
	// Get source stream
	t, _ := streams.Query(&source, me, &sourceURI)
	sourceStream := NewRemoteSpace((t.GetFieldAt(2)).(string))

	//fmt.Printf("%s: getting target stream.\n", me)
	// Get target streams
	targets, _ := streams.QueryAll(me, &target, &targetURI)
	targetStreams := make([]Space, len(targets))
	for i, target := range targets {
		targetStreams[i] = NewRemoteSpace((target.GetFieldAt(2)).(string))
	}

	for {
		//fmt.Printf("%s: waiting for data...\n", me)
		t, _ := sourceStream.Get(&timestamp, &value)
		timestamp = (t.GetFieldAt(0)).(int)
		value = (t.GetFieldAt(1)).(float32)
		//fmt.Printf("%s: forwarding data data...\n", me)
		for _, targetStream := range targetStreams {
			targetStream.Put(timestamp, value)
		}
	}
}

func aggregator(streams *Space, me string) {
	var source string
	var target string
	var sourceURI string
	var targetURI string
	var timestamp int
	var value float32
	var sum float32
	const N = 100.0

	streams.Query(&source, me, &sourceURI)
	sourceStream := NewRemoteSpace(sourceURI)
	streams.Query(me, &target, &targetURI)
	targetStream := NewRemoteSpace(targetURI)

	for {
		sum = 0.0
		//fmt.Printf("%s: collecting enough data...\n", me)
		for i := 0; i < N; i++ {
			t, _ := sourceStream.Get(&timestamp, &value)
			timestamp = (t.GetFieldAt(0)).(int)
			value = (t.GetFieldAt(1)).(float32)
			sum += value
		}
		//fmt.Printf("%s: forwarding data...\n", me)
		targetStream.Put(timestamp, sum/N)
	}
}

func display(streams *Space, me string) {
	var source string
	var sourceURI string
	var timestamp int
	var value float32

	t, _ := streams.Query(&source, me, &sourceURI)
	sourceStream := NewRemoteSpace((t.GetFieldAt(2)).(string))
	for {
		t, _ := sourceStream.Get(&timestamp, &value)
		timestamp = (t.GetFieldAt(0)).(int)
		value = (t.GetFieldAt(1)).(float32)
		fmt.Printf("DISPLAY %s : %d - %f\n", me, timestamp, value)
	}
}
