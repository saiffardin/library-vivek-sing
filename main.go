package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/mux"
)

type Book struct {
	Id, Name, Isbn string
}

const (
	API_PATH = "/apis/v1/books"
)

type library struct {
	dbHost, dbPass, dbName string
}

func main() {
	// DB_HOST is of from host:port
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost:3306"
	}

	dbPass := os.Getenv("DB_PASS")
	if dbPass == "" {
		dbPass = "password"
	}

	apiPath := os.Getenv("API_PATH")
	if apiPath == "" {
		apiPath = API_PATH
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "library"
	}

	l := library{
		dbHost: dbHost,
		dbPass: dbPass,
		dbName: dbName,
	}

	r := mux.NewRouter()
	r.HandleFunc(apiPath, l.getBooks).Methods(http.MethodGet)
	r.HandleFunc(apiPath, l.postBook).Methods(http.MethodPost)

	http.ListenAndServe(":8080", r)

}

func (l library) getBooks(w http.ResponseWriter, r *http.Request) {
	log.Println("getBooks was called")

	// open connection
	db := l.openConnection()
	// read all the books
	rows, err := db.Query("select * from books")
	if err != nil {
		log.Fatalf("querying the books table %s \n", err.Error())
	}

	books := []Book{}
	for rows.Next() {
		var id, name, isbn string
		err := rows.Scan(&id, &name, &isbn)

		if err != nil {
			log.Fatalf("while scanning the row %s \n", err.Error())
		}

		aBook := Book{
			Id:   id,
			Name: name,
			Isbn: isbn,
		}
		books = append(books, aBook)
	}

	json.NewEncoder(w).Encode(books)
	// close connection
	l.closeConnection(db)
}

func (l library) postBook(w http.ResponseWriter, r *http.Request) {
	log.Println("postBooks was called")

	// read the request into an instance of book
	// open connection
	// write all the data
	// close connection


	// read the request into an instance of book
	book := Book{}
	json.NewDecoder(r.Body).Decode(&book)

	// open connection
	db := l.openConnection()

	// write the data
	insertQuery, err := db.Prepare("insert into books values (?, ?, ?)")
	if err != nil {
		log.Fatalf("preparing the db query %s\n", err.Error())
	}
	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("while beginning the transaction %s\n", err.Error())
	}

	_, err = tx.Stmt(insertQuery).Exec(book.Id, book.Name, book.Isbn)
	if err != nil {
		log.Fatalf("execing the insert command %s\n", err.Error())
	}

	err = tx.Commit()
	if err != nil {
		log.Fatalf("while commint the transaction %s\n", err.Error())
	}
	// close the connection
	l.closeConnection(db)
}

func (l library) openConnection() *sql.DB {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s)/%s", "root", l.dbPass, l.dbHost, l.dbName))
	//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]

	if err != nil {
		log.Fatalf("opening the connection to the database %s \n", err.Error())
	}

	return db

}

func (l library) closeConnection(db *sql.DB) {
	err := db.Close()
	if err != nil {
		log.Fatalf("closing connection %s\n", err.Error())
	}
}
