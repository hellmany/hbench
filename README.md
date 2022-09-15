This library was written to test reading in multiple threads and the speed of different file systems and raid levels, since the current ones were all running in a single thread.

In response to the read and write structure of the

	type RJson struct {
		Threads int
		Bytes uint64
		Files uint64
		Seconds float64
		TimeStr string
		SpeedMBs float64
	}

Or output to console 


Compiled utilities in the repository github.com/hellmany/hbench-cli

