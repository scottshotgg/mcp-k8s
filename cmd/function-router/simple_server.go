package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ChatReq struct {
	Text string `json:"text"`
}

type ChatRes struct {
	Text string `json:"text"`
}

func chatHandler(router *Router) func(w http.ResponseWriter, r *http.Request) {
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

		res.Text = router.handleText(req.Text)

		b, err := json.Marshal(&res)
		if err != nil {
			fmt.Println("WOAH THERE WAS AN ERROR 3:", err)
			return
		}

		w.Write(b)
	}
}

// Just a dead simple server to get a dead simple frontend going
func server(r *Router) error {
	http.HandleFunc("/chat", chatHandler(r))

	port := ":9090"
	fmt.Printf("Server is running at http://localhost%s/\n", port)
	return http.ListenAndServe(port, nil)
}

func description(desc *string) string {
	if desc != nil {
		return *desc
	}

	return ""
}
