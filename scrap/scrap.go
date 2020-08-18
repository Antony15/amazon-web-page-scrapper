package scrap

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/Antony15/amazon-web-page-scrapper/logger"
	"github.com/Antony15/amazon-web-page-scrapper/utils"
	"github.com/gocolly/colly"
)

var db *mgo.Database

type scrap struct {
	Url string `json:"url"`
}

type scrapped struct {
	Url      string  `json:"url"`
	Product  product `json:"product"`
	Datetime string  `json:"datetime"`
}

type product struct {
	Name         string `json:"name"`
	ImageURL     string `json:"imageURL"`
	Stars        string `json:"stars"`
	Description  string `json:"description"`
	Price        string `json:"price"`
	TotalReviews string `json:"totalReviews"`
}

func init() {
	session, err := mgo.Dial("localhost/sellerapp_golang_mongodb")
	if err != nil {
		logger.Log.Println("Error : ", err.Error())
	}

	db = session.DB("sellerapp_golang_mongodb")
}

func home(w http.ResponseWriter, r *http.Request) {
	html := `<html>
		<head>
			<title>SellerApp Golang Test</title>
		</head>
		<body>
			<h1>SellerApp Golang Test</h1>
		</body>
	</html>`
	w.Write([]byte(html))
}

func collection() *mgo.Collection {
	return db.C("pages")
}

func scrapurl(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var (
		request  scrap
		scrapped scrapped
	)
	ResponseArray := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		ResponseArray["message"] = "Request Error"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseArray)
		return
	} else {
		scrapped.Url = request.Url
		c := colly.NewCollector()
		c.OnHTML("body", func(e *colly.HTMLElement) {
			scrapped.Product.Name = e.ChildText("span#productTitle.a-size-large.product-title-word-break")
			scrapped.Product.ImageURL = e.ChildAttr("#landingImage", "src")
			scrapped.Product.Stars = e.ChildText("span.a-icon-alt")
			utils.FormatStars(&scrapped.Product.Stars)
			scrapped.Product.Description = e.ChildText("span.a-list-item")
			scrapped.Product.TotalReviews = e.ChildText("span#acrCustomerReviewText.a-size-base")
			utils.FormatReviews(&scrapped.Product.TotalReviews)
			scrapped.Product.Price = e.ChildText("span#priceblock_ourprice.a-size-medium.a-color-price.priceBlockBuyingPriceString")
			utils.FormatPrice(&scrapped.Product.Price)
		})
		c.Visit(scrapped.Url)
		scrapped.Datetime = time.Now().Format("2006-01-02 15:04:05")
		jsonStr, err := json.Marshal(scrapped)
		checkErr(err)
		req, err := http.NewRequest("POST", "http://localhost:9999/writedocument", bytes.NewBuffer(jsonStr))
		checkErr(err)
		res, err := http.DefaultClient.Do(req)
		checkErr(err)
		defer res.Body.Close()
		if res.StatusCode != 200 {
			ResponseArray["message"] = "Request Failed to write to document"
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ResponseArray)
			return
		}
		ResponseArray["message"] = "Requested URl Successfully scrapped & saved in database"
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ResponseArray)
		return
	}
}

func writedocument(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var request scrapped
	ResponseArray := make(map[string]interface{})
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		ResponseArray["message"] = "Request Error"
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ResponseArray)
		return
	} else {
		request.Datetime = time.Now().Format("2006-01-02 15:04:05")
		if err := save(request); err != nil {
			ResponseArray["message"] = "Request Failed to write to document"
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(ResponseArray)
			return
		}
		ResponseArray["message"] = "Request Wrote to document"
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(ResponseArray)
		return
	}
}

// Save inserts or updates an scrapped to the database.
func save(scrapped scrapped) error {
	_, err := collection().Upsert(
		bson.M{"url": scrapped.Url},
		bson.M{"$set": bson.M{"product": scrapped.Product, "datetime": scrapped.Datetime}},
	)
	return err
}

func checkErr(err error) {
	if err != nil {
		logger.Log.Println("Error : ", err.Error())
	}
}

func SetupHandlers() {
	log.Println("Server started at port :9999")
	logger.Log.Println("Server started at :9999")
	r := mux.NewRouter()
	r.HandleFunc("/", home)
	r.HandleFunc("/scrapurl", scrapurl).Methods("POST")
	r.HandleFunc("/writedocument", writedocument).Methods("POST")
	if err := http.ListenAndServe(":9999", r); err != nil {
		panic(err)
	}
}
