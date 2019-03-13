package prcdapi

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// Server expose the PRCD API via HTTP.
type Server struct {
	db *Grimoire
}

// NewServer returns a new PRCD API Server.
func NewServer(g *Grimoire) *Server {
	s := &Server{db: g}
	r := mux.NewRouter()
	r.HandleFunc("/prcd/{section}", s.prcdHandler())
	r.HandleFunc("/prcd", s.prcdRandomHandler())
	r.HandleFunc("/sections", s.prcdSectionsHandler())
	r.HandleFunc("/", s.prcdRandomHandler())
	http.Handle("/", r)

	return s
}

func (s *Server) prcdSectionsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sections := s.db.GetSections()
		if _, err := w.Write([]byte(strings.Join(sections, " "))); err != nil {
			log.Print(err)
		}
	}
}

func (s *Server) prcdRandomHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		moccolo, err := s.db.FromRandomSection()
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
		} else {
			if _, err := w.Write([]byte(fmt.Sprintf("%s (%s)\n", moccolo.Text, moccolo.Author))); err != nil {
				log.Print(err)
			}
		}
	}
}

func (s *Server) prcdHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sectionName, ok := vars["section"]
		if !ok {
			http.Error(w, "Internal error: Section not found", http.StatusInternalServerError)
			return
		}

		moccolo, err := s.db.FromSection(sectionName)
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
		} else {
			if _, err := w.Write([]byte(fmt.Sprintf("%s (%s)\n", moccolo.Text, moccolo.Author))); err != nil {
				log.Print(err)
			}
		}
	}
}

// Serve starts the HTTP handler.
func (s *Server) Serve(address string) {
	log.Fatal(http.ListenAndServe(address, nil))
}
