package main

import (
	avatar "consume-api/controllers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", avatar.Index)
	http.HandleFunc("/action", avatar.Add)
	http.HandleFunc("/action/store", avatar.Store)
	http.HandleFunc("/action/update", avatar.Update)
	http.HandleFunc("/action/execute_update", avatar.ExecuteUpdate)
	http.HandleFunc("/action/delete", avatar.Delete)

	log.Print("Server is running on : http://localhost:2525")
	log.Fatal(http.ListenAndServe(":2525", nil))
}
