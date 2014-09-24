package main

import (
	"archive/tar"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"
)

func main() {
	http.HandleFunc("/untar", func(w http.ResponseWriter, r *http.Request) {
		rootDir := r.URL.Query().Get("rootDir")
		if len(rootDir) == 0 {
			rootDir = "/"
		}
		reader := tar.NewReader(r.Body)
		for {
			hdr, err := reader.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalln(err)
			}
			fileInfo := hdr.FileInfo()
			name := path.Join(rootDir, hdr.Name)
			if fileInfo.IsDir() {
				os.MkdirAll(name, fileInfo.Mode())
			} else {
				parentDir := path.Dir(name)
				_, err := os.Stat(parentDir)
				if os.IsNotExist(err) {
					os.MkdirAll(parentDir, fileInfo.Mode())
				}
				output, err := os.Create(name)
				if err != nil {
					log.Fatalln(err)
				}
				io.Copy(output, reader)
				err = output.Chmod(fileInfo.Mode())
				if err != nil {
					log.Printf("Could not set file mode: %v,  file: %v, mode: %v", err, output, fileInfo.Mode())
				}
				output.Close()
				err = os.Chtimes(name, time.Now(), hdr.ModTime)
				if err != nil {
					log.Printf("Could not change file times: %v", err)
				}
			}
		}
	})
	log.Fatal(http.ListenAndServe(":9080", nil))
}
