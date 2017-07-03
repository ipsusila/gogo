package http

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"
)

//DON'T forget to add PORT to firewall exception
var (
	fileField    = "files"
	dataDir      = "testdata/"
	uploadDir    = "testdata/upload/"
	uploadTarget = "/upload"
	serverPort   = ":8080"
	serverHost   = "http://localhost"
)

//File upload handler for testing.
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	const maxMemory = 1024 * 1024
	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//Our multipart form
	multi := r.MultipartForm
	log.Printf("File upload request received")

	//handles form data
	log.Printf("Process custom field(s)...")
	for field, value := range multi.Value {
		fmt.Printf("  Field[`%s`]: %v\n", field, value)
	}

	//Get file header.
	log.Printf("Process file(s)...")
	files := multi.File[fileField]
	for _, fh := range files {
		fmt.Printf("  %s\n", fh.Filename)
		//Handle each files
		file, err := fh.Open()
		defer file.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//Create destination file
		dstFile, err := os.Create(uploadDir + fh.Filename)
		defer dstFile.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		//Copy and overwrite destination file
		if _, err := io.Copy(dstFile, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	str := fmt.Sprintf("fields=%02d, files=%02d\n", len(multi.Value), len(multi.File))
	w.Write([]byte(str))
}

func TestMain(m *testing.M) {
	server := setup()
	go server.ListenAndServe()

	code := m.Run()
	shutdown(server)
	os.Exit(code)
}

func setup() *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc(uploadTarget, uploadHandler)

	//DON'T forget to allow connection to PORT
	//when the PC has firewall.
	server := &http.Server{
		Addr:    serverPort,
		Handler: mux,
	}

	return server
}

func shutdown(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Error shutdown HTTP server: %v", err)
	}
	log.Printf("--- DONE ---")
}

func doSubmit(submit func(string) (*http.Response, error), failOnError bool, t *testing.T) {
	url := serverHost + serverPort + uploadTarget
	resp, err := submit(url)
	if err != nil {
		if failOnError {
			t.Fatal(err)
		} else {
			t.Logf("Error: %s", err)
		}
		return
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		if failOnError {
			t.Fatal(err)
		} else {
			t.Logf("Error: %s", err)
		}
	}

	log.Printf("Response=%v", string(response))
}

func TestFieldsOnly(t *testing.T) {
	fu := NewFormUploader()
	fu.AddField("id", "Field-only-upload")
	fu.AddField("id", "Second ID")
	fu.AddField("time", time.Now().Format(time.RFC3339))
	fu.AddField("description", "Custom information")

	//perform upload
	doSubmit(fu.Post, true, t)
}

func TestSingleFileOnly(t *testing.T) {
	fu := NewFormUploader()
	fu.AddFiles(fileField, dataDir+"image01.jpg")

	//perform upload
	doSubmit(fu.Post, true, t)
}

func TestMultipleFilesOnly(t *testing.T) {
	fu := NewFormUploader()
	files := []string{
		dataDir + "image01.jpg",
		dataDir + "file01.txt",
		dataDir + "file02.pdf",
		dataDir + "file03.pdf",
	}
	fu.AddFiles(fileField, files...)

	//perform upload
	doSubmit(fu.Post, true, t)
}

func TestFilesWithFields(t *testing.T) {
	fu := NewFormUploader()

	//Add fields
	fu.AddField("id", "File and custom files")
	fu.AddField("time", time.Now().Format(time.RFC3339))
	fu.AddField("description", "Custom information")

	files := []string{
		dataDir + "image01.jpg",
		dataDir + "file01.txt",
		dataDir + "file02.pdf",
	}
	fu.AddFiles(fileField, files...)

	//perform upload
	doSubmit(fu.Post, true, t)
}

func TestPutFilesWithFields(t *testing.T) {
	fu := NewFormUploader()

	//Add fields
	fu.AddField("id", "File and custom files")
	fu.AddField("time", time.Now().Format(time.RFC3339))
	fu.AddField("description", "Custom information")

	files := []string{
		dataDir + "image01.jpg",
		dataDir + "file01.txt",
		dataDir + "conflict/file01.txt",
		dataDir + "file02.pdf",
	}
	fu.AddFiles(fileField, files...)

	//perform upload: PUT
	doSubmit(fu.Put, true, t)
}

func TestFileDoesNotExist(t *testing.T) {
	fu := NewFormUploader()
	files := []string{
		dataDir + "file01.txt",
		"desnotexist.txt",
	}
	fu.AddFiles(fileField, files...)

	//perform upload
	doSubmit(fu.Post, false, t)
}
