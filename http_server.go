package main

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

type HttpServer struct {
	port string
}

type Message struct {
	BPM int
}

func (s *HttpServer) run() error {

	router := makeMuxRouter()
	httpAddr := s.port
	log.Printf("server listen port %s", httpAddr)

	server := &http.Server{
		Addr:           ":" + httpAddr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if err := server.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func makeMuxRouter() http.Handler {
	muxRouter := mux.NewRouter()
	muxRouter.HandleFunc("/", handleGetBlockchain).Methods("GET")
	muxRouter.HandleFunc("/", handleWriteBlock).Methods("POST")
	return muxRouter
}

func handleGetBlockchain(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, r, http.StatusOK, Blockchain)
}

func handleWriteBlock(w http.ResponseWriter, r *http.Request) {
	var msg Message
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&msg); err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, r.Body)
		return
	}

	defer r.Body.Close()

	newBlock, err := generateNewBlock(Blockchain[len(Blockchain)-1], msg.BPM)

	if err != nil {
		respondWithJSON(w, r, http.StatusInternalServerError, msg)
	}

	if isBlockValid(Blockchain[len(Blockchain)-1], newBlock) {
		newBlockchain := appendChain(newBlock)
		replaceChain(newBlockchain)
		spew.Dump(Blockchain)
	}

	respondWithJSON(w, r, http.StatusCreated, Blockchain)
}

func respondWithJSON(w http.ResponseWriter, r *http.Request, code int, payload interface{}) {
	response, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(code)
	w.Write(response)
}
