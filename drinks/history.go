package drinks

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha1"
	"embed"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
)

//go:embed templates/*
var templateFS embed.FS

type HistoryManager interface {
	Save(context.Context, string, []byte) error
	GetById(ctx context.Context, id string) ([]byte, error)
}

type HistoryService struct {
	storage HistoryManager
}

type DrinkHistory map[string]DrinkDetails

type DrinkDetails struct {
	Name        string
	Date        time.Time
	Description string
	Score       []bool
}

func (d *DrinkDetails) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}

	var drinkDetails map[string]any
	if err := json.Unmarshal(data, &drinkDetails); err != nil {
		return err
	}

	rating, _ := strconv.Atoi(fmt.Sprintf("%v", drinkDetails["score"]))
	score := make([]bool, 5)
	for i := range rating {
		score[i] = true
	}

	date, err := time.Parse(time.RFC3339, fmt.Sprintf("%v", drinkDetails["date"]))
	if err != nil {
		date = time.Now()
	}

	*d = DrinkDetails{
		Name:        fmt.Sprintf("%v", drinkDetails["name"]),
		Date:        date,
		Description: fmt.Sprintf("%v", drinkDetails["description"]),
		Score:       score,
	}

	return nil
}

func (h *HistoryService) SaveHistory(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Max-Age", "3600")

	id, err := generateId()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	history, err := compress(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	if len(history) >= 102_400 {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("exceeds max size")
	}

	err = h.storage.Save(r.Context(), id, history)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	return writeShareResponse(w, id)
}

func writeShareResponse(w http.ResponseWriter, id string) error {
	type resp struct {
		ID string `json:"id"`
	}

	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(resp{ID: id})
	if err != nil {
		return err
	}

	return nil
}

func (h *HistoryService) ViewHistory(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Max-Age", "3600")

	id := r.PathValue("id")
	history, err := h.storage.GetById(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	historyJson, err := decompress(history)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return err
	}

	slog.Debug("parsing drink history", slog.Any("data", string(historyJson)))
	var drinkHistory DrinkHistory
	err = json.Unmarshal(historyJson, &drinkHistory)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return fmt.Errorf("failed to parse history json: %w", err)
	}

	tmpl := template.Must(template.ParseFS(templateFS, "templates/history.gohtml"))
	_ = tmpl.Execute(w, drinkHistory)

	return err
}

func generateId() (string, error) {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}

	hasher := sha1.New()
	_, err = hasher.Write([]byte(newUUID.String()))
	if err != nil {
		return "", err
	}

	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha[:10], nil
}

func compress(body io.Reader) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	encoder := base64.NewEncoder(base64.StdEncoding, gz)

	size, err := io.Copy(encoder, body)
	slog.Debug("compressing", slog.Any("original-size", size))
	if err != nil {
		return nil, err
	}

	err = encoder.Close()
	if err != nil {
		return nil, err
	}

	err = gz.Close()
	if err != nil {
		return nil, err
	}

	slog.Debug("compressing", slog.Any("compressed-size", len(buf.Bytes())))
	return buf.Bytes(), nil
}

func decompress(data []byte) ([]byte, error) {
	slog.Debug("decompressing", slog.Any("compressed-size", len(data)))
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	out, err := io.ReadAll(base64.NewDecoder(base64.StdEncoding, gz))
	if err != nil {
		return nil, err
	}

	err = gz.Close()
	if err != nil {
		return nil, err
	}

	slog.Debug("decompressing", slog.Any("decompressed-size", len(out)))
	return out, nil
}
