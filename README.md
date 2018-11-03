````
package main

import (
	"os"
	"gopaths"
)

func main() {
	path := gopaths.PATH("Hello\\World\\text.txt")
	
	if path.Exists()  {
		fmt.Println("File exists")
	} else {
		f, _ := path.OpenWrite()
		
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
	path := "Hello\\World\\text.txt"
	
	_, err := os.Stat(path)
	
	if err == nil || os.IsExist(err) {
		fmt.Println("File exists")
	} else {
		f, _ = os.OpenWrite(path)	
		
		f.WriteString("Hello World")
		
		f.Close()
	}
}
```