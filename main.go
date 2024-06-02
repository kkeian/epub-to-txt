package main

import (
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"io/fs"
	"log"
	"path/filepath"
	"regexp"
)

const dataDir = "data/"
const targetFileExt = ".epub"

func main() {
	// 1. Read a directory's contents,
	//    storing each matching file
	var epubs []string
	err := filepath.WalkDir(dataDir, func(path string, entry fs.DirEntry, e error) error {
		if e != nil { // error when attempting to read the directory entry
			return e
		}
		// match on file extension
		filename := entry.Name()
		found, err := filepath.Match("*"+targetFileExt, filename)
		if err != nil {
			panic(err)
		}
		if found {
			// fmt.Println(entry.Name())
			epubs = append(epubs, path)
		}
		return nil
	})
	if err != nil {
		panic(err) // if we can't walk the directory there's no point in continuing
	}

	// Process each file
	for _, epub := range epubs {
		processArchive(epub)
	}
}

func processArchive(file string) {
	fmt.Println(file)
	// Unzip/read .epub file's contents
	handle, err := zip.OpenReader(file)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close() // call this at end of enclosing func
	textFiles := filterForTextfiles(handle.File)
	// bookTitle := getTitle(textFiles[0])
	_ = getTitle(textFiles[0])
}

func filterForTextfiles(ar []*zip.File) []*zip.File {
	textFileExt := regexp.MustCompile(`\.html`)
	var htmlFiles []*zip.File
	for _, f := range ar {
		if textFileExt.Match([]byte(f.Name)) {
			htmlFiles = append(htmlFiles, f)
		}
	}
	return htmlFiles
}

func getTitle(f *zip.File) string {
	type Doc struct {
		XMLName xml.Name `xml:"html"`
		Title   []byte   `xml:"head>title"`
	}

	contentHandle, err := f.Open()
	if err != nil {
		log.Printf("Unable to open %s.\n", f.Name)
	}
	defer contentHandle.Close()

	byteChunkSize := 1024
	title := make([]byte, byteChunkSize)
	_, err = io.ReadAtLeast(contentHandle, title, byteChunkSize)
	if err != nil && err != io.ErrUnexpectedEOF {
		log.Fatalf("Error reading file %v", err)
	}

	// fmt.Printf("Bytes read : %d\n", bytesRead)
	// fmt.Printf("%s", string(title))

	var h Doc
	err = xml.Unmarshal(title, &h)
	if err != nil {
		log.Fatalf("Error unmarshalling %s %v", f.Name, err)
	}

	fmt.Printf("%v\n", string(h.Title))

	return ""
}

// func printFileStats(file *zip.File) []*zip.File {
// 	// fmt.Printf("\tUncompressed size : %d\n", file.UncompressedSize64)
// 	contentHandle, err := file.Open()
// 	if err != nil {
// 		log.Printf("Unable to open %s.\n", file.Name)
// 	}
// 	defer contentHandle.Close()

// 	// bytesRead, err := io.Copy(os.Stdout, contentHandle)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// fmt.Printf("\tBytes Read : %d\n", bytesRead)
// }
