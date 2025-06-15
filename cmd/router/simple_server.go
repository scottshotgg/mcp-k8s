package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"git.6740.io/scottshotgg/mcp-k8s/pkg/router"
)

type ChatReq struct {
	Text string `json:"text"`
}

type ChatRes struct {
	Text string `json:"text"`
}

func chatHandler(rr *router.Router) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}

		var (
			req ChatReq
			res ChatRes
		)
		// 	bodyBytes, err = io.ReadAll(r.Body)
		// )

		// if err != nil {
		// 	fmt.Println("WOAH THERE WAS AN ERROR 0:", err)
		// 	return
		// }

		// err = json.Unmarshal(bodyBytes, &req)
		var err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			fmt.Println("WOAH THERE WAS AN ERROR 1:", err)
			return
		}

		res.Text = rr.HandleText(req.Text)

		b, err := json.Marshal(&res)
		if err != nil {
			fmt.Println("WOAH THERE WAS AN ERROR 3:", err)
			return
		}

		w.Write(b)
	}
}

// Just a dead simple server to get a dead simple frontend going
func server(rr *router.Router) error {
	http.HandleFunc("/chat", chatHandler(rr))

	port := ":9090"
	return http.ListenAndServe(port, nil)
}
