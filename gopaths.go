package gopaths

import (
	"os"
	"path"
	"path/filepath"
)

type PATH []byte

func (p *PATH) Exists() bool {
	_, err := os.Stat(string(*p))
	
	if err == nil {return true}
	
	return os.IsNotExist(err)
}

func(p *PATH) Create() (file *os.File, err error) {
	file, err = os.Create(string(*p))
	
	return file, err
}

func(p *PATH) Open() (file *os.File, err error) {
	file, err = os.Open(string(*p))
	
	return file, err
}

func(p *PATH) Join(paths ...string) PATH {
	tmppath := make([]string, len(paths) + 1)
	
	tmppath = append(tmppath, string(*p))
	
	for x, _ := range paths {
		tmppath = append(tmppath, paths[x])
	}
	
	return PATH(path.Join(tmppath...))
}

func(p *PATH) ToString() string {
	return string(*p)
}

func(p *PATH) ToAbsoluteString() string {
	tmppath, _ := filepath.Abs(string(*p))
	
	return tmppath
}

func(p *PATH) Clean() PATH {
	return PATH(path.Clean(string(*p)))
}

func(p *PATH) ToAbsolutePath() PATH {
	tmppath, _ := filepath.Abs(string(*p))
	
	return PATH(tmppath)
}