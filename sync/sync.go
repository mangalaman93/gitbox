package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

var client_id = "urwvi2575ynhhu3tr5ki8br5dilsdxyt"
var client_secret = "fSkvmWjP1mNNY8PsJOXZniKkHHMNF2Ra"
var access_token = "w54XOUcU68C1taJtMIPmXAuTSVwCu5Vm"

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

type File_version struct {
	Id string
}

func get_folder_items(folder_id int, limit int) {
	// todo create struct..
	req, err := http.NewRequest("GET", "https://api.box.com/2.0/folders/"+strconv.Itoa(folder_id)+"/items?limit="+strconv.Itoa(limit)+"&offset=0", nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authorization", "Bearer "+access_token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	fmt.Println("body:")
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(body)
	fv := File_version{}
	json.Unmarshal(body, &fv)
	fmt.Println(fv.Id)
	defer resp.Body.Close()
}

func download_file(file_id int, path string) {
	req, err := http.NewRequest("GET", "https://api.box.com/2.0/files/"+strconv.Itoa(file_id)+"/content", nil)
	if err != nil {
		// handle err
	}
	req.Header.Set("Authorization", "Bearer "+access_token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))

	err = ioutil.WriteFile(path, body, 0644)
	defer resp.Body.Close()
}

func main() {
	fmt.Println("test")
	//req, err := http.NewRequest("GET", "https://api.box.com/2.0/folders/0/items?limit=4&offset=0", nil)
	//fmt.Println(err)
	//req.Header.Set("Authorization", "Bearer 4lgSSMxPJMAQMaeJu1W3c74Oi5BPiqR0")
	//client := &http.Client{}
	//resp, _ := client.Do(req)
	//fmt.Println(resp)
	//fmt.Println(resp.Body)

	//get_access_token("HqAHPQYIAaUe14lhAmIx97NDv5PMQqMQ")

	download_file(94149725862, "/tmp/jane.jpeg")

}
