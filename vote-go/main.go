package main

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	redis "github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func RedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		// Addr:     "localhost:6379",
		Password: "", // no password set
		Addr:     "redis:6379",
		DB:       0, // use default DB
	})

	return rdb
}

var characterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

// RandomString generates a random string of n length
func RandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = characterRunes[rand.Intn(len(characterRunes))]
	}
	return string(b)
}

// NewSHA1Hash generates a new SHA1 hash based on
// a random number of characters.
func NewSHA1Hash(n ...int) string {
	noRandomCharacters := 32

	if len(n) > 0 {
		noRandomCharacters = n[0]
	}

	randString := RandomString(noRandomCharacters)

	hash := sha1.New()
	hash.Write([]byte(randString))
	bs := hash.Sum(nil)

	return fmt.Sprintf("%x", bs)
}

// var vote string

// var voter_id int = 0

func main() {
	os.Setenv("OPTION_A", "Cats")
	os.Setenv("OPTION_B", "Dogs")
	option_a := os.Getenv("OPTION_A")
	option_b := os.Getenv("OPTION_B")

	var vote string

	r := gin.Default()

	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")
	// r.Use(func(c *gin.Context) {

	// 	// if c.Request.Method == "POST" {
	// 	// 	rdb := RedisClient()
	// 	// 	// voter_id++
	// 	// 	voter_id := NewSHA1Hash()
	// 	// 	vote = c.PostForm("vote")
	// 	// 	data := map[string]string{"voter_id": voter_id, "vote": vote}
	// 	// 	payload, _ := json.Marshal(data)

	// 	// 	rdb.RPush(ctx, "votes", payload)
	// 	// }

	// 	c.HTML(http.StatusOK, "index.html", gin.H{
	// 		"option_a": option_a,
	// 		"option_b": option_b,
	// 		"vote":     vote,
	// 	})
	// })

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"option_a": option_a,
			"option_b": option_b,
			"vote":     vote,
		})
	})

	r.POST("/vote", func(c *gin.Context) {
		rdb := RedisClient()
		// voter_id++
		voter_id := NewSHA1Hash()
		vote = c.PostForm("vote")
		data := map[string]string{"voter_id": voter_id, "vote": vote}
		payload, _ := json.Marshal(data)

		rdb.RPush(ctx, "votes", payload)
		// c.SetCookie("vote", vote, 10, "/", c.Request.URL.Hostname(), false, true)
		// location := url.URL{Path: "/api/callback/cookies"}
		// c.Redirect(http.StatusFound, location.RequestURI())
		c.Redirect(http.StatusFound, "/")

	})

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
