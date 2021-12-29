# deck_of_cards_backend
Backend for deck of cards

Backend written in go
Database used - MongoDB
test cases writter in python

## Setup instructions Using Docker


`sudo docker compose up` should create containers for backend as well as mongodb
The server should come up as well

Now the backend will be running at port 8080

Test cases can be run using `python3 test/test.py`

## API Usage

### Defult deck creation without shuffling

`http://localhost:8080/deck?shuffle=false`

### Default deck creation with shuffling
`http://localhost:8080/deck?shuffle=true`

### Deck creation with card codes
`http://localhost:8080/deck?shuffle=true&cards=AH,3H,2H,4S`

### Open a Deck
`http://localhost:8080/deck/{deck_id}`

### Draw two cards
`http://localhost:8080/deck/{deck_id}/cards?count=2`

### Without Docker

Assuming mongodb is already running locally at port 27017
The backend can be built using

`go mod download`
`go build server.go`

It can be run using
`./server localhost`
