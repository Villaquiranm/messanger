# Messanger

A chat simple chat system based on gRPC communication.

### Requirements

- Golang 1.15

# Architecture

The Client and Server executables listen to `SIGINT` to exit the execution

### Server

Server has two main folders, the first one is database when we admin all the stored data. We use a key-value database called BoltDB.

- database: we deal with all transactions, like sent messages, logged users etc.
- grpc: Contains the gRPC server initialization, configuration and callbacks registrations.
- grpc/services: Contains a class `ChatService` that implements the gRPC methods:

```
  rpc Join(Options) returns (LoggedUser);
  rpc GetMessages(LoggedUser) returns (Messages);
  rpc Send(Message) returns (Empty);
  rpc Disconnect(LoggedUser) returns (Empty);
```

### Client

Client also have two main folders:

**stdin**

The class that takes the input from the standard input. Contains a method `Read()` that blocks the execution until a new message is received. The received message is returned by this function.

**grpc**

Where the gRPC client is initialized. We first ask the user to input an username. _If the user already exists in the chat, the system will ask for another name_. When the user is successfully logged in the server returns a new `LoggedUser` object.

Using this object we initilize out two IO classes:

- Receiver: A package that will retrieve all non read messages and print them in the standard output. Note that we created the interface `Receiver` in the actions folder because in the future we could execute several actions when receiving a message. For exemple: write to a file, send a mail, etc.

- Sender: Using the `stdin.Reader` class we sent a new message to the server each time the `Read()` function returns a value.

## Features

- Admin send events in the chat that are received for all the users: `User x join the chat` or `user x left the chat`.
- Client leaves the chat automatically without user action. When user sends a `SIGINT` signal, the system will automatically closes the connexion.

## Build and execution

### client

```
cd client
go build .
chmod +x client
./client

```

### server

```
cd server
go build .
chmod +x server
./server

```

## TODO

- [ ] Unit tests and integration tests
- [ ] Add new actions on the folder `client/grpc/receiver/actions`
- [ ] Test different messages & usernames sizes
- [ ] Add some sort of commands like `/info` that returns the number of users in the chat, etc.
