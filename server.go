package main

import (
	"context"
	"encoding/json"
	"github.com/rs/cors"
	"github.com/gorilla/mux"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Card , Deck and Response strutures
type Card struct {
	Value string `json:"value" bson:"value"`
	Suit  string `json:"suit" bson:"suit"`
	Code  string `json:"code" bson:"code"`
}

type Deck struct {
	UUID      string `json:"deck_id" bson:"_id"`
	Shuffled  bool   `json:"shuffled" bson:"shuffled"`
	Remaining int    `json:"remaining" bson:"remaining"`
	Cards     []Card `json:"cards" bson:"cards"`
}

type DeckInfo struct {
	UUID      string `json:"deck_id"`
	Shuffled  bool   `json:"shuffled"`
	Remaining int    `json:"remaining"`
}

type DrawnCards struct {
	Cards []Card `json:"cards" bson:"cards"`
}

var collection *mongo.Collection
var allCodesMap map[string]Card
var allCodes []string
// Initialize default cards
func initializeDefaults() {
	rand.Seed(time.Now().UnixNano())
	values := []string{"ACE", "2", "3", "4", "5", "6", "7",
		"8", "9", "10", "JACK", "QUEEN", "KING"}

	suits := []string{"SPADES", "DIAMONDS", "CLUBS", "HEARTS"}

	allCodesMap = make(map[string]Card, 0)

	allCodes = make([]string, 0)

	for _, value := range values {
		for _, suit := range suits {
			currCode := value[0:1] + suit[0:1]
			allCodes = append(allCodes, currCode)
			allCodesMap[currCode] = Card{Value: value, Suit: suit, Code: currCode}
		}
	}
	log.Println("codes of default cards - ", allCodes)
}

func shuffleCards(cards []Card) {
	for i := range cards {
		j := rand.Intn(i + 1)
		cards[i], cards[j] = cards[j], cards[i]
	}
}

func createDeck(w http.ResponseWriter, r *http.Request) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	v := r.URL.Query()
	var shuffle = false
	shuffle, err := strconv.ParseBool(v.Get("shuffle"))
	if err != nil {
		log.Println("Unable to parse shuffle parameter")
		http.Error(w, "Shuffle parameter not passed properly", 400)
		return
	}
	reqUID := uuid.New().String()

	var inputCards string = v.Get("cards")
	log.Println("input Cards", inputCards)

	var cardsArray []string
	log.Println("cards param - ", inputCards)
	log.Println("All codes  - ", allCodes)
	if inputCards != "" {
		cardsArray = strings.Split(inputCards, ",")
	} else {
		// default
		cardsArray = allCodes
	}
	log.Println(cardsArray)
	var cards []Card
	//card := Card{Value: "ACE", Suit: "SPADE", Code: "AS"}
	for _, cardCode := range cardsArray {
		card, present := allCodesMap[cardCode]
		if !present {
			http.Error(w, "One or more card codes not proper", 400)
			return
		}
		cards = append(cards, card)
	}
	if shuffle {
		shuffleCards(cards)
	}

	dataDeck := Deck{UUID: reqUID, Remaining: len(cards), Shuffled: shuffle, Cards: cards}
	log.Println(dataDeck.Cards)
	deckInfo := DeckInfo{dataDeck.UUID, shuffle, len(cardsArray)}
	bsonData, err := bson.Marshal(dataDeck)
	_, err = collection.InsertOne(ctx, bsonData)
	log.Println("collection insert - ", err)
	params := mux.Vars(r)

	log.Println("CreateDeck", params, shuffle)
	json.NewEncoder(w).Encode(deckInfo)
}

func drawCards(w http.ResponseWriter, r *http.Request) {
	log.Println("Draw Cards")
	params := mux.Vars(r)
	UUID := params["uuid"]
	v := r.URL.Query()

	count, err := strconv.Atoi(v.Get("count"))
	if err != nil {
		log.Println("Unable to parse count parameter")
		http.Error(w, "Count parameter not passed properly", 400)
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var resp Deck
	err = collection.FindOne(ctx, bson.M{"_id": UUID}).Decode(&resp)
	if err != nil {
		http.Error(w, "Unable to find deck id", 404)
		return
	}

	// If count is greater than remaining cards, return all remaining cards
	if count > len(resp.Cards) {
		count = len(resp.Cards)
	}

	log.Println("collection replace one  - ", err)
	cards := DrawnCards{Cards: resp.Cards[:count]}

	resp.Cards = resp.Cards[count:]
	resp.Remaining = len(resp.Cards)
	opts := options.Replace().SetUpsert(true)
	filter := bson.M{"_id": UUID}
	_, err = collection.ReplaceOne(ctx, filter, resp, opts)

	json.NewEncoder(w).Encode(cards)
}

func openDeck(w http.ResponseWriter, r *http.Request) {
	log.Println("Open Deck")
	params := mux.Vars(r)
	UUID := params["uuid"]
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var resp Deck
	err := collection.FindOne(ctx, bson.M{"_id": UUID}).Decode(&resp)
	if err != nil {
		http.Error(w, "Unable to find deck id", 404)
		return
	}
	json.NewEncoder(w).Encode(resp)
}

func main() {

	router := mux.NewRouter()
	
	dbHost := os.Args[1]
	mongoURI := "mongodb://" + dbHost + ":27017"
	log.Println("Connecting to mongo - uri - ", mongoURI)
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	
	// client, err := mongo.NewClient(options.Client().ApplyURI("mongodb:localhost//mongodb"))
	
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	collection = client.Database("deckdb").Collection("decks")

	// Initializing default variables
	initializeDefaults()

	// create a deck
	router.HandleFunc("/deck", createDeck).Methods("GET")
	// draw cards
	router.HandleFunc("/deck/{uuid}/cards", drawCards).Methods("GET")
	//open the deck
	router.HandleFunc("/deck/{uuid}", openDeck).Methods("GET")

	handler := cors.Default().Handler(router)

	log.Fatal(http.ListenAndServe(":8080", handler))
	log.Println("Web Server started. Listening on 0.0.0.0:8080")

}
