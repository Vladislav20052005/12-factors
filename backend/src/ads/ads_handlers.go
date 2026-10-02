package ads

import (
	"ad/postgre"
	"ad/utils"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
)

func GetAds(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		log.Printf("Error in get ads request : user_id field is not found\n")
		return
	}
	log.Printf("Get ads request userId(%s)\n", userId)
	// get all ads from userId
	advertisements, err := postgre.GetAdsByUserId(userId)
	if err != nil {
		http.Error(w, "Ads not found", http.StatusNotFound)
		log.Printf("Error in get ads request userId(%s) : %v\n", userId, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(advertisements)
}

func GetAd(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "invalid request", http.StatusBadRequest)
		log.Printf("Error in get ad request : user_id field is not found\n", userId)
		return
	}
	adId := r.PathValue("id")
	log.Printf("Get ad request userId(%s), adId(%s)\n", userId, adId)
	// get all ad from userId where adId == id
	adv, err := postgre.GetAdById(adId, userId)
	if err != nil {
		http.Error(w, "Ads not found", http.StatusNotFound)
		log.Printf("Error in get ad request userId(%s), adId(%s) : %v\n", userId, adId, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"url":  adv.Url,
		"file": string(adv.File),
	})

}

func CreateAd(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "invalid request", http.StatusBadRequest)
		log.Printf("Error in create ad request : user_id field is not found\n", userId)
		return
	}
	log.Printf("Create ad request userId(%s)\n", userId)

	const maxUploadSize = 10 << 20
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, "File to large or bad request", http.StatusBadRequest)
		log.Printf("Error in create ad request userId(%s) : %v\n", userId, err)
		return
	}

	url := r.PostFormValue("URL")
	file, handler, err := r.FormFile("Image")
	fmt.Println(url, file)
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		log.Printf("Error in create ad request userId(%s) : %v\n", userId, err)
		return
	}
	defer file.Close()

	log.Printf("Got a file %s, %d in create ad request userId(%s)\n", handler.Filename, handler.Size, userId)

	// insert url and file into postgres

	adId, err := utils.GenerateId()
	if err != nil {
		http.Error(w, "Error while generating id", http.StatusInternalServerError)
		log.Printf("Error in create ad request userId(%s) : %v\n", userId, err)
		return
	}

	fileId, err := utils.GenerateId()
	if err != nil {
		http.Error(w, "Error while generating id", http.StatusInternalServerError)
		log.Printf("Error in create ad request userId(%s) : %v\n", userId, err)
		return
	}

	f, _ := io.ReadAll(file)

	err = postgre.InsertAd(adId, userId, fileId, f, url)
	if err != nil {
		http.Error(w, "Error inserting the ad", http.StatusInternalServerError)
		log.Printf("Error in create ad request userId(%s) : %v\n", userId, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id":     adId,
		"status": http.StatusAccepted,
	})
}

func EditAd(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "invalid request", http.StatusBadRequest)
		log.Printf("Error in edit ad request : user_id field is not found\n")
		return
	}
	adId := r.PathValue("id")
	log.Printf("Edit ad request userId(%s), adId(%s)\n", userId, adId)

	const maxUploadSize = 10 << 20
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, "File to large or bad request", http.StatusBadRequest)
		log.Printf("Error in edit ad request userId(%s), adId(%s) : %v\n", userId, adId, err)
		return
	}

	url := r.PostFormValue("URL")
	file, handler, err := r.FormFile("Image")
	fmt.Println(url, file)

	if file == nil {
		err := postgre.EditAdURLOnly(adId, userId, url)
		if err != nil {
			http.Error(w, "Error inserting the ad", http.StatusInternalServerError)
			log.Printf("Error in edit ad request userId(%s), adId(%s) : %v\n", userId, adId, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id":     adId,
			"status": http.StatusAccepted,
		})
		return
	}
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		log.Printf("Error in edit ad request userId(%s), adId(%s) : %v\n", userId, adId, err)
		return
	}
	defer file.Close()

	log.Printf("Got a file %s, %d in create ad request userId(%s)\n", handler.Filename, handler.Size, userId)

	// insert url and file into postgres

	f, _ := io.ReadAll(file)

	err = postgre.EditAdNewFile(adId, userId, f, url)
	if err != nil {
		http.Error(w, "Error inserting the ad", http.StatusInternalServerError)
		log.Printf("Error in edit ad request userId(%s), adId(%s) : %v\n", userId, adId, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"id":     adId,
		"status": http.StatusAccepted,
	})
}

func RemoveAd(w http.ResponseWriter, r *http.Request) {
	userId, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "invalid request", http.StatusBadRequest)
		log.Printf("Error in remove ad request : user_id field is not found\n")
		return
	}

	adId := r.PathValue("id")
	log.Printf("Remove ad request userId(%s), adId(%s)\n", userId, adId)

	// remove adId ad from postgres where userId is correct
	err := postgre.RemoveAdById(adId, userId)
	if err != nil {
		http.Error(w, "The ad does not exist", http.StatusBadRequest)
		log.Printf("Error in remove ad request userId(%s), adId(%s) : %v\n", userId, adId, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{
		"status": http.StatusAccepted,
	})
}

func GetImage(w http.ResponseWriter, r *http.Request) {
	fileId := r.PathValue("fileId")
	log.Printf("Get image request fileId(%s)\n", fileId)
	data, err := postgre.GetImage(fileId)
	if err != nil {
		http.Error(w, "Image is not found", http.StatusNotFound)
		log.Printf("Error in get image request fileId(%s) : %v\n", fileId, err)
		return
	}

	w.Header().Set("Content-Type", "image/png")
	w.Write(data)
}

func ShowAd(w http.ResponseWriter, r *http.Request) {
	token := r.PathValue("token")
	log.Printf("Show ad request token(%s)\n", token)

	// get all ads by token
	advertisements, err := postgre.GetAdsByToken(token)
	if err != nil {
		http.Error(w, "No ads", http.StatusNotFound)
		log.Printf("Error in show ad request token(%s) : %v\n", token, err)
		return
	}

	if len(advertisements) == 0 {
		http.Error(w, "No ads", http.StatusNotFound)
		log.Printf("Error in show ad request token(%s) : %v\n", token, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	i := int(rand.Int31()) % len(advertisements)
	json.NewEncoder(w).Encode(map[string]string{
		"file": advertisements[i].File,
		"url":  advertisements[i].Url,
	})
}
