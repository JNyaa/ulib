package seven

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/bodgit/sevenzip"

	"github.com/3JoB/ulib/fsutil"
	ph "github.com/3JoB/ulib/path"
)

type sevenZip struct {
	pass string
}

func New(pwd ...string) *sevenZip {
	if len(pwd) != 0 {
		return &sevenZip{pass: pwd[0]}
	}
	return &sevenZip{}
}

func (sv sevenZip) Extract(source, destination string) (extractedFiles []string, err error) {
	var i *sevenzip.ReadCloser
	if sv.pass != "" {
		i, err = sevenzip.OpenReaderWithPassword(source, sv.pass)
	} else {
		i, err = sevenzip.OpenReader(source)
	}
	if err != nil {
		return nil, err
	}
	defer i.Close()
	if err := fsutil.Mkdir(destination, 0755); err != nil {
		return nil, err
	}

	for _, f := range i.File {
		if err := sv.extractAndWriteFile(destination, f); err != nil {
			return nil, err
		}

		extractedFiles = append(extractedFiles, f.Name)
	}

	return extractedFiles, nil
}

func (sevenZip) extractAndWriteFile(destination string, f *sevenzip.File) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	path := ph.Join(destination, f.Name)
	if !strings.HasPrefix(path, ph.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("%s: illegal file path", path)
	}

	if f.FileInfo().IsDir() {
		if err = fsutil.Mkdir(path, f.Mode()); err != nil {
			return err
		}
	} else {
		if err := fsutil.Mkdir(ph.DirPath(path), f.Mode()); err != nil {
			return err
		}
		if f, err := fsutil.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode()); err != nil {
			return err
		} else {
			defer f.Close()
			if _, err := io.Copy(f, rc); err != nil {
				return err
			}
		}
	}
	return nil
}
