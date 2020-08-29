package main

import (
	"github.com/YuChaoGithub/meme-linebot/app"
	"github.com/YuChaoGithub/meme-linebot/config"
)

func main() {
	c := config.GetConfig()
	app := &app.App{}
	app.InitializeAndRun(c)
}
