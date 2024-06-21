# chat Application with Golang
***
### Description
This repo contains and is about a microservice base chat application with golang.

This chat application is a **Room** base chat app,
users can create rooms and join the rooms,
for now the chat is a message base means users by sending 
messages in the room and all the users in room will receive 
the message. chat app is implemented with websocket to be 
real-time.

Beside the chat app there is another service called notification
service that will connect to user with websocket.

Used mongoDB for database, Redis for cache, RabbitMQ for message broker.
***
### Diagram
![chat_diagram.svg](docs/chat_diagram.svg)

### How to start the app

1. First clone the repo
```shell
    git clone https://github.com/arminshfatemi/chat_application.git
```

2. Change directory to the project directory
```shell
    cd chat_application
```

3. Start the docker compose
```shell
    sudo docker compose up
```
Now all the services are up and running
***

## How the application works 
### Chat Service:
Chat service is responsible for login, signup of users and
creating, joining and sending messages to the users in
real time to users.

This service use REST APIs to create rooms, login and signup of user.

And it use websocket for realtime messages sending of messages

### Notification Service
Notification Service is responsible for sending notifications to 
users in real time with websocket

This Service have a consumer that use RabbitMQ Queues to get notifications from Chat service

### Scheduler Service
Scheduler Service will do two main task, with the given config
it will do it every given time, Ep: every minute or every hour

first task is get the old messages and log them
second task is get the old messages and put them in the archive collection

### RabbitMQ 
RabbitMQ as a message broker is going to be the way that chat 
service sends events to notification service

### MongoDB
Our database to save messages, notifications, clients and rooms

### Redis
Cache database for caching the recent messages of the rooms 
to increase the response time and load
***

## Main urls

### 1. Signup
**URL**: 127.0.0.1:8000/api/user/signup/ 

**method**: POST

**Example Json**:
```
{
"username": "username",
"email": "test@test.com",
"password": "password"
}
```
### 2. login
**URL**: 127.0.0.1:8000/api/user/login/

**method**: POST

**Example Json**:
```
{
"username": "username",
"password": "password"
}
```
### 3. Room Creation
**URL**: 127.0.0.1:8000/api/room/create/

**Auth**:
for being able to use this endpoint you need to put your login JWT in the header
example of the token :

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTg5MjU2MjMsImlkIjoiNjY3MWE1MmJjY2EyMGUzMThlNjE3ZGQ0IiwidXNlcm5hbWUiOiJ2YXJnaGEifQ.BeRG8eer-fchKmjCcbHqi0edZ7IexIyoJ7XGTaMCuJ8
```

**method**: POST

**Example Json**:
```
{
"name": "room_name",
}
```
**Description**:
this url is for creating new room 

### 4. Joining a Room
**URL**: 127.0.0.1:8000/ws/join/?name=

**Auth**:
for being able to use this endpoint you need to put your login JWT in the header
example of the token :

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTg5MjU2MjMsImlkIjoiNjY3MWE1MmJjY2EyMGUzMThlNjE3ZGQ0IiwidXNlcm5hbWUiOiJ2YXJnaGEifQ.BeRG8eer-fchKmjCcbHqi0edZ7IexIyoJ7XGTaMCuJ8
```

**Example Json**:
```
{
    "type": "message",
    "content": "hello"
}
```
**Description**:
this url is for joining th room with websocket.

given Json example is for sending a new message in the room you can send message when you are connected.

you need to specify the name of the room you want to join in url example:
127.0.0.1:8000/ws/join/?name=test_room

### 5. Notification websocket 
**URL**: 127.0.0.1:8000/ws/join-notification/

**Auth**:
for being able to use this endpoint you need to put your login JWT in the header
example of the token :

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTg5MjU2MjMsImlkIjoiNjY3MWE1MmJjY2EyMGUzMThlNjE3ZGQ0IiwidXNlcm5hbWUiOiJ2YXJnaGEifQ.BeRG8eer-fchKmjCcbHqi0edZ7IexIyoJ7XGTaMCuJ8
```
**Description**:
this endpoint is a websocket, after joining, if the room you are joined 
get a message it will send a notification in the channel to notify the user

