package gopaths

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type PATH []rune

const (
	_S_  = os.PathSeparator
	_LS_ = os.PathListSeparator
)

// Opens file for appending to the file.
// Returns
//		*os.File 	Used for reading/writing to file
//		error 		If error, it will be of type *PathError
func (p *PATH) OpenAppend() (file *os.File, err error) {
	return os.OpenFile(string(*p), os.O_APPEND|os.O_CREATE, 0777)
}

// Opens file for reading from the file.
// Returns
//		*os.File 	Used for reading/writing to file
//		error 		If error, it will be of type *PathError
func (p *PATH) OpenRead() (file *os.File, err error) {
	return os.Open(string(*p))
}

// Opens file for writing to the file.
// Returns
//		*os.File 	Used for reading/writing to file
//		error 		If error, it will be of type *PathError
func (p *PATH) OpenWrite() (file *os.File, err error) {
	return os.OpenFile(string(*p), os.O_WRONLY|os.O_CREATE, 0777)
}

// Opens file for reading and writing to the file.
// Returns
//		*os.File 	Used for reading/writing to file
//		error 		If error, it will be of type *PathError
func (p *PATH) OpenReadWrite() (file *os.File, err error) {
	return os.OpenFile(string(*p), os.O_RDWR|os.O_CREATE, 0777)
}

// Opens file for reading and writing to the file.
//		int			Open file, for reading, writing, appending, and creating....
//		os.FileMode The file permissions
// Returns
//		*os.File 	Used for reading/writing to file
//		error 		If error, it will be of type *PathError
func (p *PATH) OpenFile(flag int, perm os.FileMode) (file *os.File, err error) {
	return os.OpenFile(string(*p), flag, perm)
}

//Joins paths, adds a path separator if required.
//		[]string	an array of path to joins
//Returns
//		PATH		New path after joining paths
func (p *PATH) JoinStrings(paths ...string) PATH {
	//s := filepath.Join(paths...)
	//return PATH(filepath.Join(string(*p), s))

	strPath := strings.Builder{}
	sep := false

	isSep := func(c rune) bool {
		if c == '\\' || c == '/' {
			if !sep {
				strPath.WriteRune(_S_)
			}
			return true
		}

		strPath.WriteRune(c)
		return false
	}

	//append current PATH
	for _, c := range *p {
		sep = isSep(rune(c))
	}

	for i, v := range paths {
		for _, c := range v {
			sep = isSep(c)
		}

		if !sep {
			if i != len(paths)-1 {
				strPath.WriteRune(_S_)
				sep = true
			}
		}
	}

	return PATH(strPath.String())
}

func (p *PATH) JoinString(path string) PATH {
	//s := filepath.Join(paths...)
	//return PATH(filepath.Join(string(*p), s))

	strPath := &strings.Builder{}
	sep := false

	isSep := func(c rune) bool {
		if c == '\\' || c == '/' {
			if !sep {
				strPath.WriteRune(_S_)
			}
			return true
		}

		strPath.WriteRune(c)
		return false
	}

	//append current PATH
	for _, c := range *p {
		sep = isSep(rune(c))
	}

	if !sep && len(*p) > 0 {
		strPath.WriteRune(_S_)
		sep = true
	}

	for _, c := range path {
		sep = isSep(c)
	}

	return PATH(strPath.String())
}

//Joins paths, adds a path separator if required.
//		PATH	path to join
//Returns
//		PATH		New path after joining paths
func (p *PATH) Join(path PATH) PATH {
	return p.JoinString(string(path))
}

//Return the raw string of the path
func (p *PATH) ToString() string {
	return string(*p)
}

//Return the absolute path string
func (p *PATH) String() string {
	return p.ToAbsoluteString()
}

func (p *PATH) ToAbsoluteString() string {
	validate(p)

	tmppath, _ := filepath.Abs(string(*p))

	return tmppath
}

func (p *PATH) ToAbsolutePath() PATH {
	return PATH(p.ToAbsoluteString())
}

func (p *PATH) Clean() PATH {
	validate(p)

	return PATH(filepath.Clean(string(*p)))
}

func (p *PATH) MkdirDefault() error {
	validate(p)

	return os.Mkdir(string(*p), 0666)
}

func (p *PATH) Mkdir(perm os.FileMode) error {
	validate(p)

	return os.Mkdir(string(*p), perm)
}

func (p *PATH) MkdirAllDefault() error {
	validate(p)

	return os.MkdirAll(string(*p), 0666)
}

func (p *PATH) MkdirAll(perm os.FileMode) error {
	validate(p)

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
	validate(p)

	_, err := os.Stat(string(*p))

	return err == nil || os.IsExist(err)
}

func (p *PATH) IsDir() bool {
	validate(p)

	stat, err := os.Stat(string(*p))

	if err == nil || os.IsExist(err) {
		return stat.IsDir()
	}

	return false
}

func (p *PATH) IsFile() bool {
	validate(p)

	stat, err := os.Stat(string(*p))

	if err == nil || os.IsExist(err) {
		return !stat.IsDir()
	}

	return false
}

func (p *PATH) List() []PATH {
	if p.IsDir() {
		finfo, err := ioutil.ReadDir(string(*p))

		if err == nil {
			paths := make([]PATH, len(finfo))

			for x := range finfo {
				paths[x] = PATH(JoinPaths(string(*p), finfo[x].Name()))
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

			for x := range finfo {
				paths[x] = PATH(finfo[x].Name())
			}

			return paths
		}
	}

	return make([]PATH, 0)
}

func (p *PATH) ListStringNames() []string {
	if p.IsDir() {
		ofile, err := os.Open(string(*p))
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
	var next os.FileInfo

	if p.IsDir() {
		ofile, _ := os.Open(string(*p))
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
	}

	return func() os.FileInfo { return nil }, func() bool { return false }
}

func (p *PATH) VolumeName() string {
	validate(p)

	return filepath.VolumeName(string(*p))
}

func (p *PATH) Name() string {
	validate(p)

	return filepath.Base(string(*p))
}

func (p *PATH) Extension() string {
	validate(p)

	ext := filepath.Ext((string(*p)))

	//remove . charater
	return ext[1:]
}

/*
func (p *PATH) ChangeExtension2(ext string) PATH {
	if len(ext) == 0 {
		return *p
	}

	offset := len(*p) - 1
	for ; offset >= 0; offset-- {
		if (*p)[offset] == '.' {
			break
		}
	}

	if offset == -1 || offset == len(*p) {
		return *p
	}

	length := offset + len(ext)
	extOffset := 0

	if ext[0] == '.' {
		//Skips . at the start of ext
		extOffset++
	} else {
		//ext does not start with . so add extra one to the length
		length++
	}

	t := make(PATH, length)

	i := 0
	for ; i < offset; i++ {
		t[i] = (*p)[i]
	}

	t[i] = '.'
	i++

	//fmt.Printf("%v\n", string(t))

	for ; extOffset < len(ext); extOffset++ {
		t[i] = ext[extOffset]
		i++
	}

	return t
}
*/

func (p *PATH) ChangeExtension(ext string) PATH {
	if len(ext) == 0 {
		return *p
	}

	offset := len(*p) - 1
	for ; offset >= 0; offset-- {
		if (*p)[offset] == '.' {
			offset++
			break
		}
	}

	if offset == -1 || offset >= len(*p) {
		return *p
	}

	length := offset + len(ext)

	extOffset := 0
	if ext[0] == '.' {
		length--
		extOffset++
	}

	t := make(PATH, length)

	i := 0
	for ; i < offset; i++ {
		t[i] = (*p)[i]
	}

	for ; extOffset < len(ext); extOffset++ {
		t[i] = rune(ext[extOffset])
		i++
	}

	return t
}

func (p *PATH) NameWithoutExtension() string {
	base := p.Name()

	index := strings.LastIndex(base, ".")

	if index == -1 {
		return base
	}

	return base[:index]
}

func (p *PATH) Size() int64 {
	validate(p)

	stat, err := os.Stat(string(*p))

	if err != nil {
		return -1
	}

	return stat.Size()
}

func (p *PATH) FileInfo() os.FileInfo {
	validate(p)

	stat, err := os.Stat(string(*p))

	if err != nil {
		return nil
	}

	return stat
}

func (p *PATH) WalkDirAll() []PATH {
	validate(p)

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
	validate(p)

	list := make([]PATH, 0, 10)

	for _, n := range p.ListStringNames() {
		p := JoinPaths(string(*p), n)

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
	validate(p)

	for _, n := range p.ListStringNames() {
		p := JoinPaths(string(*p), n)

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

//This deletes the path
func (p *PATH) Delete() error {
	validate(p)

	return os.Remove(string(*p))
}

func validate(path *PATH) {
	if len(*path) == 0 {
		panic("PATH length equal zero.")
	}
}

func JoinPaths(paths ...string) PATH {
	return PATH(filepath.Join(paths...))
}
