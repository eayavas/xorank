package main

import (
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"sort"
	"time"

	"github.com/eayavas/XORank/internal/logic"
	"github.com/eayavas/XORank/internal/models"
	"github.com/eayavas/XORank/internal/store"
)

var (
	db   *store.Storage
	tmpl *template.Template
)

func main() {
	var err error
	db, err = store.NewStorage("./elo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.DB.Close()

	funcMap := template.FuncMap{"PlusOne": func(i int) int { return i + 1 }}
	tmpl = template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))

	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/logout", handleLogout)
	http.HandleFunc("/", authMiddleware(handleIndex))
	http.HandleFunc("/vote", authMiddleware(handleVote))
	http.HandleFunc("/results", authMiddleware(handleResults))

	log.Println("XORANK ONLINE: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// --- HANDLERS ---

func handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl.ExecuteTemplate(w, "login.html", nil)
		return
	}
	code := r.FormValue("passcode")

	if db.CheckPasscode(code) {
		http.SetCookie(w, &http.Cookie{
			Name: "xorank_auth", Value: code, Expires: time.Now().Add(24 * time.Hour), Path: "/",
		})
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		tmpl.ExecuteTemplate(w, "login.html", "INVALID PASSCODE")
	}
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{Name: "xorank_auth", MaxAge: -1})
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

type PageData struct {
	Left, Right *models.Item
	Finished    bool
	User        string
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	passcode := getUserFromCookie(r)
	items, _ := db.GetAllItems()

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(items), func(i, j int) { items[i], items[j] = items[j], items[i] })

	var left, right *models.Item
	found := false

	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			pairKey := getPairKey(items[i].ID, items[j].ID)
			if !db.HasVoted(passcode, pairKey) {
				left, right = items[i], items[j]
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	tmpl.ExecuteTemplate(w, "index.html", PageData{Left: left, Right: right, Finished: !found, User: passcode})
}

func handleVote(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	passcode := getUserFromCookie(r)
	winnerID := r.FormValue("winner")
	loserID := r.FormValue("loser")
	pairKey := getPairKey(winnerID, loserID)

	if db.HasVoted(passcode, pairKey) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	items, _ := db.GetAllItems()
	var winner, loser *models.Item
	for _, i := range items {
		if i.ID == winnerID {
			winner = i
		}
		if i.ID == loserID {
			loser = i
		}
	}

	if winner != nil && loser != nil {
		logic.Calculate(winner, loser)
		db.SaveVote(passcode, pairKey, winner, loser)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func handleResults(w http.ResponseWriter, r *http.Request) {
	items, _ := db.GetAllItems()
	sort.Slice(items, func(i, j int) bool {
		return items[i].Rating > items[j].Rating
	})
	tmpl.ExecuteTemplate(w, "results.html", items)
}

// --- HELPER FUNCTIONS ---

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("xorank_auth")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next(w, r)
	}
}

func getUserFromCookie(r *http.Request) string {
	c, _ := r.Cookie("xorank_auth")
	return c.Value
}

func getPairKey(id1, id2 string) string {
	if id1 < id2 {
		return id1 + "|" + id2
	}
	return id2 + "|" + id1
}
