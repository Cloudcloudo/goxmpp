

# CHECK USER JID

**GET** `/users/{ jid_localpart }`
### RESPONSE

200 - OK - user exists

404 - Not found

# REGISTER

**POST** `/users`

    `{"jid":"szymon","password":"123456"}`


# LOGIN

**POST** `/login`

    `{"jid":"szymon","password":"123456"}`

### RESPONSE
200 - OK
```
{
  "token": "HDjaPKtyfSbzBMQwsvBbhzMnGAkQmJZhraKxjHxCHFliPeab",
  "jid": "szymon@localhost/",
  "full_name": null,
  "nick_name": null,
  "company": null,
  "department": null,
  "position": null,
  "role": null,
  "street": null,
  "streer_2": null,
  "city": null,
  "state": null,
  "zip_code": null,
  "country": null,
  "about": null,
  "last_seen": "2018-01-27T11:11:11.508393Z",
  "presence": "",
  "status": "",
  "email": null,
  "www": null,
  "phone": null,
  "birthday": null,
  "avatar": null
}
```

401 - Unauthorized


# STREAM
`ws:// .... /stream?token={ TOKEN }`

**BASE DATA FRAME**

```
{
  "type": "TYPE",
  "payload": {}
}
```

## TYPES
* ## message
    ```
    {
      "type": "message",
      "payload": {
        "name": "message",
        "from": "szymon@localhost",
        "id": "as12da11",
        "to": "alicja@localhost",
        "type": "chat",
        "body": "Test wiadomości z Webscoketów"
      }
    }
    ```
    ##### Typing...
    ```
    {
      "type": "message",
      "payload": {
        "from": "alicja@localhost/dom",
        "id": "purplee0d70390",
        "to": "szymon@localhost",
        "type": "chat",
        "body": "",
        "composing": {
          "ID": ""
        }
      }
    }
    ```
    ##### Paused...
    ```
    {
      "type": "message",
      "payload": {
        "from": "alicja@localhost/dom",
        "id": "purplee0d70391",
        "to": "szymon@localhost",
        "type": "chat",
        "body": "",
        "paused": {}
      }
    }
    ```
    ##### Active...
    ```
    {
      "type": "message",
      "payload": {
        "from": "alicja@localhost/dom",
        "id": "purplee0d70391",
        "to": "szymon@localhost",
        "type": "chat",
        "body": "",
        "active": {}
      }
    }
    ```
    
* ## presence
    ```
    {
      "type": "presence",
      "payload": {
        "from": "alicja@localhost",
        "show": "away",
        "status": "Mam  na imię alicja"
      }
    }
    ```
    
    ### Status types
    
    * if show is missing or set to `av` then contact is available
    * away -- away mode
    * chat  -- ready to chat
    * dnd   -- do not disturb
    * xa    -- i will be later (extended away)
    
* ##iq

    ### get contact list
    ```
    {
      "type": "iq",
      "payload": {
        "type": "get",
        "id": "12345",
        "query": {"jabber:iq:roster"}
      }
    }
    ```
    
