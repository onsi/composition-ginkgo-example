package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"sync"
)

// a very hokey and terrible key-value store

var r *rand.Rand

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func beSlow() {
	time.Sleep(time.Duration(r.Intn(200)) * time.Millisecond)
}

func main() {
	lock := &sync.Mutex{}
	data := map[string]string{}

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {})

	http.HandleFunc("/set", func(w http.ResponseWriter, r *http.Request) {
		beSlow()
		key := r.URL.Query().Get("key")
		value := r.URL.Query().Get("value")
		lock.Lock()
		data[key] = value
		lock.Unlock()
	})

	http.HandleFunc("/get", func(w http.ResponseWriter, r *http.Request) {
		beSlow()
		key := r.URL.Query().Get("key")
		lock.Lock()
		value, ok := data[key]
		lock.Unlock()

		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		fmt.Fprint(w, value)
	})

	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		beSlow()
		key := r.URL.Query().Get("key")
		lock.Lock()
		delete(data, key)
		lock.Unlock()
	})

	http.HandleFunc("/get-prefix", func(w http.ResponseWriter, r *http.Request) {
		beSlow()
		prefix := r.URL.Query().Get("prefix")
		values := []string{}
		lock.Lock()
		for key, value := range data {
			if strings.HasPrefix(key, prefix) {
				values = append(values, value)
			}
		}
		lock.Unlock()

		json.NewEncoder(w).Encode(values)
	})

	http.HandleFunc("/delete-prefix", func(w http.ResponseWriter, r *http.Request) {
		beSlow()
		prefix := r.URL.Query().Get("prefix")
		toDelete := []string{}
		lock.Lock()
		for key := range data {
			if strings.HasPrefix(key, prefix) {
				toDelete = append(toDelete, key)
			}
		}

		for _, key := range toDelete {
			delete(data, key)
		}
		lock.Unlock()
	})

	log.Fatal(http.ListenAndServe(":9999", nil))
}
