Currently only a experiment.

````
package main

import (
	"os"
	"gopaths"
)

func main() {
	fpath := PATH("Hello\\World\\text.txt")
	
	if fpath.Exists()  {
		fmt.Println("File exists")
	} else {
		f, _ := fpath.Create()
		
		f.WriteString("Hello World")
		
		f.Close()
	}
}
```

Instead Of

```
package main

import (
	"os"
)

func main() {
	fpath := "Hello\\World\\text.txt"
	
	_, err := os.Stat(fpath)
	
	if err == nil || os.IsExist(err) {
		fmt.Println("File exists")
	} else {
		f, err = os.Create(fpath)
		
		f.WriteString("Hello World")
		
		f.Close()
	}
}
```