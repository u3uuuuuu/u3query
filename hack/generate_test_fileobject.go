//用来生成方便调试测试使用的指定格式的二进制数据文件

package hack

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math/rand"
	"os"
	"path"
	"strconv"
	"syscall"
	"time"
	"u3.com/u3query/models"
	"u3.com/u3query/tree"
	"unsafe"
)

const SplitLength = models.SplitLength

var UnitList = make([]models.Unit, 0)

func unit2Buffer(unit *models.Unit, buffer *bytes.Buffer) {
	tmp := []byte{}
	keybytes := []byte(unit.Key)
	valuebytes := []byte(unit.Value)
	tmp = append(tmp, IntToBytes(unit.KeySize)...)
	tmp = append(tmp, keybytes...)
	tmp = append(tmp, IntToBytes(unit.ValueSize)...)
	tmp = append(tmp, valuebytes...)
	buffer.Write(tmp[:len(tmp):len(tmp)])
}


//整形转换成字节
func IntToBytes(n int) []byte {
	x := int64(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, x)
	return bytesBuffer.Bytes()
}

//字节转换成整形
func BytesToInt(b []byte) int {
	bytesBuffer := bytes.NewBuffer(b)

	var x int64
	binary.Read(bytesBuffer, binary.BigEndian, &x)

	return int(x)
}



func UnitWrite(unit *models.Unit, fp *os.File) (int64, error) {
	buf := new(bytes.Buffer)
	//err = binary.Write(buf, binary.LittleEndian, int64(23))
	unit2Buffer(unit, buf)
	// 将bytes写入文件
	n, err := fp.Write(buf.Bytes())
	if err != nil {
		return 0, err
	}
	err = fp.Sync()
	return int64(n), nil
}

func UnitRead(fp *os.File, offset int64) (int64, *models.Unit, error){
	unitKeySize := make([]byte, unsafe.Sizeof(int(0)))
	unitValueSize := make([]byte, unsafe.Sizeof(int(0)))
	n1, err := fp.ReadAt(unitKeySize, offset)
	if err != nil {
		return 0, nil, err
	}
	unitKey := make([]byte, BytesToInt(unitKeySize))
	n2, _ := fp.ReadAt(unitKey, offset+int64(n1))
	n3, _ := fp.ReadAt(unitValueSize, offset+int64(n1+n2))
	unitValue := make([]byte, BytesToInt(unitValueSize))
	n4, _ := fp.ReadAt(unitValue, offset+int64(n1+n2+n3))
	n := n1+n2+n3+n4
	u := models.Unit{BytesToInt(unitKeySize), string(unitKey), BytesToInt(unitValueSize), string(unitValue)}
	return int64(n), &u, nil
}

func randomString(length int) string {
	t := time.Now().Nanosecond()
	rand.Seed(int64(t))
	str := ""
	//numRange := [2]int{48,57}   // 0-9
	//charRange := [4]int{65,90,97,122}  //A-Z a-z
	//10+26+26=10+52=62
	for i := 0; i < length; i++ {
		randint := rand.Intn(62)
		if randint < 10 {
			randint += 48
		} else if randint >= 36 {
			randint += 61
		} else {
			randint += 55
		}
		str += string(randint)
	}
	return str
}

//生成随机Unit
func randomUnit() (unit *models.Unit) {
	t := time.Now().Nanosecond()
	rand.Seed(int64(t))
	keySize := rand.Intn(20)+5
	valueSize := rand.Intn(50)+5
	key := randomString(keySize)
	value := randomString(valueSize)
	unit = &models.Unit{keySize, key, valueSize, value}
	return
}

//生成测试用数据文件
func GenerateTestDataFile() error {
	os.Remove(path.Join("tests", "test1.binary"))
	os.RemoveAll("data")
	os.Mkdir("data", 0666)

	fp, err := os.OpenFile(path.Join("tests", "test1.binary"), syscall.O_RDWR|syscall.O_APPEND|syscall.O_CREAT, 0666)
	if err != nil {
		errstring := err.Error()
		return errors.New(errstring)
	}
	defer fp.Close()

	for i := 0; i < 100; i++ {
		u := randomUnit()
		_, err := UnitWrite(u, fp)
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}
	return nil
}


func ReadFile(filedir string, num int, off int64) (*[]*models.Unit, int64, error) {
	fp, err := os.Open(filedir)
	if err != nil {
		return nil, 0, err
	}
	defer fp.Close()

	return readFile(off, num, fp)
}


//读取测试文件并返回结构化数据
func readFile(off int64, num int, fp *os.File) (*[]*models.Unit, int64, error) {
	offset := off
	ret := []*models.Unit{}
	for i := 0; i < num; i++ {
		var unit *models.Unit
		n, unit,err := UnitRead(fp, offset)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return nil, 0, err
			}
		}
		offset += n
		ret = append(ret, unit)
	}
	return &ret, offset, nil
}


func GenerateBtree() error {
	primaryKey := 0
	offset := int64(0)
	//循环从测试文件中读取数据，然后生成B+树，并持久化
	for {
		//每次读取十万个
		testFile := path.Join("tests", "test1.binary")
		units, off, err := ReadFile(testFile, SplitLength, offset)
		if err != nil {
			return err
		}
		offset += off

		//读取units，然后生成b+树，并持久化
		storagefilename := strconv.Itoa(primaryKey)+"-"+strconv.Itoa(primaryKey+SplitLength)
		btree := tree.NewBTree()
		for _, u := range *units {
			btree.Insert(primaryKey, u)
			primaryKey++
		}
		_, err = tree.SaveToDisk(btree, storagefilename)
		if models.CacheBt.MaxPrimary < primaryKey {
			models.CacheBt.MaxPrimary = primaryKey
		}
		if err != nil {
			return err
		}
		models.CacheBt.Put(storagefilename, btree)
		if len(*units) < SplitLength {
			break
		}
	}
	return nil
}