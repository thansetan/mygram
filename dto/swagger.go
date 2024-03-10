package dto

// buat nganu biar pas di swagger munculnya pas ga kurang ga lebih 👍

type UserRegister struct {
	Age      uint64 `json:"age" example:"20"`
	Email    string `json:"email" example:"budiganteng@rocketmail.com"`
	Password string `json:"password" example:"budiganteng123"`
	Username string `json:"username" example:"budiganteng"`
}

type UserLogin struct {
	Email    string `json:"email" example:"budiganteng@rocketmail.com"`
	Password string `json:"password" example:"budiganteng123"`
}

type UserUpdate struct {
	Email    string `json:"email" example:"budigantengbanget@rocketmail.com"`
	Username string `json:"username" example:"budigantengbanget"`
}

type PhotoCreate struct {
	Title   string `json:"title" example:"Gambarnya Budi Ganteng"`
	Caption string `json:"caption" example:"Ini adalah foto Budi yang sangat ganteng"`
	URL     string `json:"photo_url" example:"https://www.budiganteng.com/gambarnya-budi-ganteng.jpg"`
}

type PhotoUpdate struct {
	Title   string `json:"title" example:"Gambarnya Budi Ganteng Banget Sumpah Asli Riil"`
	Caption string `json:"caption" example:"Ini adalah foto Budi yang sangat ganteng banget sumpah asli riil"`
	URL     string `json:"photo_url" example:"https://www.budiganteng.com/ganteng.jpg"`
}

type CommentCreate struct {
	Message string `json:"message" example:"buset ganteng banget nih fotonya"`
}

type CommentUpdate struct {
	Message string `json:"message" example:"buset ganteng banget nih fotonya sumpah asli riil"`
}

type SocialMediaCreate struct {
	Name string `json:"name" example:"Twitter"`
	URL  string `json:"social_media_url" example:"https://twitter.com/budiganteng"`
}

type SocialMediaUpdate struct {
	Name string `json:"name" example:"X"`
	URL  string `json:"social_media_url" example:"https://x.com/budiganteng"`
}
