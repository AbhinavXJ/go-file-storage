package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"

	"strconv"

	"github.com/joho/godotenv"
)

type Response struct {
	Url string `json:"url"`
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

	files, err := os.ReadDir("../file_storage")
	if err != nil {
		fmt.Println("Error reading directory:", err)
	}

	fileExists := false
	for _, file := range files {
		if file.Name() == handler.Filename {
			fileExists = true
			break
		}
	}

	newFileName := handler.Filename
	randNum := rand.IntN(100)
	if fileExists {
		newFileName = strconv.Itoa(randNum) + "_" + newFileName
	}

	dst, err := os.Create("../file_storage/" + newFileName)

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
	domain := os.Getenv("DOMAIN")
	if domain == "" {
		domain = "http://localhost:8000"
	}
	res_url := domain + "/file_storage/" + newFileName

	w.Header().Set("Content-Type", "application/json")
	res := Response{
		Url: res_url,
	}
	json.NewEncoder(w).Encode(res)

}

func main() {
	err1 := godotenv.Load("../.env")
	if err1 != nil {
		fmt.Println("Error loading .env file")
	}
	fs := http.FileServer(http.Dir("../file_storage"))
	http.Handle("/file_storage/", http.StripPrefix("/file_storage/", fs))
	http.HandleFunc("/upload", uploadFile)

	fmt.Println("Server started on port 8000")

	port := os.Getenv("PORT")

	if port == "" {
		port = "8000"
	}

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
		return
	}

}
