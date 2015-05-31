package gopaths

import (
	"os"
	"path"
	"path/filepath"
)

type PATH []byte

func (p *PATH) Exists() bool {
	_, err := os.Stat(string(*p))
	
	return err == nil || os.IsExist(err)
}

func(p *PATH) Create() (file *os.File, err error) {
	file, err = os.Create(string(*p))
	
	return file, err
}

func(p *PATH) Open() (file *os.File, err error) {
	file, err = os.Open(string(*p))
	
	return file, err
}

func(p *PATH) OpenFile(flag int, perm os.FileMode) (file *os.File, err error) {
	file, err = os.OpenFile(string(*p), flag, perm)
	
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

func(p *PATH) Mkdir(perm os.FileMode) error {
	return os.Mkdir(string(*p), perm)
}

func(p *PATH) MkdirAll(perm os.FileMode) error {
	return os.MkdirAll(string(*p), perm)
}

func(p *PATH) IsDir() bool {
	stat, err := os.Stat(string(*p))
	
	if err == nil || os.IsExist(err) {
		return stat.IsDir()
	} 
	
	return false
}

/*func(p *PATH) list() []PATH {

}*/