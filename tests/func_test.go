package test

import (
	"fmt"
	"os"
	"path"
	"testing"
	"u3.com/u3query/hack"
	"u3.com/u3query/models"
)

func TestFile(t *testing.T) {
	hack.GenerateTestDataFile()
	offset := int64(0)

	fp, err := os.Open(path.Join("tests", "test1.binary"))
	if err != nil {
		t.Error(err.Error())
	}
	defer fp.Close()

	for i := 0; i < 100; i++ {
		var unit *models.Unit
		n, unit,_ := hack.UnitRead(fp, offset)
		offset += n
		fmt.Println(unit)
	}
}
