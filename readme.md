# Line Chatbot: Meme Collector (Meme藏家)
A Line chatbot which replies with a meme image when certain keywords are detected.

Add me on Line:

ID: @560xwtfv

<img src="https://i.imgur.com/PZ4xgtx.png" alt="QRCode" width=250>

[Demo Video](https://youtu.be/mHjv9NcskbA)

Usage Example:

<img src="https://i.imgur.com/VLL3J2w.jpg" alt="Demo" height=500>

# Usage
1. Add me on Line using QRCode or [this link](https://line.me/ti/p/@560xwtfv).
2. Send a meme keyword with `.jpg`. (Try `我就爛.jpg`, or `always has been.jpg`.)
3. [Nice!](https://i.imgur.com/mUUOa0v.jpg)

[Full Command List (continuously being updating)](https://meme-linebot.herokuapp.com/)

# Development & Deployment
It is difficult to run this chatbot in local environments because the chatbot is triggered via webhooks.

## Run Tests
```
./start_test_container.sh
go test ./...
./clean_test_container.sh
```

## Deploy to Heroku
```
heroku container:push -a meme-linebot web
heroku container:release -a meme-linebot web
```

# Admin APIs
Use the **uploader** tool in `./tools/uploader` to automatically upload meme images from a local directory.

## `/add`
Add a new meme entry.

Request Body:

```
{
    "admin": "the administrator secret.",
    "name": "memeName",
    "link": "imgurID"
}
```

## `/delete`
Delete an existing meme entry.

Request Body:

```
{
    "admin": "the administrator secret.",
    "name": "memeName"
}
```

# Future Plan
* Write more unit tests. Only `package models` is fully tested now.
* Redesign the frontend of the homepage.
* Implement a tracker to track the statistics of meme usage.

# Development Log
*(In reverse chronological order.)*

## 2020.09.14
* Fuzzy search supported. (i.e. Even if no meme matches the user's message, the meme with the closest name would be returned.)

## 2020.09.07
* Now the the Line ID is shown in the homepage.

## 2020.08.30
* The uploader tool will now show the error when it fails to upload image to Imgur.
* Now punctuations will be stripped out of the received messages, so that messages with redundant punctuations can also trigger the chatbot.

## 2020.08.29
* Started the project this morning and finished it before dinner. Yup.
