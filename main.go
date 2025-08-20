package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"database/sql"
	"log"

	"notes/handlers"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {

	constr := "user=postgres password=Fyers@2025 dbname=user_notes sslmode = disable"
	db, err := sql.Open("postgres", constr)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("DB connection failed ", err)
	}
	log.Println("Connected to DB")

	r := gin.Default()

	store := cookie.NewStore([]byte("secret-key")) // encryption key
	r.Use(sessions.Sessions("mysession", store))

	r.POST("/signup", func(c *gin.Context) {
		handlers.Signup(c, db)
	})
	r.POST("/signin", func(c *gin.Context) {
		handlers.Signin(c, db)
	})

	r.POST("/notes", func(c *gin.Context) {
		handlers.CreateNotes(c, db)
	})

	r.PUT("/update/:id",func(c *gin.Context){
		handlers.UpdateNotes(c,db)
	})

	r.GET("/notes", func(c *gin.Context) {
		handlers.DisplayNotes(c, db)
	})

	r.DELETE("/users/:id", func(c *gin.Context){
		handlers.DeleteNotes(c,db)
	})
	r.Run(":8080")
}
