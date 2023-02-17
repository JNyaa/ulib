package compress

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zip"
	zs "github.com/klauspost/compress/zstd"
)

type Zip struct {}

func NewZip() *Zip {
	return &Zip{}
}

func (uz Zip) Extract(source, destination string) ([]string, error) {
	decomp := zs.ZipDecompressor()
	zip.RegisterDecompressor(zs.ZipMethodPKWare, decomp)
	zip.RegisterDecompressor(zs.ZipMethodWinZip, decomp)
	r, err := zip.OpenReader(source)
	if err != nil {
		return nil, err
	}

	defer r.Close()

	if err := os.MkdirAll(destination, 0755); err != nil {
		return nil, err
	}

	var extractedFiles []string
	for _, f := range r.File {
		if err := uz.extractAndWriteFile(destination, f); err != nil {
			return nil, err
		}

		extractedFiles = append(extractedFiles, f.Name)
	}

	return extractedFiles, nil
}

func (Zip) extractAndWriteFile(destination string, f *zip.File) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	path := filepath.Join(destination, f.Name)
	if !strings.HasPrefix(path, filepath.Clean(destination)+string(os.PathSeparator)) {
		return fmt.Errorf("%s: illegal file path", path)
	}

	if f.FileInfo().IsDir() {
		if err = os.MkdirAll(path, f.Mode()); err != nil {
			return err
		}
	} else {
		err = os.MkdirAll(filepath.Dir(path), f.Mode())
		if err != nil {
			return err
		}

		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer f.Close()

		if _, err := io.Copy(f, rc); err != nil {
			return err
		}
	}

	return nil
}