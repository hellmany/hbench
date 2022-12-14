package hbench

type ConfData struct {
	Threads   int
	LimitMax  int
	Max       int
	Size      int
	RandSize  int
	Path      string
	Inter     int
	DebugInfo bool
	Extension string
}

type RJson struct {
	Threads  int
	Bytes    uint64
	Files    uint64
	Seconds  float64
	TimeStr  string
	SpeedMBs float64
}
