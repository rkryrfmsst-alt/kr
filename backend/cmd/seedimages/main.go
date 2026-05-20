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
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const metAPI = "https://collectionapi.metmuseum.org/public/collection/v1"

var paintings = []struct {
	ID    int
	Title string
	Query string
}{
	{1, "Берёзовая роща в мае", "birch grove spring landscape"},
	{2, "Осенний берег", "autumn river shore landscape"},
	{3, "Тихое озеро", "lake water reflection calm"},
	{4, "Собор в тумане", "cathedral fog architecture"},
	{5, "Портрет в золотом", "woman portrait golden warm"},
	{6, "Дама с веером", "woman fan portrait"},
	{7, "Красное и чёрное", "red black abstract"},
	{8, "Диалог форм", "geometric abstract composition"},
	{9, "Ночной рынок", "night market scene figures"},
	{10, "Падающие часы", "surrealist still life objects"},
	{11, "Улица после дождя", "Paris street boulevard night"},
	{12, "Утренний рынок", "market scene"},
	{13, "Зимний закат", "winter sunset snow landscape"},
	{14, "Цветочный натюрморт", "flowers bouquet vase still life"},
	{15, "Лошади в поле", "horse"},
	{16, "Composition VII", "abstract composition geometric"},
	{17, "Сон разума", "portrait dream surreal"},
	{18, "Мегаполис", "new york city skyscraper urban panorama"},
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
	force := false
	var onlyIDs map[int]bool

	for _, arg := range os.Args[1:] {
		if arg == "--force" {
			force = true
		} else if strings.HasPrefix(arg, "--ids=") {
			onlyIDs = map[int]bool{}
			for _, s := range strings.Split(strings.TrimPrefix(arg, "--ids="), ",") {
				if id, err := strconv.Atoi(strings.TrimSpace(s)); err == nil {
					onlyIDs[id] = true
				}
			}
		}
	}

	godotenv.Load(".env")

	outDir := "./static/paintings"
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

	client := &http.Client{Timeout: 30 * time.Second}

	for _, p := range paintings {
		if onlyIDs != nil && !onlyIDs[p.ID] {
			continue
		}

		filename := fmt.Sprintf("%d.jpg", p.ID)
		dest := filepath.Join(outDir, filename)

		if _, err := os.Stat(dest); err == nil && !force {
			fmt.Printf("[%d] пропуск — файл уже есть\n", p.ID)
			continue
		}

		fmt.Printf("[%d] %s — ищу...\n", p.ID, p.Title)
		imageURL, err := findImage(client, p.Query)
		if err != nil {
			fmt.Printf("[%d] не найдено: %v\n", p.ID, err)
			continue
		}

		if err := downloadFile(client, imageURL, dest); err != nil {
			fmt.Printf("[%d] ошибка скачивания: %v\n", p.ID, err)
			continue
		}

		imagePath := fmt.Sprintf("/static/paintings/%s", filename)
		if _, err := db.Exec(`UPDATE paintings SET image_path = $1 WHERE id = $2`, imagePath, p.ID); err != nil {
			fmt.Printf("[%d] ошибка обновления БД: %v\n", p.ID, err)
			continue
		}

		fmt.Printf("[%d] готово -> %s\n", p.ID, imagePath)
		time.Sleep(200 * time.Millisecond)
	}

	fmt.Println("Готово.")
}

func findImage(client *http.Client, query string) (string, error) {
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

	for _, id := range result.ObjectIDs {
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
