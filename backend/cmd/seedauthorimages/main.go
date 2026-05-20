package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const metAPI = "https://collectionapi.metmuseum.org/public/collection/v1"

type authorRow struct {
	ID        int    `db:"id"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}

type searchResult struct {
	Total     int   `json:"total"`
	ObjectIDs []int `json:"objectIDs"`
}

type metObject struct {
	IsPublicDomain    bool   `json:"isPublicDomain"`
	PrimaryImage      string `json:"primaryImage"`
	PrimaryImageSmall string `json:"primaryImageSmall"`
}

func main() {
	force := len(os.Args) > 1 && os.Args[1] == "--force"

	godotenv.Load(".env")

	outDir := "./static/authors"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		getenv("DB_HOST", "localhost"),
		getenv("DB_PORT", "5432"),
		getenv("DB_USER", "postgres"),
		getenv("DB_PASSWORD", ""),
		getenv("DB_NAME", "webgallery"),
		getenv("DB_SSLMODE", "disable"),
	)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("не удалось подключиться к БД:", err)
	}
	defer db.Close()

	var authors []authorRow
	if err := db.Select(&authors, `SELECT id, first_name, last_name FROM authors ORDER BY id`); err != nil {
		log.Fatal("ошибка запроса авторов:", err)
	}

	client := &http.Client{Timeout: 30 * time.Second}

	for _, a := range authors {
		filename := fmt.Sprintf("%d.jpg", a.ID)
		dest := filepath.Join(outDir, filename)

		if _, err := os.Stat(dest); err == nil && !force {
			fmt.Printf("[%d] пропуск — файл уже есть\n", a.ID)
			continue
		}

		query := portraitQuery(a.FirstName)
		fmt.Printf("[%d] %s %s — ищу (%s)...\n", a.ID, a.FirstName, a.LastName, query)

		imageURL, err := findImage(client, query, a.ID)
		if err != nil {
			fmt.Printf("[%d] не найдено: %v\n", a.ID, err)
			continue
		}

		if err := downloadFile(client, imageURL, dest); err != nil {
			fmt.Printf("[%d] ошибка скачивания: %v\n", a.ID, err)
			continue
		}

		photoPath := fmt.Sprintf("/static/authors/%s", filename)
		if _, err := db.Exec(`UPDATE authors SET photo_path = $1 WHERE id = $2`, photoPath, a.ID); err != nil {
			fmt.Printf("[%d] ошибка обновления БД: %v\n", a.ID, err)
			continue
		}

		fmt.Printf("[%d] готово -> %s\n", a.ID, photoPath)
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println("Готово.")
}

// portraitQuery выбирает запрос по полу автора:
// русские женские имена оканчиваются на 'а' или 'я'.
func portraitQuery(firstName string) string {
	runes := []rune(firstName)
	if len(runes) > 0 {
		last := runes[len(runes)-1]
		if last == 'а' || last == 'я' {
			return "woman portrait painting"
		}
	}
	return "man portrait painting"
}

func findImage(client *http.Client, query string, offset int) (string, error) {
	searchURL := fmt.Sprintf(
		"%s/search?q=%s&hasImages=true&medium=Paintings",
		metAPI, url.QueryEscape(query),
	)
	resp, err := client.Get(searchURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result searchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Total == 0 {
		return "", fmt.Errorf("ничего не найдено")
	}

	// Используем offset чтобы разные авторы получали разные портреты
	start := (offset * 7) % len(result.ObjectIDs)
	ids := append(result.ObjectIDs[start:], result.ObjectIDs[:start]...)

	for _, id := range ids {
		obj, err := getObject(client, id)
		if err != nil {
			continue
		}
		if obj.IsPublicDomain && obj.PrimaryImage != "" {
			if obj.PrimaryImageSmall != "" {
				return obj.PrimaryImageSmall, nil
			}
			return obj.PrimaryImage, nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return "", fmt.Errorf("нет подходящих изображений")
}

func getObject(client *http.Client, id int) (*metObject, error) {
	resp, err := client.Get(fmt.Sprintf("%s/objects/%d", metAPI, id))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var obj metObject
	if err := json.NewDecoder(resp.Body).Decode(&obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

func downloadFile(client *http.Client, imageURL, dest string) error {
	resp, err := client.Get(imageURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
