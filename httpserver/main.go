package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func getFiles(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List of files:")
}

func uploadFile(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Uploading file...")

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")

	if err != nil {
		fmt.Println("Error retrieving the file:", err)
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fmt.Println("Uploaded file Name : ", handler.Filename)
	fmt.Println("Uploaded file size: ", handler.Size)

	dst, err := os.Create("../file_storage/" + handler.Filename)

	if err != nil {
		fmt.Println("Error creating file:", err)
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(dst, file)

	if err != nil {
		fmt.Println("Error saving file:", err)
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	fmt.Fprintf(w, "File uploaded successfully")

}

func main() {
	http.HandleFunc("/files", getFiles)
	http.HandleFunc("/upload", uploadFile)

	fmt.Println("Server started on port 8000")

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
		return
	}

}
