package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jnka9755/go-05REGISTRATION/internal/registration"
	"github.com/jnka9755/go-05REGISTRATION/package/bootstrap"
	"github.com/jnka9755/go-05REGISTRATION/package/handler"
	"github.com/joho/godotenv"
)

func main() {

	_ = godotenv.Load()
	log := bootstrap.InitLooger()
	db, err := bootstrap.DBConnection()

	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	registrationRepository := registration.NewRepository(log, db)
	registrationBusiness := registration.NewBusiness(log, registrationRepository)

	handler := handler.NewCourseHTTPServer(ctx, registration.MakeEndpoints(registrationBusiness))

	port := os.Getenv("PORT")
	address := fmt.Sprintf("127.0.0.1:%s", port)

	server := http.Server{
		Handler:      accessControl(handler),
		Addr:         address,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	errCh := make(chan error)

	go func() {
		log.Println("Listen in ", address)
		errCh <- server.ListenAndServe()
	}()

	err = <-errCh

	if err != nil {
		log.Fatal(err)
	}
}

func accessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS, HEAD, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept,Authorization,Cache-Control,Content-Type,DNT,If-Modified-Since,Keep-Alive,Origin,User-Agent,X-Requested-With")

		if r.Method == "OPTIONS" {
			return
		}

		h.ServeHTTP(w, r)
	})
}
