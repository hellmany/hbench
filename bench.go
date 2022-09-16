package hbench

import (
	"bufio"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	bytesize "github.com/inhies/go-bytesize"
	"github.com/remeh/sizedwaitgroup"
)

var allFiles []string
var files []string

func Bench(c ConfData) (RJson, error) {

	//      pid := os.Getpid()
	//      syscall.Setpriority(syscall.PRIO_PROCESS, pid, -19)
	if c.DebugInfo {
		log.Println("Path", c.Path)
		log.Println("Threads", c.Threads)
		log.Println("Files", c.Max)
	}
	if c.LimitMax > 0 && c.LimitMax < c.Max {
		c.LimitMax = c.Max
	}

	swg := sizedwaitgroup.New(c.Threads)
	i := 0

	err := filepath.Walk(c.Path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() || (len(c.Extension) > 0 && !strings.Contains(path, c.Extension)) {
				return nil
			}

			i++
			if c.LimitMax > 0 && i > c.LimitMax {
				return io.EOF
			}
			files = append(files, path)
			return nil
		})
	if err != nil && err != io.EOF {
		if c.DebugInfo {
			log.Println("Error walk", err)
		}
		return RJson{}, err
	}

	var total_bytes uint64
	var total_files uint64
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(files), func(i, j int) { files[i], files[j] = files[j], files[i] })
	if c.DebugInfo {
		log.Println("Readed", len(files), "paths...")
	}

	if c.Max > 0 && len(files) > c.Max {
		files = files[:c.Max]
	}

	start := time.Now()
	if c.DebugInfo {
		log.Println("To process", len(files), "paths...")
	}

	for ci := 0; ci < c.Inter; ci++ {
		for _, path := range files {
			swg.Add()

			go func(path string) {
				if c.DebugInfo {
					log.Println("Reading", path, ci)
				}
				defer swg.Done()

				bytes, _ := readFile(path, c.Size*1024)

				if c.DebugInfo {
					log.Println("Readed", path, ci, total_files, bytes, "bytes")
				}
				atomic.AddUint64(&total_files, 1)
				atomic.AddUint64(&total_bytes, bytes)
				return

			}(path)
		}
	}

	swg.Wait()

	elapsed := time.Since(start)

	b := bytesize.New(float64(total_bytes))
	megabytes := b.Format("%.2f ", "megabyte", true)
	gigabytes := b.Format("%.2f ", "gigabyte", true)

	speed := float64(total_bytes) / elapsed.Seconds()
	if c.DebugInfo {
		log.Println("Readed", total_files, "files and", megabytes, " (", gigabytes, ") in", c.Threads, "threads")
		log.Printf("Speed: %.2f mb/s", speed/1024/1024)
		log.Println("Took", elapsed, "(", elapsed.Seconds(), "seconds)")
	}
	r := RJson{
		Threads:  c.Threads,
		Bytes:    total_bytes,
		Files:    total_files,
		Seconds:  elapsed.Seconds(),
		TimeStr:  elapsed.String(),
		SpeedMBs: speed / 1024 / 1024,
	}
	return r, nil
}

func readFile(p string, size int) (uint64, error) {

	f, err := os.Open(p)
	if err != nil {
		log.Println("Error reading", p, err)
		return 0, err
	}

	defer f.Close()
	var total_bytes uint64

	r := bufio.NewReader(f)
	b := make([]byte, 1024)
	t := 0
	for {
		bytes, err := r.Read(b)
		if err != nil {
			return total_bytes, err
		}
		t++
		total_bytes += uint64(bytes)
		if size > 0 && t >= size {
			return total_bytes, nil
		}
	}
	return total_bytes, nil
}
