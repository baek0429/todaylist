package model

type Post struct {
	Id          int
	Title       string
	Password    string
	Address     string
	Hour        string
	ImagePath   string
	Description string
}

func SamplePosts() []Post {
	posts := make([]Post, 10)
	for i, post := range posts {
		post.Id = i
		post.Title = "title"
		post.Address = "address"
		post.Hour = "hour"
		post.ImagePath = "img/Mountain_Bluebird.jpg"
		post.Description = "description"
		post.Password = "1"
		posts[i] = post
	}
	return posts
}

func GetAddEmptyPost() []Post {
	return []Post{Post{Title: "Type your title",
		Address:     "type your address",
		Hour:        "xx:xx~xx:xx",
		Description: "Describe your service",
		Password:    "Password"}}
}

func GetMainPosts() []Post {
	return SamplePosts()
}

func GetPostByCategory() []Post {
	return SamplePosts()
}
