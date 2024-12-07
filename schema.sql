-- Create movies table
CREATE TABLE movies (
    movie_id INTEGER PRIMARY KEY,
    title TEXT NOT NULL,
    release_year INTEGER,
    rating REAL,
    votes INTEGER
);

-- Create movies_genres table
CREATE TABLE movies_genres (
    movie_id INTEGER,
    genre TEXT,
    FOREIGN KEY (movie_id) REFERENCES movies (movie_id)
);
