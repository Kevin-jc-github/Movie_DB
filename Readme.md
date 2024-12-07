# README: Personal Movie Database with SQLite and Go

## Overview
This project demonstrates how to build a personal movie database using Go and SQLite. The database stores movie information, genres, and supports querying genres with the highest average ratings. The database setup, data population, and querying are implemented in Go to showcase practical use of relational databases in application development.

## Database Setup

### Schema Design
The database consists of two main tables:
1. **`movies` Table**:
   - Stores information about movies including:
     - `id`: Unique identifier for each movie (Primary Key).
     - `name`: Title of the movie.
     - `year`: Year of release.
     - `rank`: Movie rating (can be NULL).
   
   ```sql
   CREATE TABLE movies (
       id INTEGER PRIMARY KEY,
       name TEXT NOT NULL,
       year INTEGER,
       rank REAL
   );
   ```

2. **`movies_genres` Table**:
   - Links movies to their respective genres.
     - `movie_id`: Foreign key referencing `movies.id`.
     - `genre`: Genre name.

   ```sql
   CREATE TABLE movies_genres (
       movie_id INTEGER,
       genre TEXT NOT NULL,
       FOREIGN KEY (movie_id) REFERENCES movies (id)
   );
   ```

### Steps to Set Up the Database
1. **Create Tables**:
   - The `createTables` function initializes the schema by dropping existing tables (if any) and creating fresh ones.

2. **Populate Tables**:
   - **Movies Table**: The `populateMoviesTable` function reads data from `IMDB-movies.csv` and populates the `movies` table. Faulty lines in the CSV are detected and either fixed or skipped.
   - **Movies Genres Table**: The `populateMoviesGenresTable` function reads data from `IMDB-movies_genres.csv` and populates the `movies_genres` table.

3. **Clear Tables**:
   - The `clearTables` function ensures that tables are emptied before re-populating to avoid duplication.

## Query Execution
The project includes functionality to execute SQL queries on the database. For example:
- **Find the highest-rated genres**:
  - The query calculates the average rank for each genre and sorts them in descending order:

    ```sql
    SELECT genre, AVG(rank) AS avg_rank
    FROM movies
    JOIN movies_genres ON movies.id = movies_genres.movie_id
    WHERE rank IS NOT NULL
    GROUP BY genre
    ORDER BY avg_rank DESC;
    ```

- Results are displayed in a formatted table:
  ```
  Genre           Average Rank
  -----------------------------
  Comedy          8.50
  Drama           8.30
  Action          7.80
  Documentary     7.50
  ```

## Enhancements and Future Development

### Adding a Personal Collection
To expand the database, a new table can be introduced to store information about the user's personal movie collection:

**`personal_collection` Table**:
- `movie_id`: References `movies.id`.
- `location`: Location where the movie is stored (e.g., "DVD Shelf" or "Hard Drive").
- `personal_rating`: User's personal rating of the movie.

```sql
CREATE TABLE personal_collection (
    movie_id INTEGER,
    location TEXT NOT NULL,
    personal_rating REAL,
    FOREIGN KEY (movie_id) REFERENCES movies (id)
);
```

### Plans for Drawing on the Database
1. **Purpose**:
   - Provide a centralized, personal movie management system.
   - Enable detailed queries, such as finding highly rated movies within specific genres or viewing movies stored at a particular location.

2. **User Interactions**:
   - **Add Movies**: Users can add their favorite movies to the personal collection along with location and personal ratings.
   - **Search Movies**: Search movies by title, genre, or personal rating.
   - **Filter and Sort**: Filter movies based on storage location, release year, or IMDb rating.
   - **Insights and Recommendations**: Use queries to generate personalized recommendations based on favorite genres and personal ratings.

3. **Advantages Over IMDb**:
   - IMDb provides general data, but this application offers a personalized approach to movie tracking and management.
   - Customizable data fields, such as personal ratings and storage locations, cater to individual preferences.

### Further Application Development
1. **Recommendation System**:
   - Build a recommendation system based on personal ratings and watched history.
   - Use collaborative filtering or content-based filtering algorithms to suggest movies.

2. **Enhanced User Interface**:
   - Develop a web or desktop GUI using Go frameworks such as Gin (for web) or Fyne (for desktop).
   - Allow users to interact with the database visually, with forms for data input and tables for query results.

3. **Analytics Dashboard**:
   - Add visualizations, such as bar charts or pie charts, to show:
     - Most-watched genres.
     - Average ratings by genre.
     - Personal ratings vs IMDb ratings.

4. **Mobile App Integration**:
   - Extend the application to mobile platforms for on-the-go access to the personal collection.

---

## How to Run the Project
1. Place `IMDB-movies.csv` and `IMDB-movies_genres.csv` in the project directory.
2. Run the application:
   ```bash
   go run main.go
   ```
3. View the query results in the terminal.

## Conclusion
This project demonstrates how to create and interact with a relational database in Go using SQLite. By adding a personal collection and enhancing user interaction, this application has the potential to evolve into a comprehensive movie management tool, surpassing IMDb in personalized features.

