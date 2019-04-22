package main

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
    "github.com/gorilla/mux"
    "github.com/spf13/viper"
    "log"
    "math/rand"
    "net/http"
)

var db *sql.DB

func main() {
    // Instantiate the configuration
    viper.SetConfigName("config")
    viper.AddConfigPath(".")
    viper.ReadInConfig()

    // Instantiate the database
    var err error
    dsn := viper.GetString("mysql_user") + ":" + viper.GetString("mysql_password") + "@tcp(" + viper.GetString("mysql_host") + ":3306)/" + viper.GetString("mysql_database") + "?collation=utf8mb4_unicode_ci&parseTime=true"
    db, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Instantiate the mux router and assign mux as the HTTP handler
    r := mux.NewRouter()
    r.HandleFunc("/s", ShortenHandler).Queries("url", "")
    r.HandleFunc("/{token:[a-zA-Z0-9]+}", ShortenedUrlHandler)
    r.HandleFunc("/", CatchAllHandler)
    http.Handle("/", r)
    http.ListenAndServe(":8080", nil)
}

// Shortens a given URL passed through in the request.
// If the URL has already been shortened, returns the existing URL.
// Writes the short URL in plain text to w.
func ShortenHandler(w http.ResponseWriter, r *http.Request) {
    // Check if the url parameter has been sent along (and is not empty)
    url := r.URL.Query().Get("url")

    if url == "" {
        http.Error(w, "", http.StatusBadRequest)
        return
    }

    // Get the short URL out of the config
    if !viper.IsSet("short_url") {
        http.Error(w, "", http.StatusInternalServerError)
        return
    }
    short_url := viper.GetString("short_url")

    // Check if url already exists in the database
    var token string
    err := db.QueryRow("SELECT `token` FROM `redirect` WHERE `url` = ?", url).Scan(&token)
    if err == nil {
        // The URL already exists! Return the shortened URL.
        w.Write([]byte(short_url + "/" + token))
        return
    }

    // generate a token and validate it doesn't
    // exist until we find a valid one.
    var exists = true
    for exists == true {
        token = generateToken()
        err, exists = tokenExists(token)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    }

    // Insert it into the database
    stmt, err := db.Prepare("INSERT INTO `redirect` (`token`, `url`) VALUES (?, ?)")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    _, err = stmt.Exec(token, url)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Write([]byte(short_url + "/" + token))
}

// generateToken will generate a random token to be used as shorten link.
func generateToken() string {
    // It doesn't exist! Generate a new token for it
    var chars = []rune("0123456789abcdefghijklmnopqrstuvwxyz")
    s := make([]rune, 6)
    for i := range s {
        s[i] = chars[rand.Intn(len(chars))]
    }

    return string(s)
}

// tokenExists will check whether the token already exists in the database
func tokenExists(token string) (e error, exists bool) {
    err := db.QueryRow("SELECT EXISTS(SELECT * FROM `redirect` WHERE `token` = ?)", token).Scan(&exists)
    if err != nil {
        return err, false
    }

    return nil, exists
}

// Handles a requested short URL.
// Redirects with a 301 header if found.
func ShortenedUrlHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Check if a token exists
    vars := mux.Vars(r)
    token, ok := vars["token"]
    if !ok {
        http.Error(w, "", http.StatusBadRequest)
        return
    }

    // 2. Check if the token exists in the database
    var url string
    err := db.QueryRow("SELECT `url` FROM `redirect` WHERE `token` = ?", token).Scan(&url)

    if err != nil {
        http.NotFound(w, r)
        return
    }

    // Redirect the user to the URL.
    http.Redirect(w, r, url, http.StatusMovedPermanently)
}

// Catches all other requests to the short URL domain.
// If a default URL exists in the config redirect to it.
func CatchAllHandler(w http.ResponseWriter, r *http.Request) {
    // 1. Get the redirect URL out of the config
    if !viper.IsSet("default_url") {
        // The reason for using StatusNotFound here instead of StatusInternalServerError
        // is because this is a catch-all function. You could come here via various
        // ways, so showing a StatusNotFound is friendlier than saying there's an
        // error (i.e. the configuration is missing)
        http.NotFound(w, r)
        return
    }

    // 2. If it exists, redirect the user to it
    http.Redirect(w, r, viper.GetString("default_url"), http.StatusMovedPermanently)
}
