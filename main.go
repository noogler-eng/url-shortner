package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type URL struct {
	Id          string    `json:"id"`
	OriginalUrl string    `json:"orgurl"`
	ShortUrl    string    `json:"shorturl"`
	CreatedAt   time.Time `json:"createdat"`
}

// map contains key and value
var urlDb = make(map[string]URL)

func generateShortUrl(originalURL string) string {
	// 1. hasing the incomming string
	hasher := md5.New()
	_, err := hasher.Write([]byte(originalURL))
	if err != nil {
		panic(err)
	}
	// fmt.Printf("type:%T\n", hasher)
	// fmt.Println("value: ", value)
	// fmt.Println("hasher: ", hasher)

	// 2. hashing returns the array of distributed bytes
	// summation of bytes
	data := hasher.Sum(nil)
	// fmt.Println("data after sum: ", data)
	// fmt.Printf("data:%T\n", data)

	// 3. encoding the bytes into an string using hex
	encodedData := hex.EncodeToString(data)

	// 4. retuning string of 8 charaters
	return string(encodedData[:8])
}


func createURL(originalURL string) string {
	// 1. calling short url and storing database with same id as of shortURL
	shortURL := generateShortUrl(originalURL)
	id := shortURL
	urlDb[id] = URL{
		Id:          id,
		OriginalUrl: originalURL,
		ShortUrl:    shortURL,
		CreatedAt:   time.Now(),
	}
	return shortURL
}


func getURL(origninalURL string) (URL, error) {
	shortURL := createURL(origninalURL)
	url, ok := urlDb[shortURL]
	if !ok {
		return URL{}, errors.New("URL NOT FOUND")
	}

	return url, nil
}


func main() {
	fmt.Println("url shortner")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("GET method", r.Host)
		// fmt.Println(r.URL.Path)
		str := r.URL.Path[1:]
		// fmt.Println("getting path", str)
		originalURLFromMap := urlDb[str].OriginalUrl
		// fmt.Println("getting original url from map: ",originalURLFromMap);
		http.Redirect(w, r, originalURLFromMap, http.StatusSeeOther)



		// Fprintf standardly write on that writer which we provides to it
		// _, err := fmt.Fprintf(w, "hello sharad")
		// if err != nil {
		// 	panic(err)
		// }
	})

	http.HandleFunc("/getShortURL", func(w http.ResponseWriter, r *http.Request) {
		type Data struct {
			Url string 	`json:"url"`
		}
		var data Data;
		// json to struct
		err := json.NewDecoder(r.Body).Decode(&data);
		if err!=nil {
			http.Error(w, "Invalid Body", http.StatusBadRequest)
		}

		originalURL := data.Url
		urlObj, err_1 := getURL(originalURL)
		// fmt.Println(urlObj)
		if err_1!=nil {
			panic(err_1)
		}

		new_shortURL := strings.Join([]string{"http://" + r.Host, urlObj.ShortUrl}, "/")
		var sentData = Data{
			Url: new_shortURL,
		}
		w.Header().Set("content-type", "application/json")
		// struct to json
		json.NewEncoder(w).Encode(sentData)
	})


	// starting HTTP server at port :8080
	fmt.Println("starting server")
	err := http.
		ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}


// handler function manages all itself GET, POST, ..
// it send the request on basic of data we send