# Shire-Client

## TODO

- [x] info server
  - [x] API
  - [x] database
- [x] project overall structure
  - [x] core package
  - [x] diagram
- [x] one to one message
- [x] group message 
  - [x] start raft
  - [x] watch message change
- [x] add persistence
- [x] frontend improvement
  - [x] search user
  - [x] fetch data from App
  - [x] group create
  - [x] UI improvement
- [ ] configuration panel
- [x] CI/CD
- [ ] testing
- [ ] documentation

## Start development

### Run in development mode 

```shell
wails dev
```
### Build Project

```shell
wails build
```

## Implementation Detail

### Project Structure
TODO

### Message Transmit

All message transmit is handled by the package `github.com/templari/shire-client/core`. The message handling process is mainly done by `Core.SendMessage` and `Core.ReceiveMessage`, `Core.Subscribe` is provided for subscribing to all messages goes through the package.

#### One to one

Assume a user called Alice want to send message to Bob.

1. Alice will fetch user list from info server.
    ```
    GET {{info_server}}/users
    ```

    The response will contains a list of users.

    ```
    [
      {
        "id": 1,
        "name": "Alice",
        "address": "192.168.1.52",
        "port": 51463,
        "rpcPort": 51464,
        "createdAt": "2023-04-21T14:10:06.235Z",
        "updatedAt": "2023-04-21T14:10:06.256Z"
      },
      // ...
    ]
    ```

2. Alice will then make a http post request to Bob
   ```
   POST {{user.address}}:{{user.port}}/message

   {
     "from": {{alice.id}},
     "to": {{bob.id}},
     "content": "Hello Bob"
   }
   ```

3. After Bob received the request, Bob will then response to Alice with a OK status code. Then both Alice and Bob will pass the message to all subscribers.
   

#### Group

Assume a user called Alice want to have a group another 2 members, Bob and Charlie. 

1. Alice will fetch user list from info server. Detail see [One to one](#one-to-one)


2. Alice will send a request to info server with members' id which Alice wish to have in the group.

    ```
    POST {{info_server}}/group

    [
      {
        "userId": {{alice.id}},
      },
      {
        "userId": {{bob.id}},
      },
      {
        "userId": {{charlie.id}},
      }
    ]
    ```
   The response will contains a unique group id.
   ```
   {
     "groupUsers": [
     //users...
     ],
     "id": {{groupId}},
     "createdAt": "2023-04-22T04:27:35.001Z",
     "updatedAt": "2023-04-22T04:27:35.001Z"
   }
   ```

3. With group created, each member will start start their raft and then start transimit group message.
    
    3.1 Alice will make http post request to each members of the group to ask them to prepare their raft and set up RPC service

    ```
    POST {{user.address}}:{{user.port}}/group

    {
      "groupId": {{groupId}}
    }
    ```

    3.2 When a member received the request, say Bob, Bob will then register and serve RPC service.
      
    ```go
     // todo 
    ```

    3.3 After Bob has started the service and begin to listen. Bob will make a http put request to info server to update the group info. Then bob will response to Alice's call with a success message.

    ```
    PUT {{info_server}}/groups/{{groupId}}/users/{{userId}} 

     {
         "userPort": {{port}}
     }
    ```

    3.4 When Alice received success message from all members, Alice will and send a message to each member to start the raft.


    ```
    POST {{user.address}}:{{user.port}}/group/start

    {
      "groupId": {{groupId}}
    }
    ```

    3.5 When a member received the request, say Bob, Bob will then start the raft. As for the organizer of the group, Alice will start the raft when 3.4 is done.

    ```go

    // todo
    ```

    Till now, all member of the group has started the raft and is ready to transmit message.
    