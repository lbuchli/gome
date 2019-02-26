package graphics

import (
	"os"
	"testing"
)

func TestObjectFileReaderData(t *testing.T) {
	reader := OBJFileReader{}
	f, err := os.Open("/home/lukas/go/src/gitlocal/gome/testfiles/test1.obj")
	if err != nil {
		t.Fatal("Could not read object file")
	}

	t.Fatal(reader.Data(f))
}
