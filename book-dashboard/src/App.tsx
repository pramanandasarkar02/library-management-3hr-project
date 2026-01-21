import { useEffect, useState } from "react"
import type { Book } from "./types";

const booksUrl = "http://localhost:8080/books";


function App() {
  const [books, setBooks] = useState<Book[]>([])
  
  useEffect(() => {
    const fetchBooks = async () => {
      const response = await fetch(booksUrl);
      if(!response.ok){
        throw new Error("Failed to fetch books data")
      }
      const data: Book[] = await response.json();
      setBooks(data)

    }
    fetchBooks()
  })

  return (
<div className="min-h-screen bg-gray-100 p-6">
  <h1 className="text-3xl font-bold mb-6 text-center">Books</h1>

  <ul className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
    {books.map((book) => (
      <li
        key={book.isbn}
        className="bg-white rounded-lg shadow-md p-4 flex flex-col"
      >
        <img
          src={book.thumbnail}
          alt={book.title}
          className="h-48 w-full object-contain mb-4"
        />

        <h3 className="text-lg font-semibold">{book.title}</h3>
        <p className="text-sm text-gray-500 mb-2">{book.subtitle}</p>

        <p className="text-sm">
          <span className="font-medium">Authors:</span>{" "}
          {book.authors.join(", ")}
        </p>

        <p className="text-sm">
          <span className="font-medium">Categories:</span>{" "}
          {book.categories.join(", ")}
        </p>

        <p className="text-sm mt-2 line-clamp-3 text-gray-700">
          {book.description}
        </p>

        <div className="mt-auto pt-4 text-sm text-gray-600">
          <p>{book.publishedyear}</p>
          <p>{book.averagerating} ({book.ratingCount})</p>
          <p>{book.numPages} pages</p>
        </div>
      </li>
    ))}
  </ul>
</div>
 
  )
}

export default App
