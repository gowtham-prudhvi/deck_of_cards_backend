import requests
import json

# Test cases for testing the service.

if __name__ == "__main__":
    positiveCases = 0
    positiveSuccess = 0
    # Unshuffled default deck
    positiveCases += 1
    deck_id = None
    resp = requests.get("http://localhost:8080/deck?shuffle=false")
    if resp.status_code == 200:
        jsonResp = resp.json()
        # print(jsonResp)
        if "deck_id" in jsonResp and "shuffled" in jsonResp and "remaining" in jsonResp:
            deck_id = jsonResp["deck_id"]
            shuffled = jsonResp["shuffled"]
            remaining = jsonResp["remaining"]
            if (not shuffled) and remaining == 52:
                positiveSuccess += 1
     
    print("number of positive tests completed - %s - success - %s" % (positiveCases, positiveSuccess))           
    
    #Deck with card codes 
    positiveCases += 1
    resp = requests.get("http://localhost:8080/deck?shuffle=true&cards=AH,3H,2H,4S")
    if resp.status_code == 200:
        jsonResp = resp.json()
        if "deck_id" in jsonResp and "shuffled" in jsonResp and "remaining" in jsonResp:
            deck_id = jsonResp["deck_id"]
            shuffled = jsonResp["shuffled"]
            remaining = jsonResp["remaining"]
            if (shuffled) and remaining == 4:
                positiveSuccess += 1

    print("number of positive tests completed - %s - success - %s" % (positiveCases, positiveSuccess))

    # shuffled default deck
    positiveCases += 1
    resp = requests.get("http://localhost:8080/deck?shuffle=true")
    if resp.status_code == 200:
        jsonResp = resp.json()
        if "deck_id" in jsonResp and "shuffled" in jsonResp and "remaining" in jsonResp:
            deck_id = jsonResp["deck_id"]
            shuffled = jsonResp["shuffled"]
            remaining = jsonResp["remaining"]
            if (shuffled) and remaining == 52:
                positiveSuccess += 1
    
    print("number of positive tests completed - %s - success - %s" % (positiveCases, positiveSuccess))

    # Get the deck from deck_id
    positiveCases += 1
    resp = requests.get("http://localhost:8080/deck/%s" % deck_id)
    if resp.status_code == 200:
        
        jsonResp = resp.json()
        # print(jsonResp)
        if "deck_id" in jsonResp and "shuffled" in jsonResp and "remaining" in jsonResp:
            deck_id = jsonResp["deck_id"]
            shuffled = jsonResp["shuffled"]
            remaining = jsonResp["remaining"]
            if (shuffled) and remaining == 52 and (len(jsonResp["cards"]) == 52):
                positiveSuccess += 1
        
    print("number of positive tests completed - %s - success - %s" % (positiveCases, positiveSuccess))

    # Draw 2 cards from the deck
    positiveCases += 1
    resp = requests.get("http://localhost:8080/deck/%s/cards?count=2" % deck_id)
    if resp.status_code == 200:
        
        jsonResp = resp.json()
        # print(jsonResp)
        if "cards" in jsonResp:
            cards = jsonResp["cards"]
            if len(cards) == 2:
                positiveSuccess += 1
    print("number of positive tests completed - %s - success - %s" % (positiveCases, positiveSuccess))  
    
    # Open the deck - 50 cards should remain as two are drawn          
    positiveCases += 1
    resp = requests.get("http://localhost:8080/deck/%s" % deck_id)
    if resp.status_code == 200:
        
        jsonResp = resp.json()
        # print(jsonResp)
        if "deck_id" in jsonResp and "shuffled" in jsonResp and "remaining" in jsonResp:
            deck_id = jsonResp["deck_id"]
            shuffled = jsonResp["shuffled"]
            remaining = jsonResp["remaining"]
            if (shuffled) and remaining == 50 and (len(jsonResp["cards"]) == 50):
                positiveSuccess += 1
    
    print("number of positive tests completed - %s - success - %s" % (positiveCases, positiveSuccess))
    
    #Draw for more than number of cards in the deck(Only 50 left)
    positiveCases += 1
    resp = requests.get("http://localhost:8080/deck/%s/cards?count=100" % deck_id)
    if resp.status_code == 200:
        
        jsonResp = resp.json()
        # print(jsonResp)
        if "cards" in jsonResp:
            cards = jsonResp["cards"]
            if len(cards) == 50:
                positiveSuccess += 1
    
    print("number of positive tests completed - %s - success - %s" % (positiveCases, positiveSuccess))
    
    # All cards are drawn so no card should remain
    positiveCases += 1
    resp = requests.get("http://localhost:8080/deck/%s" % deck_id)
    if resp.status_code == 200:
        
        jsonResp = resp.json()
        if "deck_id" in jsonResp and "shuffled" in jsonResp and "remaining" in jsonResp:
            deck_id = jsonResp["deck_id"]
            shuffled = jsonResp["shuffled"]
            remaining = jsonResp["remaining"]
            if (shuffled) and remaining == 0 and (len(jsonResp["cards"]) == 0):
                positiveSuccess += 1
    
    print("number of positive tests completed - %s - success - %s" % (positiveCases, positiveSuccess))
    
    
    negativeCases = 0
    negativeSuccess = 0
    
    #Count passed is not integer
    negativeCases += 1
    resp = requests.get("http://localhost:8080/deck/%s/cards?count=fdafd" % deck_id)
    if resp.status_code == 400:
        negativeSuccess += 1
    print("number of negative tests completed - %s - success - %s" % (negativeCases, negativeSuccess))
    
    #UUID is not found with draw cards
    negativeCases += 1
    resp = requests.get("http://localhost:8080/deck/dfadsfjdasffdffa/cards?count=1")
    if resp.status_code == 404:
        negativeSuccess += 1
    print("number of negative tests completed - %s - success - %s" % (negativeCases, negativeSuccess))
    
    #UUID is not found with open deck
    negativeCases += 1
    resp = requests.get("http://localhost:8080/deck/dfadsfjdasffdffa")
    if resp.status_code == 404:
        negativeSuccess += 1
    print("number of negative tests completed - %s - success - %s" % (negativeCases, negativeSuccess))
    
    #shuffle is not true or false
    negativeCases += 1
    resp = requests.get("http://localhost:8080/deck?shuffle=fdadfdaf")
    if resp.status_code == 400:
        negativeSuccess += 1
    print("number of negative tests completed - %s - success - %s" % (negativeCases, negativeSuccess))
    
    #cards passed while creating deck are not proper
    negativeCases += 1
    resp = requests.get("http://localhost:8080/deck?shuffle=true&cards=ddfadasffda")
    if resp.status_code == 400:
        negativeSuccess += 1
    print("number of negative tests completed - %s - success - %s" % (negativeCases, negativeSuccess))