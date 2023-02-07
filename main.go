package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// json fields are lowercase so that we can serialize the attributes properly and distinguish between
// a Go attribute and json attribute
type song struct {
	ID       string `json:id`
	Title    string `json:title`
	Artist   string `json:author`
	Quantity int    `json:quantity`
}

var playlist = []song{
	{ID: "1", Title: "Everything", Artist: "Lil Baby", Quantity: 8},
	{ID: "2", Title: "0 to 100", Artist: "Drake", Quantity: 7},
	{ID: "3", Title: "Prague", Artist: "SL", Quantity: 2},
}

// handling route of getting every playlist in library in json form
func getSongs(c *gin.Context) {
	//this will give us nicely formated json format of the playlists slice. Able to serialize the playlists slice
	c.IndentedJSON(http.StatusOK, playlist)
}

func songByID(c *gin.Context) {
	id := c.Param("id")
	song, err := getsongByID(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "playlist not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, song)
}

func checkoutSong(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"Message": "Missing id query parameter"})
		return
	}
	song, err := getsongByID(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "playlist not found"})
		return
	}

	if song.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "playlist not available"})
		return
	}

	song.Quantity -= 1
	c.IndentedJSON(http.StatusOK, song)
}
func checkInSong(c *gin.Context) {
	id, ok := c.GetQuery("id")
	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"Message": "Missing id query parameter"})
		return
	}
	song, err := getsongByID(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "playlist not found"})
		return
	}

	song.Quantity += 1
	c.IndentedJSON(http.StatusOK, song)
}

func getsongByID(id string) (*song, error) {
	for i, b := range playlist {
		if b.ID == id {
			return &playlist[i], nil
		}
	}

	return nil, errors.New("this song is not in playlist")
}

func createSong(c *gin.Context) {
	var newSong song

	//error response message
	if err := c.BindJSON(&newSong); err != nil {
		return
	}

	playlist = append(playlist, newSong)
	c.IndentedJSON(http.StatusCreated, newSong)
}
func main() {
	//setting up router for handling different routes
	router := gin.Default()
	//if this route is prompted, the getSongs() func will be called. Will return in json form
	router.GET("/playlists", getSongs)
	router.POST("/playlists", createSong)
	router.GET("/playlists/:id", songByID)
	router.PATCH("/checkout", checkoutSong)
	router.PATCH("checkin", checkoutSong)
	router.Run("localhost:3000")
}
