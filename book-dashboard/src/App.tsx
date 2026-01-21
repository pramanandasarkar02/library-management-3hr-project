import { useEffect, useState } from "react";
import type { Book } from "./types";

const booksUrl = "http://localhost:8080/books";

function App() {
  const [books, setBooks] = useState<Book[]>([]);
  const [expandedIsbn, setExpandedIsbn] = useState<string | null>(null);
  const [searchText, setSearchText] = useState<string>("");

  useEffect(() => {
    const fetchBooks = async () => {
      const response = await fetch(booksUrl);
      if (!response.ok) {
        throw new Error("Failed to fetch books data");
      }
      const data: Book[] = await response.json();
      setBooks(data);
    };
    fetchBooks();
  }, []);

  const toggleDescription = (isbn: string) => {
    setExpandedIsbn(expandedIsbn === isbn ? null : isbn);
  };


  const handleSearch = () => {
    
  }

  return (
    <div className="min-h-screen bg-gray-100 p-6">
      <h1 className="text-3xl font-bold mb-6 text-center">Books</h1>

      <div className="w-full mb-6 flex gap-2">
        <input
          type="text"
          placeholder="Search by title, author, or category"
          className="w-full border rounded px-4 py-2 outline-none focus:border-blue-500"
          value={searchText}
          onChange={(e) => setSearchText(e.target.value)}
        />
        <button onClick={handleSearch} className="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600">
          Search 
        </button>
      </div>

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
              {Array.isArray(book.authors)
                ? book.authors.join(", ")
                : book.authors}
            </p>

            <p className="text-sm">
              <span className="font-medium">Categories:</span>{" "}
              {Array.isArray(book.categories)
                ? book.categories.join(", ")
                : book.categories}
            </p>

            <p
              onClick={() => toggleDescription(book.isbn)}
              className={`text-sm mt-2 text-gray-700 cursor-pointer ${
                expandedIsbn === book.isbn ? "" : "line-clamp-3"
              }`}
            >
              {book.description}
            </p>

            <p
              onClick={() => toggleDescription(book.isbn)}
              className="text-xs text-blue-600 mt-1 cursor-pointer"
            >
              {expandedIsbn === book.isbn ? "Show less" : "Read more"}
            </p>

            <div className="mt-auto pt-4 text-sm text-gray-900">
              <p>publishedyear: {book.publishedYear}</p>
              <p>averagerating: {book.averageRating}</p>
              <p>ratingCount: {book.ratingCount}</p>
              <p>pages: {book.numPages}</p>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
}

export default App;
