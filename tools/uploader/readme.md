# Meme Auto-Uploader for Admin
An auto-uploader for the admin to upload new memes to the Line Meme Bot.

## Usage
```
export IMGUR_CLIENT_ID=<the imgur client id>
export ADMIN_SECRET=<admin secret key>
go run . [absolute path of the directory containing meme images]
```

*Only files with extensions `.jpg`, `.jpeg`, `.gif`, and `.png` will be processed.*

The the filename will be the keyword for the meme.

## Program Structure
1. Loop through all the image files in the directory.
2. Upload the image file to imgur, obtaining the url of the image.
3. Upload the (name, url) pair to the meme line bot server.