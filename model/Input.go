package model

type Input struct {
	Title       string
	MainDisplay string
	Posts       []Post
}

func GetMainInput() Input {
	return Input{Title: "Todaylist", MainDisplay: "maindisplay", Posts: GetMainPosts()}
}

func GetAddEmptyInput() Input {
	return Input{Title: "Todaylist", Posts: GetAddEmptyPost()}
}
