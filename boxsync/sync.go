package boxsync

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"errors"
	"os"
	"bytes"
	"mime/multipart"
	"path/filepath"
	"io"
	"fmt"
	"os/exec"
)

var Client_Id = "urwvi2575ynhhu3tr5ki8br5dilsdxyt"
var Client_Secret = "fSkvmWjP1mNNY8PsJOXZniKkHHMNF2Ra"
var Access_Token = "jITEhGRQWgCGHxpwPudlaa8ZhiSh74h5"

//func get_access_token(auth_code string) {
//	// do manually for now...
//	body := strings.NewReader("grant_type=authorization_code&code=" + auth_code + "&client_id=" + client_id + "&client_secret=" + client_secret)
//	req, err := http.NewRequest("POST", "https://api.box.com/oauth2/token", body)
//	if err != nil {
//		// handle err
//	}
//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//
//	resp, err := http.DefaultClient.Do(req)
//	if err != nil {
//		// handle err
//	}
//
//	fmt.Println(resp)
//
//	defer resp.Body.Close()
//}

type File_Version struct {
	Id string
}

type Entry struct {
	Type         string
	Id           string
	File_version File_Version
	Sequence_Id  string
	Etag         string
	Name         string
}

type Folder_Items struct {
	Total_Count int
	Entries     []Entry
	Offset      int
	Limit       int
}

func Get_Folder_Items(folder_id string) []Entry {
	// todo create struct..
	req, err := http.NewRequest("GET", "https://api.box.com/2.0/folders/" + folder_id + "/items?limit=100&offset=0", nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authorization", "Bearer " + Access_Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fi := Folder_Items{}
	json.Unmarshal(body, &fi)
	defer resp.Body.Close()
	return fi.Entries
}

func Download_File_By_Id(file_id string, path string) error {
	req, err := http.NewRequest("GET", "https://api.box.com/2.0/files/" + file_id + "/content", nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer " + Access_Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	body, _ := ioutil.ReadAll(resp.Body)

	err = ioutil.WriteFile(path, body, 0644)
	defer resp.Body.Close()
	return err
}

func Download_File_By_Path(boxpath string, path string) error {
	// boxpath: A/B/C/a.txt
	// path: X/Y/Z/e.txt
	id, Type, err := Get_Object_Id(boxpath)
	if (err != nil) {
		return err
	}
	if (Type != "file") {
		return errors.New("Path doesn't point to a file!")
	}
	return Download_File_By_Id(id, path)
}

func Download_Folder(boxpath string, localpath string) error {
	// boxpath: A/B/C (strict!)
	// localpath: X/Y/Z/ (strict!)
	// Creates: X/Y/Z/C/...

	p := Remove_Slashes_At_Ends(boxpath)
	dirs := strings.Split(p, "/")

	id, ty, err := Get_Object_Id(boxpath)
	if err != nil {
		return err
	}
	if ty != "folder" {
		return errors.New("Path is not to a folder")
	}
	curr_name := dirs[len(dirs) - 1]
	os.MkdirAll(localpath + "/" + curr_name, 0755)
	curr_contents := Get_Folder_Items(id)
	for _, entry := range curr_contents {
		if entry.Type == "file" {
			err = Download_File_By_Path(boxpath + "/" + entry.Name, localpath + curr_name + "/" + entry.Name)
			if err != nil {
				return err
			}
		} else if entry.Type == "folder" {
			err = Download_Folder(boxpath + "/" + entry.Name, localpath + curr_name + "/")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Remove_Slashes_At_Ends(path string) string {
	ret := path
	if (string(ret[0]) == "/") {
		ret = ret[1:len(ret)]
	}
	if (string(ret[len(ret) - 1]) == "/") {
		ret = ret[0:len(ret) - 1]
	}
	return ret
}

func Get_Object_Id(boxpath string) (string, string, error) {
	// returns id, type, err
	// Object can be file or folder

	path := Remove_Slashes_At_Ends(boxpath)
	dirs := strings.Split(path, "/")
	curr_contents := Get_Folder_Items("0")
	var curr_id string
	found := 0
	for i, dir := range dirs {
		if (i == len(dirs) - 1) {
			break
		}
		for _, entry := range curr_contents {
			if (strings.Compare(entry.Name, dir) == 0) {
				curr_id = entry.Id
				found = 1
			}
		}
		if (found == 0) {
			return "", "", errors.New("Invalid path")
		}
		curr_contents = Get_Folder_Items(curr_id)
	}
	found = 0
	for _, entry := range curr_contents {
		if (strings.Compare(entry.Name, dirs[len(dirs) - 1]) == 0) {
			found = 1
			return entry.Id, entry.Type, nil
		}
	}
	return "", "", errors.New("Invalid path")
}

//func Upload_File_By_Id(file_id string, path string) error {
//	data, err := os.Open(path)
//	if err != nil {
//		return err
//	}
//	defer data.Close()
//	req, err := http.NewRequest("PUT", "https://upload.box.com/api/2.0/files/" + file_id + "/content", data)
//	req.Header.Set("Authorization", "Bearer e7E2PvIYVSjfTvdk6akkyFO1KBE9Wv4d")
//	req.Header.Set("Content-Type", "text/plain")
//	resp, err := http.DefaultClient.Do(req)
//	body, _ := ioutil.ReadAll(resp.Body)
//	fmt.Println(string(body))
//	return nil
//}

func Upload_File_By_Id(file_id string, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("filename", filepath.Base(path))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)

	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://upload.box.com/api/2.0/files/" + file_id + "/content", body)
	req.Header.Set("Authorization", "Bearer e7E2PvIYVSjfTvdk6akkyFO1KBE9Wv4d")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	_, err = http.DefaultClient.Do(req)
	//_, _ := ioutil.ReadAll(resp.Body)
	////fmt.Println(string(body2))
	return err
}

//func Upload_New_File(parent_id string, path string) error {
//
//	p := Remove_Slashes_At_Ends(path)
//	dirs := strings.Split(p, "/")
//	name := dirs[len(dirs) - 1]
//	parent := map[string]string{
//		"id":        parent_id,
//	}
//
//	file, err := os.Open(path)
//	if err != nil {
//		return err
//	}
//	defer file.Close()
//
//	body := &bytes.Buffer{}
//
//	writer := multipart.NewWriter(body)
//	part, err := writer.CreateFormFile("filename", filepath.Base(path))
//	if err != nil {
//		return err
//	}
//	_, err = io.Copy(part, file)
//
//	_ = writer.WriteField("name", name)
//	w1, _ := writer.CreateFormField("parent")
//
//	w1.Write(parent)
//
//	err = writer.Close()
//	if err != nil {
//		return err
//	}
//
//	req, err := http.NewRequest("POST", "https://upload.box.com/api/2.0/files/content", body)
//	req.Header.Set("Authorization", "Bearer e7E2PvIYVSjfTvdk6akkyFO1KBE9Wv4d")
//
//	req.Header.Set("Content-Type", writer.FormDataContentType())
//	_, err = http.DefaultClient.Do(req)
//	//_, _ := ioutil.ReadAll(resp.Body)
//	////fmt.Println(string(body2))
//
//	return err
//}

func Upload_New_File(parent_id string, path string) error {

	p := Remove_Slashes_At_Ends(path)
	dirs := strings.Split(p, "/")
	name := dirs[len(dirs) - 1]
	jsonMap := map[string]string{
		"id": parent_id,
	}
	json_str, _ := json.Marshal(jsonMap)

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)
	_ = writer.WriteField("name", name)
	_ = writer.WriteField("parent", string(json_str))

	part, err := writer.CreateFormFile("filename", filepath.Base(path))
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)

	//w1, _ := writer.CreateFormField("parent")

	//w1.Write(json_str)

	err = writer.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://upload.box.com/api/2.0/files/content", body)
	req.Header.Set("Authorization", "Bearer " + Access_Token)

	req.Header.Set("Content-Type", writer.FormDataContentType())
	////b, _ := ioutil.ReadAll(req.Body)
	////fmt.Println(string(b))
	////
	//resp, err := http.DefaultClient.Do(req)
	//body2, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(body2))
	cmd_arr := []string{"https://upload.box.com/api/2.0/files/content", "-H", "\"Authorization: Bearer " + Access_Token + "\"", "-X", "POST", "-F", "attributes='{\"name\":\"" + name + "\", \"parent\":{\"id\":\"" + parent_id + "\"}}'", "-F", "file=@" + path}
	fmt.Println(cmd_arr)
	c := exec.Command("curl", cmd_arr...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err = c.Run()
	if err != nil {
		fmt.Println("Error: ", err)
	}

	//fmt.Println(string(b1))
	return err
}
