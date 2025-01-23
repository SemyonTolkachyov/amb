# Anonymous message board
## Description

This board is a microservice application for anonymous messages.

**Features**
- Sending messages
- Search messages

## Application components
##### Services
- Message - saving messages to the database and notifying other services via nats
- Pusher - managing and maintaining client connections via websockets
- Search - search for messages using elastic search
##### Other
- Nats - server for the NATS Messaging System
- Postgres - server for Postgresql database
- Elastic search - search engine
- Nginx - reverse proxy server


## Commands

**Backend services**
***Docker-compose***
For development mode...
```
docker-compose up --watch
```

For production mode...
```
docker-compose up
```

**Web client**
Client requires [Node.js](https://nodejs.org/) v18.3+ to run.

Install the dependencies.

```sh
cd client
npm i
```
And start the client app
for development mode...
```
npm run dev
```
for production mode...

```sh
npm run build
```
