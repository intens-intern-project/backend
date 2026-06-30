package main

import (
	"database/sql"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Counter struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func main() {
	dsn := "host=db port=5432 user=admin dbname=admin password=admin sslmode=disable"

	var err error
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost"},
		AllowMethods:     []string{"GET", "USER", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.GET("/counter", getCounter)
	router.PUT("/counter/plus", incrementCounter)
	router.PUT("/counter/reset", resetCounter)

	router.Run(":8080")
}

func getCounter(c *gin.Context) {
	var counter Counter
	err := db.QueryRow("SELECT id, name, value FROM counters WHERE name = $1", "default").Scan(&counter.ID, &counter.Name, &counter.Value)

	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, counter)
}

func incrementCounter(c *gin.Context) {
	_, err := db.Exec("UPDATE counters SET value = value + 1 WHERE name = $1", "default")
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	getCounter(c)
}

func resetCounter(c *gin.Context) {
	_, err := db.Exec("UPDATE counters SET value = 0 WHERE name = $1", "default")
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	getCounter(c)
}
