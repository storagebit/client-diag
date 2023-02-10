/*The MIT License (MIT)
Copyright © 2020 StorageBIT.ch
Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated
documentation files (the “Software”), to deal in the Software without restriction, including without limitation the
rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit
persons to whom the Software is furnished to do so, subject to the following conditions:
The above copyright notice and this permission notice shall be included in all copies or substantial portions of the
Software.
THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO
THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.IN NO EVENT SHALL THE AUTHORS
OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR
OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package main

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func CreateTarball(tarballFilePath string, filePaths []string) error {

	file, err := os.Create(tarballFilePath)
	if err != nil {
		log.Println(formatRed("Error creating archive: " + err.Error()))
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(formatRed("Error closing archive: " + err.Error()))
		}
	}(file)

	gzipWriter := gzip.NewWriter(file)
	defer func(gzipWriter *gzip.Writer) {
		err := gzipWriter.Close()
		if err != nil {
			log.Println(formatRed("Error closing gzip writer: " + err.Error()))
		}
	}(gzipWriter)

	tarWriter := tar.NewWriter(gzipWriter)
	defer func(tarWriter *tar.Writer) {
		err := tarWriter.Close()
		if err != nil {
			log.Println(formatRed("Error closing tar writer: " + err.Error()))
		}
	}(tarWriter)

	for _, filePath := range filePaths {
		err := addFileToTarWriter(filePath, tarWriter)
		if err != nil {
			log.Println(formatRed("Error adding file to archive: " + err.Error()))
		}
	}

	return nil
}

func addFileToTarWriter(filePath string, tarWriter *tar.Writer) error {
	file, err := os.Open(filePath)
	if err != nil {
		log.Println(formatRed("Error opening file:" + err.Error()))
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(formatRed("Error closing file: " + err.Error()))
		}
	}(file)

	stat, err := file.Stat()
	if err != nil {
		log.Println(formatRed("Error getting file stats: " + err.Error()))
	}

	header := &tar.Header{
		Name:    strings.ReplaceAll(sClientDiagBundleArchiveName, ".tar.gz", "") + "/" + filepath.Base(filePath),
		Size:    stat.Size(),
		Mode:    int64(stat.Mode()),
		ModTime: stat.ModTime(),
	}

	err = tarWriter.WriteHeader(header)
	if err != nil {
		log.Println(formatRed("Error writing header: " + err.Error()))
	}

	_, err = io.Copy(tarWriter, file)
	if err != nil {
		log.Println(formatRed("Error copying file data into tarball: " + err.Error()))
	}

	return nil
}
