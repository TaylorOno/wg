package drinks

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"golang.org/x/net/html"
)

func init() {
	functions.HTTP("base", WG().ServeHTTP)
}

func WG() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", redirectHome)
	mux.HandleFunc("/menu", DrinkMenu)
	mux.HandleFunc("POST /history", History)
	mux.HandleFunc("GET /history/{id}", History)

	return mux
}

func redirectHome(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/wg-drinks/menu", http.StatusMovedPermanently)
}

// DrinkMenu renders the enhanced drink menu with tracking and rating support.
func DrinkMenu(w http.ResponseWriter, r *http.Request) {
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

	historyService := HistoryService{storage: NewTursoDatabase()}

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
