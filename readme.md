# Line Chatbot: Meme Collector (Meme藏家)
A Line chatbot which replies a meme image when certain keywords are detected.

![Add me on Line](https://i.imgur.com/PZ4xgtx.png)

# Usage
1. Add me on Line using QRCode or [this link](https://line.me/ti/p/@560xwtfv).
2. Send a meme keyword with `.jpg`. (Try `我就爛.jpg`, or `always has been.jpg`, check the full keyword list on [here](https://meme-linebot.herokuapp.com/).)
3. Nice!

# Development
It is difficult to run this chatbot in local environments because the chatbot is triggered via webhooks.

## Run Tests
```
./start_test_container.sh
go test ./...
./clean_test_container.sh
```

# Admin APIs
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
* Write an administrative web frontend for the author so that he can upload new memes easier.
* Redesign the frontend of the homepage.
* Write more unit tests.

# Development Log
## 2020.08.29
* Started the project this morning and finished it before dinner. Yup.