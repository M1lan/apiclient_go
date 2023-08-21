package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

var apiBaseURL = "https://reqres.in/api"

// a user is a collection of 5 attributes
type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"avatar"`
}

// helper function to get total amount of users
func getUserCount() int {
	url := apiBaseURL + "/users"
	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject struct {
		Total int `json:"total"`
	}
	err = json.Unmarshal(responseData, &responseObject)
	if err != nil {
		log.Fatal(err)
	}

	return responseObject.Total
}

// write to disk a HTML table generated from template + userlist
func generateHTMLTable(users []User, filename string) error {
	const tableTemplate = `
         <html lang="en-US">
            <head>
               <link href="./style.css" rel="stylesheet">
            </head>
                <table>
                        <tr>
                            <th>ID</th>
                            <th>Email</th>
                            <th>Firstname</th>
                            <th>Lastname</th>
                            <th>Avatar</th>
                        </tr>
                        {{range .}}
                        <tr>
                            <td>{{.ID}}</td>
                            <td>{{.Email}}</td>
                            <td>{{.FirstName}}</td>
                            <td>{{.LastName}}</td>
                            <td><img src="{{.Avatar}}" alt="Avatar or photo"></td>
                        </tr>
                        {{end}}
                </table>
         </html>
        `

	tmpl, err := template.New("table").Parse(tableTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	err = tmpl.Execute(file, users)
	if err != nil {
		return err
	}

	return nil
}

func fetchUsersFromAPI() ([]User, error) {
	// TODO: remove -1 during the demo to get all users!
	userCount := getUserCount() - 1

	url := apiBaseURL + "/users?per_page=" + strconv.Itoa(userCount)
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Expected status code 200 but received %d", response.StatusCode)
	}

	var userList struct {
		Data []User `json:"data"`
	}
	err = json.NewDecoder(response.Body).Decode(&userList)
	if err != nil {
		return nil, err
	}

	return userList.Data, nil
}

func main() {
	users, err := fetchUsersFromAPI()
	if err != nil {
		log.Fatal(err)
	}

	err = generateHTMLTable(users, "./gh_pages/index.html")
	if err != nil {
		log.Fatal(err)
	}
}
