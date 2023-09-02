package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

type Album struct {
	id     int
	title  string
	artist string
	price  float32
}

var dbConn *pgx.Conn

func main() {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	databaseUrl := "postgres://postgres:123456@localhost:5432/golang"

	// Abra uma conexão com o banco de dados
	var err error
	dbConn, err = pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		fmt.Println("Erro ao conectar ao banco de dados", err)
		return
	}
	fmt.Println("\nConnected to Database!")
	// defer faz com que a função seja executada no final da função atual
	defer dbConn.Close(context.Background())

	var albums []Album

	albums, err = getAlbums()
	if err != nil {
		fmt.Printf("Houve um erro ao requisitar a função getAlbums: %v", err)
	}
	printAlbums(albums)

	albums, err = getAlbumsByArtist("John Coltrane")
	if err != nil {
		fmt.Printf("Houve um erro ao requisitar a função getAlbumsByArtist: %v", err)
	}
	printAlbums(albums)

	album, err := getAlbumById(5)
	fmt.Println()
	if err != nil {
		fmt.Printf("Houve um erro ao requisitar a função getAlbumById: %v", err)
	} else {
		fmt.Println(album)
	}

	// album = Album{
	// 	title:  "The Modern Sound of Betty Carter",
	// 	artist: "Betty Carter",
	// 	price:  49.99,
	// }

	// id, err := addAlbums(album)
	// fmt.Println()
	// if err != nil {
	// 	fmt.Printf("Houve um erro ao requisitar a função addAlbums: %v", err)
	// } else {
	// 	fmt.Printf("Album adicionado com sucesso, id: %d", id)
	// }

	// album = Album{
	// 	title:  "Lush Life",
	// 	artist: "John Coltrane",
	// 	price:  49.99,
	// }

	// id, err := addAlbumsTx(album)
	// fmt.Println()
	// if err != nil {
	// 	fmt.Printf("Houve um erro ao requisitar a função addAlbumsTx: %v", err)
	// } else {
	// 	fmt.Printf("Album adicionado com sucesso, id: %d", id)
	// }

	// err = deleteAlbum(6)
	// fmt.Println()
	// if err != nil {
	// 	fmt.Printf("Houve um erro ao requisitar a função deleteAlbum: %v", err)
	// }else
	// {
	// 	fmt.Printf("Nenhum erro ao requisitar a função deleteAlbum")
	// }

	album = Album{
		title:  "Lush Life",
		artist: "John Coltrane",
		price:  10,
	}

	err = updateAlbum(7, album)
	fmt.Println()
	if err != nil {
		fmt.Printf("Houve um erro ao requisitar a função updateAlbum: %v", err)
	} else {
		fmt.Printf("Nenhum erro ao requisitar a função updateAlbum")
	}
}

func updateAlbum(id int64, alb Album) error {
	sqlQuery := `UPDATE album SET title = $1, artist = $2, price = $3 WHERE id = $4`
	result, err := dbConn.Exec(context.Background(), sqlQuery, alb.title, alb.artist, alb.price, id)
	if err != nil {
		return fmt.Errorf("updateAlbum %d: %v", id, err)
	}
	if result.RowsAffected() == 0 {
		fmt.Printf("\nNenhum registro com o ID %d encontrado para atualização\n", id)
	}

	return nil
}

func deleteAlbum(id int64) error {
	sqlQuery := `DELETE FROM album WHERE id = $1`
	result, err := dbConn.Exec(context.Background(), sqlQuery, id)
	if err != nil {
		return fmt.Errorf("deleteAlbum %d: %v", id, err)
	}
	if result.RowsAffected() == 0 {
		// return fmt.Errorf("nenhum registro com o ID %d encontrado para exclusão", id)
		// log.Printf("nenhum registro com o ID %d encontrado para exclusão", id)
		fmt.Printf("\nnenhum registro com o ID %d encontrado para exclusão\n", id)
	}

	return nil
}

func printAlbums(albums []Album) {
	fmt.Println()
	for _, album := range albums {
		fmt.Println(album)
	}
}

func addAlbumsTx(alb Album) (int64, error) {
	var id int64
	tx, err := dbConn.Begin(context.Background())
	if err != nil {
		return 0, fmt.Errorf("addAlbums: %v", err)
	}
	defer tx.Rollback(context.Background())

	sqlQuery := `INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id`
	row := tx.QueryRow(context.Background(), sqlQuery, alb.title, alb.artist, alb.price)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("addAlbums: %v", err)
	}

	if err := tx.Commit(context.Background()); err != nil {
		return 0, fmt.Errorf("addAlbums: %v", err)
	}

	return id, nil
}

func addAlbums(alb Album) (int64, error) {
	var id int64
	sqlQuery := `INSERT INTO album (title, artist, price) VALUES ($1, $2, $3)`
	row := dbConn.QueryRow(context.Background(), sqlQuery, alb.title, alb.artist, alb.price)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("addAlbums: %v", err)
	}
	return id, nil
}

func getAlbums() ([]Album, error) {
	var albums []Album

	// Execute a consulta SQL
	sqlQuery := `SELECT * FROM album`
	rows, err := dbConn.Query(context.Background(), sqlQuery)
	if err != nil {
		log.Fatalf("Erro ao executar a consulta: %v", err)
	}
	defer rows.Close()

	// Iterando sobre os resultados
	for rows.Next() {
		// Declaração de variáveis para armazenar os valores retornados
		var (
			id     int
			title  string
			artist string
			price  float32
		)

		// Insira os valores retornados pelo banco nas variáveis declaradas
		if err := rows.Scan(&id, &title, &artist, &price); err != nil {
			log.Fatalf("Erro ao ler os resultados: %v", err)
		}
		alb := Album{
			id:     id,
			title:  title,
			artist: artist,
			price:  price,
		}
		albums = append(albums, alb)
	}

	// Verifica erros após a iteração sobre os resultados
	if err := rows.Err(); err != nil {
		fmt.Println("Erro após a iteração sobre os resultados:", err)
	}
	return albums, nil
}

func getAlbumsByArtist(artistName string) ([]Album, error) {
	var albums []Album

	sqlQuery := `SELECT * FROM album WHERE artist = $1`
	rows, err := dbConn.Query(context.Background(), sqlQuery, artistName)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Houve um erro ao executar a query, %s", err))
	}
	defer rows.Close()

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.id, &alb.title, &alb.artist, &alb.price); err != nil {
			return nil, fmt.Errorf("albumByArtist %q: %v", artistName, err)
		}
		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", artistName, err)
	}

	return albums, nil
}

func getAlbumById(id int64) (Album, error) {
	var album Album

	sqlQuery := `SELECT * FROM album WHERE id = $1`
	row := dbConn.QueryRow(context.Background(), sqlQuery, id)
	if err := row.Scan(&album.id, &album.title, &album.artist, &album.price); err != nil {
		if err == pgx.ErrNoRows {
			return album, fmt.Errorf("none album with id %d", id)
		}
		return album, fmt.Errorf("albumById %d: %v", id, err)
	}

	return album, nil
}
