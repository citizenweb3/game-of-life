Example of game on golang.

# How to run
1. install golang (https://go.dev/)
2. run 
```go run ./main.go```



# How to use

send query to localhost:8000/
or
visit localhost:8000/ui/user and use ui.
To quick start use GenerateRandomSystem.
ps. result of battle and process you could find in console


## Examples
### POST|PATHC|DELETE
example of post|patch|delete request:
```
curl --location --request POST 'localhost:8000/cards/mint' \
--header 'Content-Type: application/json' \
--data-raw '{
    "CardId": "7c1d92fff9aec9aa3d9c38044a99d2cb54f4581ee8606a90f9f9a31934d1b9f9",
    "UserID": "ddd",
    "Hp": 10,
    "Level": 4,
    "Strength": 9,
    "Accuracy": 59
}'
```

### GET
example of GET request:
```
curl --location --request GET 'localhost:8000/system/user?user_id=ddd'
```


## Handles
### System
POST /generate - generate all objects automatic. card has 50% to add to set, user has 50% to be open to battle. 
```
	UserCountFrom int - minimum users
	UserCountTo   int - maximum users

	CardCountFrom int - minimum cards in user collection
	CardCountTo   int - maximum cards in user collection
```

GET /system/user/list- return user list

POST /system/user/add - add new user
```
    UserID     string // optional (for empty value uses random string)
    Volts      uint64
    Amperes    uint64
    Cyberlinks uint64
    Kw         uint64
    Random     bool // true - generate random params, false - use params from request
```

GET /system/user?user_id=<USER_ID> - return user params (volts, Amperes...) 

POST /system/gotothefuture - move currnet time to the future (for check freeze)


### Cards
POST /cards/mint - Generate new card for user
```
    UserID string
```

GET /cards/list?owner=<USER_ID>  return users card

GET /cards/info?card_id=<CARD_ID> return cards parametrs 

POST /cards/transfer  send card to new user

```
    Executor string
    CardID   string
    To       string
```

POST /cards/burn delete card
```
    Executor string
    CardID   string
```

POST /cards/freeze freeze card
```
    Executor string
    CardID1   string
    CardID2   string
```

POST /cards/unfreeze unfreze card (delete freezed cards and create new mixed card). CardID - is any of CardID1, CardID2 from freezed
```
    Executor string
    CardID   string
```

### Set

GET /set/actual?user_id=<USER_ID> - return actual set
  
POST /set/card add new card to set
```
    Executor string
    CardID   string
```

PATCH /set/card change in set CardIDLast to CardIDNew
```
    Executor   string
    CardIDLast string
    CardIDNew  string
```
DELETE /set/card remove card from set
```
    Executor string
    CardID   string
```

POST /set/attribute set users attribute for card *NumInSet* in set
```
    Executor string
    NumInSet uint8
    Hp       uint64
    Level    uint8
    Strength uint64
    Accuracy uint64
```

GET/set/attribute?user_id=<USER_ID> return users attribute for cards


### Battle
POST /battle/start Start battle. Rival must be avalible for ballte
```
    Executor string
    Rival    string
```

GET /battle/isopen?user_id=<USER_ID> return open|close users status for battle

POST /battle/ready set battles open|close status for executor 
```
    Executor string
    Ready    bool
```
