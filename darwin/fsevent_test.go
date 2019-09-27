package darwin_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	// Frameworks
	mac "github.com/djthorpe/gopi-hw/darwin"
)

func Test_FSStream_000(t *testing.T) {
	if stream, err := mac.NewEventStream([]string{"/"}, 0, 0, 0, 0); err != nil {
		t.Error(err)
	} else {
		defer mac.ReleaseEventStream(stream)
	}
}

func Test_FSStream_001(t *testing.T) {
	stop := make(chan struct{})
	go func(stop chan struct{}) {
		// wait for 1 second and then stop
		time.Sleep(time.Second)
		close(stop)
	}(stop)
	t.Log("Start RL")
	if err := mac.RunLoopInCurrentThread(stop); err != nil {
		t.Error(err)
	}
	t.Log("Stop RL")
}

func Test_FSStream_002(t *testing.T) {
	// Create a stream and add to run loop
	latest := mac.LatestEventID()
	if stream, err := mac.NewEventStream([]string{"/tmp"}, 0, latest, 0, mac.FS_STREAM_CREATE_FLAG_FILEEVENTS); err != nil {
		t.Error(err)
	} else {
		defer mac.ReleaseEventStream(stream)
		runloop := mac.CurrentRunLoop()
		mac.StartEventStreamInRunloop(stream, runloop)
		defer mac.StopEventStreamInRunLoop(stream, runloop)

		stop := make(chan struct{})

		// Background thread will stop later
		go func(stop chan struct{}) {
			// wait for 2 seconds and then stop
			time.Sleep(time.Second * 2)
			close(stop)
		}(stop)

		t.Log("Start RL")
		if err := mac.RunLoopInCurrentThread(stop); err != nil {
			t.Error(err)
		}
		t.Log("Stop RL")

	}
}
func Test_FSStream_003(t *testing.T) {
	for f := mac.FS_STREAM_FLAG_MIN; f <= mac.FS_STREAM_FLAG_MAX; f <<= 1 {
		str := fmt.Sprint(f)
		if strings.HasPrefix(str, "FS_STREAM_FLAG_") {
			t.Logf("%08X => %v", uint(f), str)
		} else {
			t.Errorf("Expecting prefix FS_STREAM_FLAG_: %08X => %v", uint(f), str)
		}
	}
}
func Test_FSStream_004(t *testing.T) {
	// Create two streams and add to run loop
	latest := mac.LatestEventID()
	if stream1, err := mac.NewEventStream([]string{"/tmp"}, 0, latest, 0, mac.FS_STREAM_CREATE_FLAG_FILEEVENTS); err != nil {
		t.Error(err)
	} else if stream2, err := mac.NewEventStream([]string{"/tmp"}, 0, latest, 0, mac.FS_STREAM_CREATE_FLAG_FILEEVENTS); err != nil {
		t.Error(err)
	} else {

		defer mac.ReleaseEventStream(stream1)
		defer mac.ReleaseEventStream(stream2)

		runloop := mac.CurrentRunLoop()
		mac.StartEventStreamInRunloop(stream1, runloop)
		mac.StartEventStreamInRunloop(stream2, runloop)
		defer mac.StopEventStreamInRunLoop(stream1, runloop)
		defer mac.StopEventStreamInRunLoop(stream2, runloop)

		stop := make(chan struct{})

		// Callback
		mac.SetEventCallback(0, func(evt *mac.FSEvent) {
			t.Log(evt)
		})

		// Background thread will stop later
		go func(stop chan struct{}) {
			// wait for 1 second
			time.Sleep(time.Second * 1)

			// Create a file
			if file, err := os.Create("/tmp/test_file"); err != nil {
				t.Error(err)
			} else if err := file.Close(); err != nil {
				t.Error(err)
			}

			// wait for 1 second and then remove file
			time.Sleep(time.Second * 1)

			if err := os.Remove("/tmp/test_file"); err != nil {
				t.Error(err)
			}

			// wait for 1 second and then stop
			time.Sleep(time.Second * 1)

			close(stop)
		}(stop)

		t.Log("Start RL")
		if err := mac.RunLoopInCurrentThread(stop); err != nil {
			t.Error(err)
		}
		t.Log("Stop RL")

	}
}
