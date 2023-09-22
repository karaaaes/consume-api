package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

type Avatar struct {
	ID             int64  `json:"id"`
	AvatarName     string `json:"avatar_name"`
	AvatarImage    string `json:"avatar_image"`
	AvatarUsername string `json:"avatar_username"`
	AvatarPassword string `json:"avatar_password"`
	AvatarEmail    string `json:"avatar_email"`
}

type Response struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	URL    string `json:"url"`
}

type AvatarResponseAll struct {
	Data     []Avatar `json:"data"`
	Response Response `json:"response"`
}

type AvatarResponse struct {
	Data     Avatar   `json:"data"`
	Response Response `json:"response"`
}

var BASE_URL = "http://localhost:8000/api"

func Index(w http.ResponseWriter, r *http.Request) {
	// Buat struct untuk respons JSON
	var avatarsResponseAll AvatarResponseAll
	response, err := http.Get(BASE_URL + "/avatar")
	if err != nil {
		log.Print(err)
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		http.Error(w, "Failed to fetch data", http.StatusInternalServerError)
		return
	}

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&avatarsResponseAll); err != nil {
		log.Print(err)
		http.Error(w, "Failed to decode response", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"avatars": avatarsResponseAll.Data, // Gunakan avatarsResponse.Data sebagai data avatar
	}

	temp, err := template.ParseFiles("view/index.html")
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := temp.Execute(w, data); err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func Add(w http.ResponseWriter, r *http.Request) {
	temp, _ := template.ParseFiles("view/add.html")
	temp.Execute(w, nil)
}

func Store(w http.ResponseWriter, r *http.Request) {
	maxID := getMaxId()
	r.ParseMultipartForm(10 << 20) // Mengizinkan upload hingga 10 MB (sesuaikan dengan kebutuhan Anda)

	// Mengambil file dari form
	file, _, err := r.FormFile("avatar_image")
	if err != nil {
		// Handle error
		log.Println(err)
		return
	}
	defer file.Close()

	// Mengambil nama file dari header
	fileHeader := r.MultipartForm.File["avatar_image"][0]
	avatarImageName := fileHeader.Filename

	fmt.Println("Nama File : " + avatarImageName)

	newPost := Avatar{
		ID:             maxID + 1,
		AvatarName:     r.FormValue("avatar_name"),
		AvatarImage:    "images/" + avatarImageName, // Gunakan nama file asli sebagai AvatarImage
		AvatarUsername: r.FormValue("avatar_username"),
		AvatarPassword: r.FormValue("avatar_password"),
		AvatarEmail:    r.FormValue("avatar_email"),
	}

	jsonValue, _ := json.Marshal(newPost)
	buffer := bytes.NewBuffer(jsonValue)
	req, err := http.NewRequest(http.MethodPost, BASE_URL+"/avatar", buffer)
	if err != nil {
		log.Print(err)
	}
	req.Header.Set("Content-Type", "application/json; charset = UTF-8")

	httpClient := &http.Client{}
	res, err := httpClient.Do(req)

	if err != nil {
		log.Print(err)
	}
	defer res.Body.Close()

	var avatarResponse AvatarResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&avatarResponse); err != nil {
		log.Print(err)
	}

	if avatarResponse.Response.Code == 200 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func Update(w http.ResponseWriter, r *http.Request) {
	// Mendapatkan ID dari query parameter
	id := r.URL.Query().Get("id")
	concat := BASE_URL + "/avatar/" + id

	// Mengirim permintaan GET ke endpoint Avatar dengan ID tertentu
	res, err := http.Get(concat)
	if err != nil {
		log.Print(err)
		http.Error(w, "Error 1 - Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	// Membaca respons JSON ke dalam struct Avatar
	var avatar AvatarResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&avatar); err != nil {
		log.Print(err)
		http.Error(w, "Error 2 - Internal Server Error", http.StatusInternalServerError)
		return
	}

	fmt.Println(avatar)

	// Menyiapkan data yang akan dikirimkan ke template
	data := map[string]interface{}{
		"avatar": avatar.Data,
	}

	if avatar.Response.Code == 200 {
		// Parsing dan mengeksekusi template
		temp, err := template.ParseFiles("view/edit.html")
		if err != nil {
			log.Print(err)
			http.Error(w, "Error 3 - Internal Server Error", http.StatusInternalServerError)
			return
		}

		if err := temp.Execute(w, data); err != nil {
			log.Print(err)
			http.Error(w, "Error 4 - Internal Server Error", http.StatusInternalServerError)
			return
		}
	}
}

func ExecuteUpdate(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // Mengizinkan upload hingga 10 MB (sesuaikan dengan kebutuhan Anda)

	// Mengambil nama file dari header
	fileHeader := r.MultipartForm.File["avatar_image"]
	var avatarImageName string

	// Jika ada file yang diunggah, gunakan nama file baru
	if len(fileHeader) > 0 {
		avatarImageName = "images/" + fileHeader[0].Filename
	} else {
		// Jika tidak ada file yang diunggah, gunakan nama file yang ada (avatar_image_old)
		avatarImageName = r.FormValue("avatar_image_old")
	}

	newPost := Avatar{
		AvatarName:     r.FormValue("avatar_name"),
		AvatarImage:    avatarImageName, // Gunakan nama file asli atau yang lama sebagai AvatarImage
		AvatarUsername: r.FormValue("avatar_username"),
		AvatarPassword: r.FormValue("avatar_password"),
		AvatarEmail:    r.FormValue("avatar_email"),
	}

	jsonValue, _ := json.Marshal(newPost)
	buffer := bytes.NewBuffer(jsonValue)
	req, err := http.NewRequest(http.MethodPut, BASE_URL+"/avatar/"+r.FormValue("avatar_id"), buffer)
	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")

	httpClient := &http.Client{}
	res, err := httpClient.Do(req)

	if err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	var avatarResponse AvatarResponse
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&avatarResponse); err != nil {
		log.Print(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if avatarResponse.Response.Code == 200 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	concat := BASE_URL + "/avatar/" + id

	fmt.Println(concat)
	req, err := http.NewRequest(http.MethodDelete, concat, nil)
	if err != nil {
		log.Print(err)
	}

	httpClient := &http.Client{}
	res, err := httpClient.Do(req)
	if err != nil {
		log.Print(err)
	}

	defer res.Body.Close()

	fmt.Println(res.StatusCode)
	if res.StatusCode == 200 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}

}

func getMaxId() int64 {
	// Buat struct untuk respons JSON
	var avatarsResponse struct {
		Data []Avatar `json:"data"`
	}

	response, err := http.Get(BASE_URL + "/avatar")
	if err != nil {
		log.Print(err)
		return 0 // Mengembalikan 0 jika terjadi kesalahan
	}

	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&avatarsResponse); err != nil {
		log.Print(err)
		return 0 // Mengembalikan 0 jika terjadi kesalahan
	}

	// Mencari ID paling besar
	var maxID int64
	for _, avatar := range avatarsResponse.Data {
		if avatar.ID > maxID {
			maxID = avatar.ID
		}
	}
	return maxID
}
