package gopaths

import (
	//"fmt"
	"io/ioutil"
	"os"
	//"path"
	"path/filepath"
	"strings"
)

type PATH []byte //rune

const (
	_S_  = os.PathSeparator
	_LS_ = os.PathListSeparator
)

func (p *PATH) Create() (file *os.File, err error) {
	return os.Create(string(*p))
}

func (p *PATH) Open() (file *os.File, err error) {
	return os.Open(string(*p))
}

func (p *PATH) OpenFile(flag int, perm os.FileMode) (file *os.File, err error) {
	return os.OpenFile(string(*p), flag, perm)
}

func (p *PATH) JoinString(paths ...string) PATH {
	s := filepath.Join(paths...)

	return PATH(filepath.Join(string(*p), s))
}

func (p *PATH) Join(paths ...PATH) PATH {
	tmppath := make([]string, len(paths)+1)

	tmppath = append(tmppath, string(*p))

	for x, _ := range paths {
		tmppath = append(tmppath, string(paths[x]))
	}

	return PATH(filepath.Join(tmppath...))
}

func (p *PATH) ToString() string {
	return string(*p)
}

func (p *PATH) ToAbsoluteString() string {
	tmppath, _ := filepath.Abs(string(*p))

	return tmppath
}

func (p *PATH) ToAbsolutePath() PATH {
	tmppath, _ := filepath.Abs(string(*p))
	//println(tmppath)
	return PATH(tmppath)
}

func (p *PATH) Clean() PATH {
	return PATH(filepath.Clean(string(*p)))
}

func (p *PATH) MkdirDefault() error {
	return os.Mkdir(string(*p), 0666)
}

func (p *PATH) Mkdir(perm os.FileMode) error {
	return os.Mkdir(string(*p), perm)
}

func (p *PATH) MkdirAllDefault() error {
	return os.MkdirAll(string(*p), 0666)
}

func (p *PATH) MkdirAll(perm os.FileMode) error {
	return os.MkdirAll(string(*p), perm)
}

func (p *PATH) MkdirAllParentDefault() error {
	str := string(*p)

	return os.MkdirAll(str[:len(p.Name())], 0666)
}

func (p *PATH) MkdirAllParent(perm os.FileMode) error {
	str := string(*p)

	return os.MkdirAll(str[:len(p.Name())], perm)
}

func (p *PATH) Exists() bool {
	_, err := os.Stat(string(*p))

	return err == nil || os.IsExist(err)
}

func (p *PATH) IsDir() bool {
	stat, err := os.Stat(string(*p))

	if err == nil || os.IsExist(err) {
		return stat.IsDir()
	}

	return false
}

func (p *PATH) List() []PATH {
	if p.IsDir() {
		finfo, err := ioutil.ReadDir(string(*p))

		if err == nil {
			paths := make([]PATH, len(finfo))

			for x, _ := range finfo {
				paths[x] = PATH(JoinString(string(*p), finfo[x].Name()))
			}

			return paths
		}
	}

	return make([]PATH, 0)
}

func (p *PATH) ListNames() []PATH {
	if p.IsDir() {
		finfo, err := ioutil.ReadDir(string(*p))

		if err == nil {

			paths := make([]PATH, len(finfo))

			for x, _ := range finfo {
				paths[x] = PATH(finfo[x].Name())
			}

			return paths
		}
	}

	return make([]PATH, 0)
}

func (p *PATH) ListStringNames() []string {
	if p.IsDir() {
		ofile, err := p.Open()
		defer ofile.Close()

		if err == nil {
			names, err := ofile.Readdirnames(0)

			if err == nil {
				return names
			}
		}
	}

	return make([]string, 0)
}

func (p *PATH) ListInfo() []os.FileInfo {
	if p.IsDir() {
		finfo, err := ioutil.ReadDir(string(*p))

		if err != nil {
			return make([]os.FileInfo, 0)
		}

		return finfo
	} else {
		return make([]os.FileInfo, 0)
	}
}

func (p *PATH) IterFileInfos() (func() os.FileInfo, func() bool) {
	var next os.FileInfo = nil

	if p.IsDir() {
		ofile, _ := p.Open()
		//Next()
		return func() os.FileInfo {
				return next
			},
			//HasMore()
			func() bool {
				stat, err := ofile.Readdir(1)

				if err != nil {
					next = nil
					ofile.Close()
				} else {
					next = stat[0]
				}

				return next != nil
			}
	} else {
		return func() os.FileInfo { return nil }, func() bool { return false }
	}
}

func (p *PATH) VolumeName() string {
	return filepath.VolumeName(string(*p))
}

func (p *PATH) Name() string {
	return filepath.Base(string(*p))
}

func (p *PATH) Extension() string {
	return filepath.Ext((string(*p)))
}

func (p *PATH) NameWithoutExtension() string {
	base := p.Name()

	index := strings.LastIndex(base, ".")

	if index == -1 {
		return base
	} else {
		return base[:index]
	}
}

func (p *PATH) Size() int64 {
	stat, err := os.Stat(string(*p))

	if err != nil {
		return -1
	}

	return stat.Size()
}

func (p *PATH) WalkDirAll() []PATH {
	return p.WalkDirPath(func(PATH) (bool, bool) {
		return true, false
	}, func(PATH) bool {
		return true
	})
}

func (p *PATH) WalkDirInclude(acceptfile func(PATH) bool) []PATH {
	return p.WalkDirPath(func(PATH) (bool, bool) {
		return true, false
	}, acceptfile)
}

func (p *PATH) WalkDirPath(acceptdir func(PATH) (bool, bool), acceptfile func(PATH) bool) []PATH {
	list := make([]PATH, 0, 10)

	for _, n := range p.ListStringNames() {
		p := JoinString(string(*p), n)

		if !p.Exists() {
			continue
		}

		if p.IsDir() {
			if walk, include := acceptdir(p); walk {
				if include {
					list = append(list, p)
				}

				files := p.WalkDirPath(acceptdir, acceptfile)

				if len(files) > 0 {
					list = append(list, files...)
				}
			}
		} else {
			if acceptfile(p) {
				list = append(list, p)
			}
		}
	}

	return list
}

func (p *PATH) WalkDirDoAll(do func(path PATH)) {
	p.WalkDirDo(func(PATH) (bool, bool) {
		return true, false
	}, func(PATH) bool {
		return true
	}, do)
}

func (p *PATH) WalkDirDoInclude(acceptfile func(PATH) bool, do func(path PATH)) {
	p.WalkDirDo(func(PATH) (bool, bool) {
		return true, false
	}, acceptfile, do)
}

func (p *PATH) WalkDirDo(acceptdir func(PATH) (bool, bool), acceptfile func(PATH) bool, do func(path PATH)) {
	for _, n := range p.ListStringNames() {
		p := JoinString(string(*p), n)

		if !p.Exists() {
			continue
		}

		if p.IsDir() {
			if walk, include := acceptdir(p); walk {
				if include {
					do(p)
				}

				p.WalkDirDo(acceptdir, acceptfile, do)
			}
		} else {
			if acceptfile(p) {
				do(p)
			}
		}
	}
}

func JoinString(paths ...string) PATH {
	return PATH(filepath.Join(paths...))
}
