package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	completedPrefix           = "_completed_"
	uploadCooldown            = 1
	imgurUploadEndpoint       = "https://api.imgur.com/3/image"
	memeLineBotUploadEndpoint = "https://meme-linebot.herokuapp.com/add"
)

var validSuffixes = []string{".jpg", ".jpeg", ".png", ".gif"}
var imgurClientID string
var adminSecret string

func main() {
	counter := 0

	if len(os.Args) < 2 {
		fmt.Println("Usage: uploader [directory absolute path]")
		return
	}

	// Get secret environment variables.
	imgurClientID = os.Getenv("IMGUR_CLIENT_ID")
	adminSecret = os.Getenv("ADMIN_SECRET")

	// Get the files from the directory.
	dirPath := os.Args[1]
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Loop through all the files, ignoring subdirectories.
	for _, file := range files {
		if file.IsDir() || file.Name()[0] == '.' || !validExtension(file.Name()) || strings.HasPrefix(file.Name(), completedPrefix) {
			continue
		}

		time.Sleep(uploadCooldown * time.Second)

		// Upload to imgur
		fmt.Printf("Uploading %v to imgur...", file.Name())
		url, err := uploadToImgur(dirPath + "/" + file.Name())
		if err != nil {
			fmt.Printf("failed. Err: %v\n", err)
			continue
		}

		fmt.Println("completed. Now uploading to meme bot database...")

		// Upload to meme bot database.
		splitted := strings.Split(file.Name(), ".")
		err = uploadToMemeDatabase(splitted[0], url)
		if err != nil {
			fmt.Printf("failed. Err: %v\n", err)
			continue
		}

		fmt.Println("completed.")

		// Rename the file if successful.
		os.Rename(dirPath+"/"+file.Name(), dirPath+"/"+completedPrefix+file.Name())

		counter++
	}

	fmt.Println("Successfully uploaded", counter, "images to the meme database.")
}

func validExtension(filename string) bool {
	for _, val := range validSuffixes {
		if strings.HasSuffix(filename, val) {
			return true
		}
	}
	return false
}

func uploadToImgur(filename string) (string, error) {
	// Open file.
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Buffer for multipart form.
	var buf = new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	part, err := writer.CreateFormFile("image", "filename")
	if err != nil {
		return "", err
	}

	// Copy the file to the form.
	io.Copy(part, f)
	writer.Close()

	// Request.
	req, err := http.NewRequest("POST", imgurUploadEndpoint, buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Client-ID "+imgurClientID)

	// Perform post request.
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	respStruct := struct {
		Data struct {
			Link string `json:"link"`
		} `json:"data"`
	}{}

	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		return "", err
	}

	// Only need the imgur ID and extension.
	splitted := strings.Split(respStruct.Data.Link, "/")

	return splitted[len(splitted)-1], nil
}

func uploadToMemeDatabase(name, url string) error {
	reqStruct := struct {
		Admin string `json:"admin"`
		Name  string `json:"name"`
		Link  string `json:"link"`
	}{
		adminSecret,
		name,
		url,
	}

	body, err := json.Marshal(reqStruct)
	if err != nil {
		return err
	}

	// Post request.
	req, err := http.NewRequest("POST", memeLineBotUploadEndpoint, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	// Send request.
	_, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	return nil
}
