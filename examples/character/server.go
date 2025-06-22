package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/karashiiro/bingode"
	"github.com/xivapi/godestone/v2"
)

type Server struct {
	db      *Database
	scraper *godestone.Scraper
	port    string
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Data      interface{} `json:"data"`
	Cached    bool        `json:"cached"`
	Timestamp time.Time   `json:"timestamp"`
}

func NewServer(dbPath string, port string) (*Server, error) {
	db, err := NewDatabase(dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	scraper := godestone.NewScraper(bingode.New(), godestone.EN)

	return &Server{
		db:      db,
		scraper: scraper,
		port:    port,
	}, nil
}

func (s *Server) getCharacterHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		s.writeError(w, http.StatusBadRequest, "Invalid character ID")
		return
	}

	characterID := uint32(id)

	// 嘗試從快取中獲取
	cachedData, found, err := s.db.GetCharacter(characterID)
	if err != nil {
		log.Printf("Database error: %v", err)
		s.writeError(w, http.StatusInternalServerError, "Database error")
		return
	}

	if found {
		// 返回快取的資料
		var character interface{}
		if err := json.Unmarshal(cachedData, &character); err != nil {
			log.Printf("Failed to unmarshal cached data: %v", err)
			s.writeError(w, http.StatusInternalServerError, "Failed to process cached data")
			return
		}

		response := SuccessResponse{
			Data:      character,
			Cached:    true,
			Timestamp: time.Now(),
		}
		s.writeJSON(w, http.StatusOK, response)
		return
	}

	// 快取中沒有，從 Lodestone 抓取
	log.Printf("Fetching character %d from Lodestone", characterID)
	character, err := s.scraper.FetchCharacter(characterID)
	if err != nil {
		log.Printf("Failed to fetch character %d: %v", characterID, err)
		s.writeError(w, http.StatusNotFound, "Character not found or fetch failed")
		return
	}

	// 儲存到快取
	if err := s.db.SaveCharacter(characterID, character); err != nil {
		log.Printf("Failed to cache character %d: %v", characterID, err)
		// 繼續處理，只記錄錯誤
	}

	response := SuccessResponse{
		Data:      character,
		Cached:    false,
		Timestamp: time.Now(),
	}
	s.writeJSON(w, http.StatusOK, response)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	}
	s.writeJSON(w, http.StatusOK, response)
}

func (s *Server) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func (s *Server) writeError(w http.ResponseWriter, statusCode int, message string) {
	response := ErrorResponse{Error: message}
	s.writeJSON(w, statusCode, response)
}

func (s *Server) Start() error {
	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/api/character/{id:[0-9]+}", s.getCharacterHandler).Methods("GET")
	r.HandleFunc("/health", s.healthHandler).Methods("GET")

	// Static route for root
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "FF14 Character Data API Server\n")
		fmt.Fprintf(w, "Usage: GET /api/character/{id}\n")
		fmt.Fprintf(w, "Health: GET /health\n")
	}).Methods("GET")

	log.Printf("Starting server on port %s", s.port)
	log.Printf("API endpoint: http://localhost:%s/api/character/{id}", s.port)
	
	return http.ListenAndServe(":"+s.port, r)
}

func (s *Server) Close() error {
	return s.db.Close()
}