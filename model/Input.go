package model

import "net/url"

type Input struct {
	Title         string
	MainDisplay   string
	Posts         []Post
	BlobActionURL *url.URL
}

func GetMainInput() Input {
	return Input{Title: "Todaylist", MainDisplay: "maindisplay", Posts: GetMainPosts()}
}

func GetAddEmptyInput() Input {
	return Input{Title: "Todaylist"}
}
