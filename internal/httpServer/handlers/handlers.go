package handlers

import (
	"crypto/sha512"
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Handlers struct {
	shutdownCh chan bool
	infoMap *InfoMap
	metrics *Stats
	waitToStore time.Duration
}

func NewHandlers(shutdownCh chan bool) *Handlers {
	return &Handlers{
		infoMap: NewInfoMap(),
		metrics: &Stats{},
		shutdownCh: shutdownCh,
		waitToStore: 5 * time.Second,
	}
}

func (h *Handlers) Hash(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	r.ParseForm()
	password, ok := r.Form["password"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	key := h.infoMap.GetKey()
	keyStr := strconv.Itoa(key)

	// TODO figure out a way to do testing of graceful shutdowns in unit tests, if possible
	//time.Sleep(20 * time.Second)
	go func() {
		time.Sleep(h.waitToStore)
		sha_512 := sha512.New()
		sha_512.Write([]byte(password[0]))
		b64 := base64.StdEncoding.EncodeToString(sha_512.Sum(nil))
		h.infoMap.Store(key, b64)
	}()

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(keyStr))

	taken := time.Since(start)
	h.metrics.Update(taken)
}

func (h *Handlers) HashID(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(strings.Trim(r.URL.Path, "/"), "hash/")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	hashedPassword := h.infoMap.Load(id)

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(hashedPassword))
}

func (h *Handlers) Shutdown(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	h.shutdownCh <- true
	w.WriteHeader(http.StatusAccepted)
}

func (h *Handlers) Stats(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	response, err := h.metrics.ToJSON()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write(response)
}
