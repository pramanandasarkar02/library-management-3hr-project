package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"pramanandasarkar02/library-management/internal"
	"strconv"
	"strings"
	"sync"

	"github.com/gocql/gocql"
)

var BooksDataColumn = []string{
	"isbn10",
	"title",
	"subtitle",
	"authors",
	"categories",
	"thumbnail",
	"description",
	"published_year",
	"average_rating",
	"num_pages",
	"ratings_count",
}
var BookDbColumn = []string{
	"isbn", //isbn10
	"title",
	"subtitle",
	"authors",    
	"categories", 
	"thumbnail",
	"description",
	"published_year",
	"average_rating",
	"num_pages",
	"ratings_count",
}

var DataReadSession *gocql.Session

func insert_csv_data(session *gocql.Session) {
	file, err := os.Open("./data/processed_data.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	isFirstLine := true
	booksCount := 0

	for scanner.Scan() {
		line := scanner.Text()
		if isFirstLine {
			isFirstLine = false
			continue
		}

		values := strings.Split(line, "℧")
		insertQuery := `
		INSERT INTO books (
			isbn,
			title,
			subtitle,
			authors,
			categories,
			thumbnail,
			description,
			published_year,
			average_rating,
			num_pages,
			ratings_count
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

		publishedYear, _ := strconv.Atoi(values[7])
		avgRatingFloat64, _ := strconv.ParseFloat(values[8], 32)
		numPages, _ := strconv.Atoi(values[9])
		ratingCount, _ := strconv.Atoi(values[10])

		avgRating := float32(avgRatingFloat64)

		err = session.Query(insertQuery).Bind(
			values[0],
			values[1],
			values[2],
			values[3],
			values[4],
			values[5],
			values[6],
			publishedYear,
			avgRating,
			numPages,
			ratingCount,
		).Exec()

		if err != nil {
			log.Fatal(err)
		}
		booksCount += 1
	}
	DataReadSession = session
	log.Println("Inserted Books Count: ", booksCount)

}

func init_db() {
	// cluster config
	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Port = 9042
	cluster.Consistency = gocql.Quorum
	

	// create session
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal("Connection err: ", err)
	}
	defer session.Close()

	// remove old keyspace
	err = session.Query(`
		DROP KEYSPACE IF EXISTS library_keyspace
	`).Exec()
	if err != nil {
		log.Fatal("Cannot drop keyspace:", err)
	}

	// create keyspace
	err = session.Query(`
		Create keyspace if not exists library_keyspace
		with replication = {
			'class': 'SimpleStrategy',
			'replication_factor': 1
		}
	`).Exec()

	if err != nil {
		log.Fatal(err)
	}

	session.Close()

	cluster.Keyspace = "library_keyspace"

	session, err = cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	err = session.Query(`
		Create table if not exists books(
			isbn text primary key,
			title text,
			subtitle text,
			authors text,
			categories text,
			thumbnail text,
			description text,
			published_year int,
			average_rating float,
			num_pages int,
			ratings_count int
		)
	`).Exec()

	if err != nil {
		log.Fatal(err)
	}

	insert_csv_data(session)
	session.Close()

	DataReadSession, err = cluster.CreateSession()
	if err != nil{
		log.Fatal(err)
	}

	log.Println("Database Initialize.....")
}



func getBooks(w http.ResponseWriter, r *http.Request){
	enableCors(&w)
	w.Header().Set("Content-Type", "application/json")
	
	var books []internal.Book
	iter := DataReadSession.Query(`
		SELECT 
			isbn,
			title,
			subtitle,
			authors,
			categories,
			thumbnail,
			description,
			published_year,
			average_rating,
			num_pages,
			ratings_count
		FROM books;
	`).Iter()

	var book internal.Book
	for iter.Scan(
		&book.ISBN,
		&book.Title,
		&book.SubTitle,
		&book.Authors,
		&book.Categories,
		&book.Thumbnail,
		&book.Description,
		&book.PublishedYear,
		&book.AverageRating,
		&book.NumPages,
		&book.RatingCount,
	) {
		books = append(books, book)
	}

	if err := iter.Close(); err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}


	json.NewEncoder(w).Encode(books)
}
func enableCors(w *http.ResponseWriter) {
(*w).Header().Set("Access-Control-Allow-Origin", "*")
(*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
(*w).Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
}

func startServer(wg *sync.WaitGroup){
	defer wg.Done()
	http.HandleFunc("/books", getBooks)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

func main() {
	init_db()

	var wg sync.WaitGroup

	wg.Add(1)

	



	go startServer(&wg)

	wg.Wait()
	
}
