# Gocc

Gocc is a peer to peer terminal chat application. This provides encrypted
communication between clients without having to go through a central server
of any kind. Additionally, people can disconnect and connect to any peers
without killing the chat room.

### Why Use Gocc?

Gocc is for users who care about secure and untracked communication. Data
is never sent to remote servers, so as long as you trust everyone in the chat
room to not save data, then no data is tracked.

## Installation

> go get github.com/jpanda109/gocc

## Usage

(assuming go's bin is in your path) run for various flags
```
gocc -h
```

#### Flags
--port, -p: This is the port that other peers will be able to connect to in order
to join a chat room  
--connect, -c: Specifies the address of the peer you wish to connect to  
--name, -n: Specifies the name that you'll be recognized by in the chatroom

#### Example
User 1:
> gocc -p 8080 -n user1

User 2:
> gocc -p 7000 -n user2 -c <remote_ip>:8080
