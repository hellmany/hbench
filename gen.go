package hbench

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync/atomic"
	"time"

	bytesize "github.com/inhies/go-bytesize"
	"github.com/remeh/sizedwaitgroup"
)

var letterRunes = []rune("abcdefghijk")

//var letterRunesBig = []rune("abcdefghijkl")

func Gen(c ConfData) {

	rand.Seed(time.Now().UnixNano())

	var total_bytes uint64
	var total_files uint64
	start := time.Now()

	fmt.Println("Path", c.Path)
	fmt.Println("Threads", c.Threads)
	fmt.Println("Files", c.Max)

	swg := sizedwaitgroup.New(c.Threads)

	for i := 1; i <= c.Max; i++ {
		swg.Add()
		dir := c.Path + "/" + RandStringRunes(2) + "/" + RandStringRunes(2)
		file := RandStringRunes(10) + ".test.mp4"
		go func() {
			defer swg.Done()

			bytes := createfile(dir, file, c.Size+rand.Intn(c.RandSize), i)
			atomic.AddUint64(&total_bytes, bytes)
			atomic.AddUint64(&total_files, 1)

		}()
	}

	swg.Wait()

	elapsed := time.Since(start)

	b := bytesize.New(float64(total_bytes))
	megabytes := b.Format("%.2f ", "megabyte", true)
	gigabytes := b.Format("%.2f ", "gigabyte", true)

	speed := float64(total_bytes) / elapsed.Seconds()

	log.Println("Writed", total_files, "files and", megabytes, " (", gigabytes, ") in", c.Threads, "threads")
	log.Printf("Speed: %.2f mb/s", speed/1024/1024)
	log.Println("Took", elapsed, "(", elapsed.Seconds(), ") seconds")

}
func createfile(dir string, file string, size int, c int) uint64 {

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.MkdirAll(dir, os.ModePerm)
	}
	fmt.Println("Writing", dir+"/"+file, size, "mb", c)
	f, _ := os.OpenFile(dir+"/"+file, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0655)
	defer f.Close()

	var total_bytes uint64
	for i := 1; i <= size*1024; i++ {
		d := make([]byte, 1024)
		rand.Read(d)

		//fmt.Println("Writing", dir+"/"+file, i, "mb")
		//return
		//_ := os.WriteFile(dir+"/"+path, d, 0644)
		bytes, _ := f.Write(d)
		total_bytes += uint64(bytes)
	}
	fmt.Println("Writed", dir+"/"+file, size, "mb", c)
	return total_bytes
}

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
