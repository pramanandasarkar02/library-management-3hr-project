package internal

import (
	"encoding/json"
	"github.com/gocql/gocql"
	"log"
	"net/http"
	"strconv"
)

var dataSession *gocql.Session

func StartServer(session *gocql.Session) {
	dataSession = session

	http.HandleFunc("/books", getBooks)
	http.HandleFunc("/trending-book", getTrendingBooks)
	http.HandleFunc("/recent-book", getRecentBooks)
	http.HandleFunc("/search-book", searchBooks)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

/* ---------------- RECENT BOOKS ---------------- */

func getRecentBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		enableCors(w)
		w.WriteHeader(http.StatusOK)
		return
	}

	enableCors(w)
	w.Header().Set("Content-Type", "application/json")

	limit := getLimit(r, 10)
	var books []Book

	query := `
		SELECT isbn, title, subtitle, authors, categories, thumbnail,
		       description, published_year, average_rating, num_pages, ratings_count
		FROM books
		ORDER BY published_year DESC
		LIMIT ? ALLOW FILTERING
	`

	iter := dataSession.Query(query, limit).Iter()
	defer iter.Close()

	for {
		var book Book
		if !iter.Scan(
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
			break
		}
		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)
}



func getTrendingBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		enableCors(w)
		w.WriteHeader(http.StatusOK)
		return
	}

	enableCors(w)
	w.Header().Set("Content-Type", "application/json")

	limit := getLimit(r, 10)
	var books []Book

	log.Println("get trending request")

	query := `
		SELECT isbn, title, subtitle, authors, categories, thumbnail,
		       description, published_year, average_rating, num_pages, ratings_count
		FROM books
		ORDER BY ratings_count DESC
		LIMIT ? ALLOW FILTERING
	`

	iter := dataSession.Query(query, limit).Iter()
	defer iter.Close()

	for {
		var book Book
		if !iter.Scan(
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
			break
		}
		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)
}

/* ---------------- ALL BOOKS ---------------- */

func getBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		enableCors(w)
		w.WriteHeader(http.StatusOK)
		return
	}

	enableCors(w)
	w.Header().Set("Content-Type", "application/json")

	var books []Book

	iter := dataSession.Query(`
		SELECT isbn, title, subtitle, authors, categories, thumbnail,
		       description, published_year, average_rating, num_pages, ratings_count
		FROM books
	`).Iter()
	defer iter.Close()

	for {
		var book Book // FIXED: declared inside loop
		if !iter.Scan(
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
			break
		}
		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)
}

/* ---------------- SEARCH BOOKS ---------------- */

func searchBooks(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		enableCors(w)
		w.WriteHeader(http.StatusOK)
		return
	}

	enableCors(w)
	w.Header().Set("Content-Type", "application/json")

	search := r.URL.Query().Get("query")
	if search == "" {
		http.Error(w, "query parameter is required", http.StatusBadRequest)
		return
	}

	limit := getLimit(r, 10)
	var books []Book

	query := `
		SELECT isbn, title, subtitle, authors, categories, thumbnail,
		       description, published_year, average_rating, num_pages, ratings_count
		FROM books
		WHERE title = ?
		   OR authors CONTAINS ?
		   OR categories CONTAINS ?
		LIMIT ? ALLOW FILTERING
	`

	iter := dataSession.Query(query, search, search, search, limit).Iter()
	defer iter.Close()

	for {
		var book Book
		if !iter.Scan(
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
			break
		}
		books = append(books, book)
	}

	json.NewEncoder(w).Encode(books)
}

/* ---------------- HELPERS ---------------- */

func getLimit(r *http.Request, def int) int {
	if l := r.URL.Query().Get("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 {
			return v
		}
	}
	return def
}
