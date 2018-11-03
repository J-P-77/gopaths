package gopaths

import (

	//"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestChangeExtension(t *testing.T) {
	type Test struct {
		path, change, result string
	}

	tests := []Test{
		Test{"youtube-files-dl.txt", ".go", "youtube-files-dl.go"},
		Test{"youtube-files-dl", "go", "youtube-files-dl"},
		Test{"youtube-files-dl", ".go", "youtube-files-dl"},
		Test{".youtube-files-dl", ".go", ".go"},
		Test{".", ".go", "."},
		Test{"youtube-files-dl.txt", "d", "youtube-files-dl.d"},
		Test{"youtube-files-dl.txt", ".d", "youtube-files-dl.d"},
	}

	for i, v := range tests {
		path := PATH(v.path)

		path = path.ChangeExtension(v.change)

		if path.ToString() != v.result {
			t.Errorf("%v [%v]\n", i, path.ToString())
		}
	}
}

func TestExists(t *testing.T) {
	path := PATH(".\\testdata")

	if !path.Exists() {
		t.Errorf("Path does not exists... %v", path.String())
	}
}

func TestIsDir(t *testing.T) {
	d := PATH(".\\testdata")

	if !d.IsDir() {
		t.Errorf("Path should be a directory... %v\n", d.String())
	}

	f := PATH(".\\testdata\\testfile.txt")

	if f.IsDir() {
		t.Errorf("Path should not be a directory... %v\n", d.String())
	}
}
func TestToAbsolutePath(t *testing.T) {
	type Test struct {
		path, result string
	}

	wd, err := os.Getwd()

	if err != nil {
		t.Errorf("Failed to get working directory... %v\n", err)
	}

	tests := []Test{
		Test{".\\testdata\\testfile.txt", wd + "\\testdata\\testfile.txt"},
	}

	for i, v := range tests {
		path := PATH(v.path)

		if path.ToAbsoluteString() != v.result {
			t.Errorf("%v [%v]\n", i, path.String())
		}
	}
}

func TestOpenRead(t *testing.T) {
	f := PATH(".\\testdata\\testfile.txt")

	rfile, err := f.OpenRead()
	defer rfile.Close()

	if err != nil {
		t.Error("Error Opening file " + f.ToAbsoluteString())
	}

	data, err := ioutil.ReadAll(rfile)

	if err != nil {
		t.Errorf("Failed reading file... %v\n", err)
	}

	const teststring = "This is just a test"
	if string(data) != teststring {
		t.Errorf("Failed reading file... [%v] does not match [%v]\n", err, teststring)
	}
}

func TestJoinString(t *testing.T) {
	type Test struct {
		paths  []string
		result string
	}

	tests := []Test{
		Test{[]string{"test/", "\\test\\", "\\test\\"}, "test\\test\\test\\"},
		Test{[]string{"\\test\\", "test/", "\\test\\"}, "\\test\\test\\test\\"},
		Test{[]string{"\\test\\", "test", "\\test/"}, "\\test\\test\\test\\"},
		Test{[]string{"c:\\test\\", "test", "\\test/"}, "c:\\test\\test\\test\\"},
		Test{[]string{"C:\\test", "test", "test//test"}, "C:\\test\\test\\test\\test"},
	}

	for i, v := range tests {
		path := PATH("")

		for _, p := range v.paths {
			path = path.JoinString(p)
		}

		if path.ToString() != v.result {
			t.Errorf("%v Test:%v [%v] != [%v]\n", i, v.paths, path.ToString(), v.result)
		}
	}
}
func TestJoinStrings(t *testing.T) {
	type Test struct {
		paths  []string
		result string
	}

	tests := []Test{
		Test{[]string{"test/", "\\test\\", "\\test\\"}, "test\\test\\test\\"},
		Test{[]string{"\\test\\", "test/", "\\test\\"}, "\\test\\test\\test\\"},
		Test{[]string{"\\test\\", "test", "\\test/"}, "\\test\\test\\test\\"},
		Test{[]string{"c:\\test\\", "test", "\\test/"}, "c:\\test\\test\\test\\"},
		Test{[]string{"C:\\test", "test", "test//test"}, "C:\\test\\test\\test\\test"},
	}

	for i, v := range tests {
		path := PATH("")

		path = path.JoinStrings(v.paths...)

		if path.ToString() != v.result {
			t.Errorf("%v Test:%v [%v] != [%v]\n", i, v.paths, path.ToString(), v.result)
		}
	}
}

func TestIterFileInfos(t *testing.T) {
	type Test struct {
		path   string
		result []string
	}

	tests := []Test{
		Test{".\\testdata\\testdir", []string{"1", "2.gotest", "3.txt"}},
	}

	for i, v := range tests {
		p := PATH(v.path)

		counter := 0
		for next, hasNext := p.IterFileInfos(); hasNext(); {
			name := next().Name()

			if name == v.result[counter] {
				t.Log(name)
				counter++
			}
		}

		if counter != len(v.result) {
			t.Errorf("%v Test:%v Failed\n", i, p.String())
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

	items := p1.WalkDirPath(func(PATH) (bool, bool) {
		return false, true
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

func BenchmarkList(b *testing.B) {
	p := PATH("C:\\Users\\Justin")

	b.ResetTimer()
	for x := 0; x < b.N; x++ {
		p.List()
	}
}
