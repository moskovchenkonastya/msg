package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
//	"go/token"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

// описание используемых структур для базы данных

// Таблица: Хранение пользователей

type authData struct{
	userName     string
	userPassword string
	// passwordHash []byte ?
}

/*

    CREATE TABLE `authData` (
        `uid` INTEGER PRIMARY KEY AUTOINCREMENT,
        `username` VARCHAR(64) NULL,
        `passwordHash` VARCHAR(64) NULL,
        `created` DATE NULL
    );


 */

func handler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
	if err != nil {
		log.Printf("can't print to connection: %s", err)
	}
}

func myJSONHandler(w http.ResponseWriter, r *http.Request) {
	resp := struct {
		Code    string
		Message string
	}{
		Code:    "OK",
		Message: fmt.Sprintf("Hi there, I love %s!", strings.TrimPrefix(r.URL.Path, "/json/")),
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Printf("can't json marshal %+v: %s", resp, err)
		return
	}

	if _, err := w.Write(respBytes); err != nil {
		log.Printf("can't write to connection: %s", err)
	}
}

func connection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./foo.db")

	return db, err
}

func templateHandler(w http.ResponseWriter, r *http.Request) {
	const tmpl = `<h1>{{.Title}}</h1>
  <a href="http://golang.org">GO</a>`
	t, err := template.New("tmpl").Parse(tmpl)
	if err != nil {
		log.Printf("can't create template: %s", err)
		return
	}
	ctx := struct {
		Title string
	}{
		Title: fmt.Sprintf("Hi there, I love %s!", strings.TrimPrefix(r.URL.Path, "/tmpl/")),
	}
	if err := t.Execute(w, ctx); err != nil {
		log.Printf("can't execute and print template: %s", err)
	}
}

func authLogin(w http.ResponseWriter, r *http.Request) {

	resp := struct {
		Code    string
		Message string
		TokenAuth   string
	}{}

	// если авторизация выполнена

	// r.URL.Path

	var auth authData;

	auth.userName = r.FormValue("username")
	auth.userPassword = r.FormValue("password")



	if checkDatabase(auth) == true {
		resp.Code =  "OK"
		resp.Message = "User logon - successful!"
		resp.TokenAuth = "" // func authGenerateToken()
	} else {
		resp.Code =  "ERROR"
		resp.Message = "Logon error!"

	}

	// если авторизация не выполнена

	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Printf("can't json marshal %+v: %s", resp, err)
		return
	}

	if _, err := w.Write(respBytes); err != nil {
		log.Printf("can't write to connection: %s", err)
	}


}

func authRegister(w http.ResponseWriter, r *http.Request) {

	resp := struct {
		Code    string
		Message string
		TokenAuth   string
	}{
		TokenAuth : "",
	}

	// чтение параметров
	//username := r.FormValue("username")
	// password := r.FormValue("password")

	// если пользователь создан
	if dbUserRegister(authData{userName: r.FormValue("username"), userPassword: r.FormValue("password")}) {
		resp.Code =  "OK"
		resp.Message = "User logon - successful!"
		resp.TokenAuth = "28094092380948b3203028409328094"
	} else {
		//если пользователь не создан
		resp.Code =  "ERROR"
		resp.Message = "User cant create!"
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Printf("can't json marshal %+v: %s", resp, err)
		return
	}

	if _, err := w.Write(respBytes); err != nil {
		log.Printf("can't write to connection: %s", err)
	}
}

func authForgive(w http.ResponseWriter, r *http.Request) {

	resp := struct {
		Code    string
		Message string
	}{
	}

	// чтение параметров
	//username := r.FormValue("username")

	// если пользователь есть
	if true {
		resp.Code =  "OK"
		resp.Message = "New password!"
	} else {
		//если пользователь не создан
		resp.Code =  "ERROR"
		resp.Message = "User cant create!"
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		log.Printf("can't json marshal %+v: %s", resp, err)
		return
	}

	if _, err := w.Write(respBytes); err != nil {
		log.Printf("can't write to connection: %s", err)
	}
}

func auth(authdata authData) (bool) {
	// TODO: logAuth
	return checkDatabase(authdata)
}

func checkDatabase(authdata authData) (bool) {
	//connect to database

	db, _ := connection()

	rows, _ := db.Query("SELECT id FROM authData WHERE username=? and passwordHash?", authdata.userName, authdata.userPassword)

	if rows.Next() {
		//consoleLog("Data is correct")
		return true

	}
	rows.Close() //good habit to close
	db.Close()

	return false
}

func dbUserRegister(authdata authData) (bool) {
	//connect to database

	db, err := connection()

	stmt, err := db.Prepare("INSERT INTO authData(username, passwordHash, created) values(?,?,?)")
	res, err := stmt.Exec(authdata.userName, authdata.userPassword, "2017-08-21")

	_, err = res.LastInsertId()

	if err == nil {
		//consoleLog("Data is correct")
		return true

	}
	db.Close()

	return false
}

func dbUserForgive(authdata authData) (bool) {
	//connect to database

	db, err := connection()

	stmt, err := db.Prepare("UPDATE authData SET passwordHash=? WHERE username=?")
	_, err = stmt.Exec(authdata.userPassword, authdata.userName)

	if err == nil {
		//consoleLog("Data is correct")
		return true

	}
	db.Close()

	return false
}






func main() {
	http.HandleFunc("/auth/login", authLogin)
	http.HandleFunc("/auth/register", authRegister)
	http.HandleFunc("/auth/forgive", authForgive)
//	http.HandleFunc("/tmpl/", templateHandler)
//	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
