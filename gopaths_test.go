package gopaths

import (
	//"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"
)

func exists(path PATH, t *testing.T) {
	if !path.Exists() {
		t.Error("File does not exists " + path.ToAbsoluteString())
	}
}

func TestExists(t *testing.T) {
	d := PATH("C:\\Users\\Justin\\temp-dir")

	if !d.Exists() {
		t.Error("File does not exists " + d.ToAbsoluteString())
	}
}

func TestIsDir(t *testing.T) {
	d := PATH("C:\\Users\\Justin\\temp-dir")

	if d.IsDir() == false {
		t.Error("Path should be a directory")
	}

	f := PATH("C:\\Users\\Justin\\mac addresses.txt")

	if f.IsDir() == true {
		t.Error("Path should not be a directory")
	}

	t.Log("[isDir=" + strconv.FormatBool(d.IsDir()) + "] " + d.ToAbsoluteString())
	t.Log("[isDir=" + strconv.FormatBool(f.IsDir()) + "] " + f.ToAbsoluteString())

}
func TestToAbsolutePath(t *testing.T) {
	f := PATH("HelloWorld\\Justin")

	wd, _ := os.Getwd()

	if f.ToAbsoluteString() != wd+"\\HelloWorld\\Justin" {
		t.Error("Paths Not Equal")
	}
}

func TestOpenFile(t *testing.T) {
	f := PATH("gopaths_test.go")

	exists(f, t)

	ofile, err := f.OpenFile(os.O_RDONLY, 0)
	defer ofile.Close()

	if err != nil {
		t.Error("Error Opening file " + f.ToAbsoluteString())
	} else {
		_, err := ioutil.ReadAll(ofile)

		if err != nil {
			t.Error(err)
		}
	}
}

func TestJoinString(t *testing.T) {
	f := PATH("HelloWorld\\Justin")

	f = f.JoinString("gopaths_test.go")

	if f.ToString() != "HelloWorld\\Justin\\gopaths_test.go" {
		t.Error(f.ToString() + " != HelloWorld\\Justin\\gopaths_test.go")
	}
}

func TestJoin(t *testing.T) {
	a := PATH("HelloWorld\\Justin")
	b := PATH("gopaths_test.go")

	c := a.Join(b)

	if c.ToString() != "HelloWorld\\Justin\\gopaths_test.go" {
		t.Error(c.ToString() + " != HelloWorld\\Justin\\gopaths_test.go")
	}
}

func TestIterFileInfos(t *testing.T) {
	p := PATH("./")
	//fmt.Println(p.ToString())

	for next, hasnext := p.IterFileInfos(); hasnext(); {
		name := next().Name()

		t.Log(name)

		switch name {
		case "gopaths.go":
		case "gopaths_test.go":
		case "README.md":
		case "setup.bat":
		case "test.bat":
		default:
			t.Error()
		}
	}
}

func TestVolumeName(t *testing.T) {
	p1 := PATH("C:\\Users\\Justin")

	if p1.VolumeName() != "C:" {
		t.Error("Not C drive", p1.VolumeName())
	}

	p2 := PATH("./")
	p2 = p2.ToAbsolutePath()
	if p2.VolumeName() != "C:" {
		t.Error("Not C drive", p2.VolumeName())
	}
}

func TestName(t *testing.T) {
	p1 := PATH("C:/Users/Justin")

	if p1.Name() == "Justin" {
		t.Log(p1.Name())
	} else {
		t.Fail()
	}
}

func TestWalkDirPath(t *testing.T) {
	goal := 368
	p1 := PATH("C:\\Users\\Justin\\temp-dir\\build")

	items := p1.WalkDirPath(func(PATH) bool {
		return true
	},
		func(PATH) bool {
			return true
		})

	if len(items) != goal {
		for _, i := range items {
			t.Log(string(i))
		}
		t.Logf("Count: %v Goal was %v", len(items), goal)

		t.Fail()
	}
}

func BenchmarkListString(b *testing.B) {
	/*p := PATH("C:\\Users\\Justin")

	b.ResetTimer()
	for x := 0; x < b.N; x++ {
		p.ListString()
	}*/
}

func BenchmarkList(b *testing.B) {
	p := PATH("C:\\Users\\Justin")

	b.ResetTimer()
	for x := 0; x < b.N; x++ {
		p.List()
	}
}
