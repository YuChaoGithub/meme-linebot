package app

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
)

const (
	homeTemplateFilePath = "./ui/html/home.html"
	greetingMemeName     = "bonjour.jpg"
	farewellMemeName     = "adios.jpg"
)

var validSuffixes = []string{".jpg", ".png", ".gif", ".jpeg"}
var punctuations = map[rune]struct{}{
	'?': {}, '!': {}, '.': {}, ',': {}, '\'': {}, ';': {}, ':': {}, '-': {}, '(': {}, ')': {}, '"': {},
	'。': {}, '，': {}, '！': {}, '？': {}, '、': {}, '：': {}, '；': {}, '）': {}, '（': {},
	'‘': {}, '’': {}, '“': {}, '”': {},
}

// homepageHandler renders the home page listing all available memes.
func (a *App) homepageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}

	t, ok := a.pageTemplates["home.html"]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error getting the page template.")
		return
	}

	buf := new(bytes.Buffer)

	memes, err := a.memeModel.GetAll()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error fetching memes from the database.")
		return
	}

	err = t.Execute(buf, memes)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("Error executing the page template.")
	}

	buf.WriteTo(w)
}

// callbackHandler handles line message callbacks.
func (a *App) callbackHandler(w http.ResponseWriter, r *http.Request) {
	// Check if it is a valid line callback request.
	events, err := a.bot.ParseRequest(r)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	// Loop through events and reply to relevant ones.
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			if textMessage, ok := event.Message.(*linebot.TextMessage); ok {
				a.replyWithMeme(event.ReplyToken, textMessage.Text)
			}
		} else if event.Type == linebot.EventTypeMemberJoined {
			a.replyWithMeme(event.ReplyToken, greetingMemeName)
		} else if event.Type == linebot.EventTypeMemberLeft {
			a.replyWithMeme(event.ReplyToken, farewellMemeName)
		}
	}
}

// addMeme is used by the admin to add a meme entry.
func (a *App) addMeme(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	// Retrieve the request body.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Unmarshal json.
	req := struct {
		Admin string `json:"admin"`
		Name  string `json:"name"`
		Link  string `json:"link"`
	}{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the admin key is matched.
	if req.Admin != a.adminSecret {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Insert to the database.
	err = a.memeModel.Insert(req.Name, req.Link)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
	}

	// Success.
	w.WriteHeader(http.StatusCreated)
}

// deleteMeme is used by the admin to delete a meme entry.
func (a *App) deleteMeme(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		return
	}

	// Retrieve the request body.
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Unmarshal json.
	req := struct {
		Admin string `json:"admin"`
		Name  string `json:"name"`
	}{}
	err = json.Unmarshal(body, &req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Check if the admin key is matched.
	if req.Admin != a.adminSecret {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Insert to the database.
	err = a.memeModel.Delete(req.Name)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}

	// Success.
	w.WriteHeader(http.StatusNoContent)
}

// replyWithMeme is a helper function which replies to the event (with the replyToken) with
// a meme named memeName. It does nothing if no such meme exists.
func (a *App) replyWithMeme(replyToken string, memeName string) {
	formattedName := strings.TrimSpace(strings.ToLower(memeName))

	// Check if the format is correct, that is, it has a trailing .jpg, .png, etc.
	isValid := false
	for _, val := range validSuffixes {
		if strings.HasSuffix(formattedName, val) {
			formattedName = strings.TrimSuffix(formattedName, val)
			isValid = true
			break
		}
	}
	if !isValid {
		return
	}

	// Get rid of punctuations.
	cleanedName := []rune{}
	for _, val := range formattedName {
		if _, ok := punctuations[val]; !ok {
			cleanedName = append(cleanedName, val)
		}
	}

	// Get the meme from the database.
	memeURL, err := a.memeModel.Get(string(cleanedName))
	if err != nil {
		// No such meme exists.
		return
	}

	// Reply with the meme image.
	_, err = a.bot.ReplyMessage(replyToken, linebot.NewImageMessage(memeURL, memeURL)).Do()
	if err != nil {
		log.Printf("Error sending reply message with the meme <%v>, link <%v>.\n", memeName, memeURL)
		log.Println(err)
	}
}
