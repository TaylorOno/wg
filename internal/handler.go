package drinks

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"golang.org/x/net/html"
	"wg/internal/storage"
)

func init() {
	functions.HTTP("drinks", Drinks)
	functions.HTTP("history", History)
}

// Drinks renders the enhanced drink menu with tracking and rating support.
func Drinks(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Max-Age", "3600")

	doc, err := AddDrinkExtensions()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	err = html.Render(w, doc)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, err)
		return
	}

	return
}

// History supports sharing of users drink history.
//
// POST will save the users history to an external data store.
// GET will render history given a share key.
func History(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Max-Age", "3600")

	historyService := HistoryService{storage: storage.NewTursoDatabase()}

	switch r.Method {
	case http.MethodPost:
		if err := historyService.SaveHistory(w, r); err != nil {
			slog.Error(err.Error())
			return
		}
	case http.MethodGet:
		if err := historyService.ViewHistory(w, r); err != nil {
			slog.Error(err.Error())
			return
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
