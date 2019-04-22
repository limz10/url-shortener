# URL Shortener

## Features
* Generates short URLs using only `[a-z0-9]` characters.
* Doesn’t create multiple short URLs when you try to shorten the same URL. In this case, the script will simply return the existing short URL for that long URL.

## Install
Docker and Docker Compose will be required to install this software.

Unzip the files, and run the following terminal command:

```
sudo docker-compose up --build
```

Your app will be running at `localhost:8080`

## Usage
### Redirect
For instance, go to [localhost:8080/cogo](localhost:8080) 
It will direct you to the homepage of Cogo Labs

### Generate short URL token
Send a POST request to `localhost:8080/s`

For instance, [localhost:8080/s?url=https://apple.com](localhost:8080/s?url=https://apple.com) will generate a short URL token.  
You may save this token and go to `localhost:8080/THIS_TOKEN`.

It will direct you to the official page of Apple.

### Hitting already shortened URLs
This software keeps unique URL tokens.  If you send the same POST request many times, this software won't generate different tokens.

