package files

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

// FindWithExtension returns a list of files that contain the given
// extensions in the requested baseDirectory
func FindWithExtension(exts []string, baseDirectory string) []string {
	//	Sanity check that the source directory seems to exist
	if _, err := os.Stat(baseDirectory); os.IsNotExist(err) {
		log.Panic("The directory doesn't exist " + baseDirectory)
	}

	var files []string
	filepath.Walk(baseDirectory, func(path string, f os.FileInfo, _ error) error {
		//	If it's a file...
		if !f.IsDir() {
			//	See if its extension matches one we're looking for...
			if contains(exts, filepath.Ext(f.Name())) {
				//	If it does, Add it to the pile of file results
				files = append(files, path)
			}
		}
		return nil
	})

	//	Return the list of files found
	return files
}

// Copy copies the contents from src to dst using io.Copy.
// If dst does not exist, CopyFile creates it with permissions perm;
// otherwise CopyFile truncates it before writing.
func Copy(src, dst string, perm os.FileMode) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()
	_, err = io.Copy(out, in)
	return
}

// contains returns true if the target slice contains the item 'e'
func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
