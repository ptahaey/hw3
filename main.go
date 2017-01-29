package main

import (
	"net/http"
	"strconv"
	"hw3/hook"
)

func init() {
	hook.Init()
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hooks", func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			rw.Write([]byte("\n"))
		}()
		inputValues := r.URL.Query()
		if r.Method == http.MethodGet {
			if len(inputValues) > 2 {
				http.Error(rw, "Wrong number of input parameters", http.StatusBadRequest)
				return
			}
			url, topic := inputValues.Get("url"), inputValues.Get("topic")
			hooksList := hook.GetHooks(url, topic)
			if len(hooksList) > 0 {
				rw.Header().Set("Content-Type", "application/json")
				rw.Write(hooksList)
				return
			}
			rw.WriteHeader(http.StatusNoContent)
			return
		}
		if r.Method == http.MethodPut {
			if len(inputValues) != 3 {
				http.Error(rw, "Wrong number of input parameters", http.StatusBadRequest)
				return
			}
			url, topic := inputValues.Get("url"), inputValues.Get("topic")
			max_failures, _ := strconv.Atoi(inputValues.Get("max_failures"))

			success, error := hook.CreateHook(url, topic, max_failures)
			if !success {
				http.Error(rw, error.Error(), http.StatusBadRequest)
				return
			}

			rw.WriteHeader(http.StatusOK)
			return
		}

		if r.Method == http.MethodDelete {
			if len(inputValues) > 2 {
				http.Error(rw, "Wrong number of input parameters", http.StatusBadRequest)
				return
			}
			url, topic := inputValues.Get("url"), inputValues.Get("topic")
			hook.DeleteHooks(url, topic)

			rw.WriteHeader(http.StatusNoContent)
			return
		}


	})
	mux.HandleFunc("/topics/", func(rw http.ResponseWriter, r *http.Request) {
		defer func() {
			rw.Write([]byte("\n"))
		}()
		inputValues := r.URL.Query()
		if r.Method == http.MethodPut {
			//topic := strings.TrimPrefix(string(r.URL.Path), "/topics") ;
			if len(inputValues) !=1 {
				http.Error(rw, "Wrong number of input parameters", http.StatusBadRequest)
				return
			}
			text := inputValues.Get("text")
			urlsList := hook.PutTopics(string(r.URL.Path), text)
			rw.WriteHeader(http.StatusAccepted)
			if len(urlsList) > 0 {
				rw.Header().Set("Content-Type", "application/json")
				rw.Write(urlsList)
			}
			return
		}
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return

	})
	http.ListenAndServe(":8844", mux)
}
