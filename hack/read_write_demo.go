package hack

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path"
	"time"
	"unsafe"
)

//纯文本写
func txtWrite() {

	// 1. create bytes with string
	data := []byte("hello world\n")

	// 2. create file to write
	file, _ := os.Create("data.txt")
	defer file.Close()
	bytes, _ := file.Write(data)

	fmt.Printf("Wrote %d bytes to file \n", bytes)
}

//纯文本读
func txtRead() {

	// 1. open file to read
	file, _ := os.Open("data.txt")
	defer file.Close()

	// 2. read all bytes
	stats, _ := file.Stat()

	var size int64 = stats.Size()
	bytes := make([]byte, size)

	bufr := bufio.NewReader(file)
	n, _ := bufr.Read(bytes)

	// 3. convert bytes to string
	fmt.Println(string(bytes), n)
}

//二进制文件写
func binWrite() {
	t := time.Now().Nanosecond()
	fp, _ := os.Create(path.Join("bin", "numbers.binary"))
	defer fp.Close()

	rand.Seed(int64(t))

	buf := new(bytes.Buffer)
	for i := 0; i < 10; i++ {
		binary.Write(buf, binary.LittleEndian, int32(i))
		fp.Write(buf.Bytes())
	}

	// bin file contains: 0~9
}

//二进制文件读
func binRead() {
	fp, _ := os.Open(path.Join("bin", "numbers.binary"))
	defer fp.Close()

	data := make([]byte, 4)
	var k int32
	for {
		data = data[:cap(data)]

		// read bytes to slice
		n, err := fp.Read(data)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			break
		}

		// convert bytes to int32
		data = data[:n]
		binary.Read(bytes.NewBuffer(data), binary.LittleEndian, &k)

		fmt.Println(k)
	}
}

//结构体示例
type MyData struct {
	_ [1]byte
	Y int32
	X int32
	Z int32
}

func structWrite() {
	fp, _ := os.Create(path.Join("bin", "struct.binary"))
	defer fp.Close()

	// 将结构体转成bytes, 按照字段的声明顺序，但是"_"被放在最后
	data := &MyData{X:1, Y:2, Z:3}
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, data)

	// 将bytes写入文件
	fp.Write(buf.Bytes())
	fp.Sync()
}

func structRead() {
	fp, _ := os.Open(path.Join("bin", "struct.binary"))
	defer fp.Close()

	// 创建byte slice, 以读取bytes. 此处MyData的size为16，因为有字节对齐
	dataBytes := make([]byte, unsafe.Sizeof(MyData{}))
	data := MyData{}
	n, _ := fp.Read(dataBytes)
	dataBytes = dataBytes[:n]

	// 将bytes转成对应的struct
	binary.Read(bytes.NewBuffer(dataBytes), binary.LittleEndian, &data)
	fmt.Println(data)
}
