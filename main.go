package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// var cookieHandler = securecookie.New(
// 	securecookie.GenerateRandomKey(64),
// 	securecookie.GenerateRandomKey(32))

func getUserName(request *http.Request) (userName string) {
	cookie, err := request.Cookie("session")
	if err != nil {
		return "Error"
	} else {
		userName = cookie.Value
	}

	// if cookie, err := request.Cookie("session"); err == nil {
	// 	cookieValue := make(map[string]string)
	// 	if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
	// 		userName = cookieValue["session"]
	log.Println(">>>>>> " + userName)
	// 	}
	// }
	return userName
}

type APIControllerV1 struct{}

func (c *APIControllerV1) getUser(w http.ResponseWriter, r *http.Request) {
	userName := getUserName(r)
	if userName == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	// Handle request
}

type APIControllerV2 struct{}

func (c *APIControllerV2) getUser(w http.ResponseWriter, r *http.Request) {
	userName := getUserName(r)
	log.Println(userName)

	if userName == "" {
		http.Error(w, "Unauthorized v2", http.StatusUnauthorized)
		return
	} else if userName == "cookie_value" {
		resp := make(map[string]string)
		resp["message"] = "Success"
		resp["cookie value"] = userName
		jsonResp, err := json.Marshal(resp)
		if err != nil {
			log.Fatalf("Error happened in JSON marshal. Err: %s", err)
		}
		w.Write(jsonResp)
	}
	// Handle request
}

func (c *APIControllerV1) getCookie(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err == http.ErrNoCookie {
		// Cookie doesn't exist, so create a new one and set it in the response
		cookie = &http.Cookie{
			Name:   "session",
			Value:  "cookie_value",
			Path:   "/",
			MaxAge: 3600, // cookie expiration time in seconds
		}
		http.SetCookie(w, cookie)
	} else if err != nil {
		// Handle other errors if any
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Use the cookie value
	log.Printf("Cookie value: %s", cookie.Value)
}

func main() {
	router := mux.NewRouter()

	// API version 1
	v1 := router.PathPrefix("/api/v1").Subrouter()
	v1Controller := &APIControllerV1{}

	v1.HandleFunc("/getcookie", v1Controller.getCookie).Methods("POST")
	// API version 1
	// v1.Use(authMiddleware) // Apply authentication middleware to all routes in this version
	v1.HandleFunc("/user", v1Controller.getUser).Methods("GET")

	// API version 2
	v2 := router.PathPrefix("/api/v2").Subrouter()
	v2.Use(authMiddleware) // Apply authentication middleware to all routes in this version
	v2Controller := &APIControllerV2{}
	v2.HandleFunc("/user", v2Controller.getUser).Methods("GET")

	// Start the server
	log.Print("server started on 8080")
	http.ListenAndServe(":8080", router)

}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := r.Cookie("session"); err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
