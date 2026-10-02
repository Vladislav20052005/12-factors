package make_server

import (
	"ad/ads"
	"ad/auth"
	"fmt"
	"log"
	"net/http"
	"os"
)

func Create() *http.Server {
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatalln("PORT is not defined\n")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST   /api/auth/login", auth.LogIn)
	mux.HandleFunc("POST   /api/auth/signup", auth.CreateUser)
	mux.HandleFunc("DELETE /api/auth", auth.ValidationMiddleware(auth.RemoveUser))
	mux.HandleFunc("GET    /api/auth", auth.ValidationMiddleware(auth.GetUser))

	mux.HandleFunc("GET    /api/ads", auth.ValidationMiddleware(ads.GetAds))
	mux.HandleFunc("GET    /api/ads/{id}", auth.ValidationMiddleware(ads.GetAd))
	mux.HandleFunc("POST   /api/ads", auth.ValidationMiddleware(ads.CreateAd))
	mux.HandleFunc("PUT    /api/ads/{id}", auth.ValidationMiddleware(ads.EditAd))
	mux.HandleFunc("DELETE /api/ads/{id}", auth.ValidationMiddleware(ads.RemoveAd))

	mux.HandleFunc("GET    /api/files/{fileId}", ads.GetImage)

	mux.HandleFunc("GET    /api/show/{token}", ads.ShowAd)

	wrappedServer := CORSMiddleware(mux)

	log.Printf("The server is listening on port :%s\n", port)

	return &http.Server{Addr: fmt.Sprintf(":%s", port), Handler: wrappedServer}
}
