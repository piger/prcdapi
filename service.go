package prcdapi

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Server struct {
	db *Grimoire
}

func NewServer(g *Grimoire) *Server {
	s := &Server{db: g}
	r := mux.NewRouter()
	r.HandleFunc("/prcd/{section}", s.prcdHandler())
	r.HandleFunc("/prcd", s.prcdRandomHandler())
	r.HandleFunc("/sections", s.prcdSectionsHandler())
	http.Handle("/", r)

	return s
}

func (s *Server) prcdSectionsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sections := make([]string, 0, len(s.db.Sections))
		for section := range s.db.Sections {
			sections = append(sections, section)
		}

		w.Write([]byte(strings.Join(sections, " ")))
	}
}

func (s *Server) prcdRandomHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		moccolo, err := s.db.FromRandomSection()
		if err != nil {
			log.Print(err)
			return
		}
		w.Write([]byte(fmt.Sprintf("%s (%s)\n", moccolo.Text, moccolo.Author)))
	}
}

func (s *Server) prcdHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sectionName, ok := vars["section"]
		if !ok {
			log.Print("no section?")
			return
		}
		moccolo, err := s.db.FromSection(sectionName)
		if err != nil {
			log.Print(err)
			return
		}
		w.Write([]byte(fmt.Sprintf("%s (%s)\n", moccolo.Text, moccolo.Author)))
	}
}

func (s *Server) Serve() {
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
