package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql" // better practice: "_" for loosely-coupled
)

var db *sql.DB // bud prctice: Package variable. Just for simplicity.

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func main() {

	// Capture connection properties.
	//cfg := mysql.Config{
	//	User:                 os.Getenv("DBUSER"),
	//	Passwd:               os.Getenv("DBPASS"),
	//	Net:                  "tcp",
	//	Addr:                 "127.0.0.1:3306",
	//	DBName:               "recordings",
	//	AllowNativePasswords: true, // added for native password authentication(?)
	//}
	//formatDSN := cfg.FormatDSN()

	// good practice: Create a structure and a method like Config or FormatDSN above by myself because mysql package is "_"
	// Below is just for simplicity.
	formatDSN := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/recordings?checkConnLiveness=false&maxAllowedPacket=0", os.Getenv("DBUSER"), os.Getenv("DBPASS"))

	// Get a database handle.
	var err error
	db, err = sql.Open("mysql", formatDSN) // Validate arguments without creating a connection.
	if err != nil {
		log.Fatal("Open: ", err)
	}

	pingErr := db.Ping() // Confirm connection.
	if pingErr != nil {
		log.Fatal("Ping: ", pingErr)
	}
	fmt.Println("Connected!")

	// Test the query. (expect multiple rows)
	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)

	// Hard-code ID 2 here to test the query. (expect a row)
	alb, err := albumByID(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", alb)

	// Test the execute INSERT
	albID, err := addAlbum(Album{
		Title:  "The Modern Sound of Betty Carter",
		Artist: "Betty Carter",
		Price:  49.99,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added album: %v\n", albID)
}

// albumsByArtist queries for albums that have the specified artist name.
func albumsByArtist(name string) ([]Album, error) {
	// An albums slice to hold data from returned rows.
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = ?", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}

// albumByID queries for the album with the specified ID.
func albumByID(id int64) (Album, error) {
	// An album to hold data from the returned row.
	var alb Album

	row := db.QueryRow("SELECT * FROM album WHERE id = ?", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == sql.ErrNoRows {
			return alb, fmt.Errorf("albumsById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumsById %d: %v", id, err)
	}
	return alb, nil
}

// addAlbum adds the specified album to the database,
// returning the album ID of the new entry
func addAlbum(alb Album) (int64, error) {
	result, err := db.Exec("INSERT INTO album (title, artist, price) VALUES (?, ?, ?)", alb.Title, alb.Artist, alb.Price)
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addAlbum: %v", err)
	}
	return id, nil
}
