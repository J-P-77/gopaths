package gopaths

import (
	"os"
	//"fmt"
	"testing"
)

func TestToAbsolutePath(t *testing.T) {
	f := PATH("HelloWorld\\Justin")
	
	wd, _ := os.Getwd()
	
	if f.ToAbsoluteString() != wd + "\\HelloWorld\\Justin" {
		t.Error("Paths Not Equal")
	}
}