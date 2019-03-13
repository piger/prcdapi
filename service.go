package prcdapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func respondJSON(w http.ResponseWriter, payload interface{}, code int) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write([]byte(response))
	}
}

func respondJSONError(w http.ResponseWriter, message string, code int) {
	respondJSON(w, map[string]string{"error": message}, code)
}

// Server expose the PRCD API via HTTP.
type Server struct {
	db *Grimoire
}

// NewServer returns a new PRCD API Server.
func NewServer(g *Grimoire) *Server {
	s := &Server{db: g}
	r := mux.NewRouter()

	rJSON := r.Headers("Content-Type", "application/json").Subrouter()
	rJSON.HandleFunc("/prcd", s.prcdRandomJSONHandler())
	rJSON.HandleFunc("/prcd/{section}", s.prcdJSONHandler())
	rJSON.HandleFunc("/sections", s.prcdSectionsJSONHandler())
	rJSON.HandleFunc("/", s.prcdRandomJSONHandler())

	rJSONa := r.Headers("Accept", "application/json").Subrouter()
	rJSONa.HandleFunc("/prcd", s.prcdRandomJSONHandler())
	rJSONa.HandleFunc("/prcd/{section}", s.prcdJSONHandler())
	rJSONa.HandleFunc("/sections", s.prcdSectionsJSONHandler())
	rJSONa.HandleFunc("/", s.prcdRandomJSONHandler())

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
		if _, err := w.Write([]byte(strings.Join(sections, " ") + "\n")); err != nil {
			log.Print(err)
		}
	}
}

func (s *Server) prcdSectionsJSONHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sections := s.db.GetSections()
		respondJSON(w, sections, 200)
	}
}

func (s *Server) prcdRandomHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		moccolo, section, err := s.db.FromRandomSection()
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal error: fetching moccolo.", http.StatusInternalServerError)
		} else {
			if _, err := w.Write([]byte(fmt.Sprintf("%s (%s) [%s]\n", moccolo.Text, moccolo.Author, section))); err != nil {
				log.Print(err)
			}
		}
	}
}

func (s *Server) prcdRandomJSONHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		moccolo, section, err := s.db.FromRandomSection()
		if err != nil {
			log.Print(err)
			respondJSONError(w, "Error fetching a random Moccolo.", http.StatusInternalServerError)
		} else {
			moccolo.Section = section
			respondJSON(w, moccolo, 200)
		}
	}
}

func (s *Server) prcdHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sectionName, ok := vars["section"]
		if !ok {
			http.Error(w, "Internal error: Section not found.", http.StatusInternalServerError)
			return
		}

		moccolo, err := s.db.FromSection(sectionName)
		if err != nil {
			log.Print(err)
			http.Error(w, "Internal error: fetching moccolo from section.", http.StatusInternalServerError)
		} else {
			if _, err := w.Write([]byte(fmt.Sprintf("%s (%s)\n", moccolo.Text, moccolo.Author))); err != nil {
				log.Print(err)
			}
		}
	}
}

func (s *Server) prcdJSONHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sectionName, ok := vars["section"]
		if !ok {
			respondJSONError(w, "Internal error: Section not found.", http.StatusInternalServerError)
			return
		}

		moccolo, err := s.db.FromSection(sectionName)
		if err != nil {
			log.Print(err)
			respondJSONError(w, "Internal error: fetching moccolo from section.", http.StatusInternalServerError)
		} else {
			respondJSON(w, moccolo, 200)
		}
	}
}

// Serve starts the HTTP handler.
func (s *Server) Serve(address string) {
	log.Fatal(http.ListenAndServe(address, nil))
}
