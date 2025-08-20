package handlers

import (

	"github.com/gin-contrib/sessions"
	//"github.com/gin-contrib/sessions/cookie"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Signup(c *gin.Context, db *sql.DB) {
	fmt.Println(">>>>>>>>>", db)
	var user Users
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "All fields are required"})
		return
	}
	hashed_password, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Failed": err})
		return
	}
	query := `INSERT INTO users(username,password) VALUES($1,$2)`
	_, err = db.Exec(query, user.Username, string(hashed_password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{"message": user.Username + " signed up successfully"})
}
func Signin(c *gin.Context, db *sql.DB) {
	session := sessions.Default(c)
	
	var val_user Users
	if err := c.ShouldBindJSON(&val_user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "All fields are required"})
		return
	}
	var userID int
	
	var stored_hashedPass string
	query := `select id,password from users where username = $1`
	err := db.QueryRow(query, val_user.Username).Scan(&userID,&stored_hashedPass)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "Invalid Credentials"})
		return
	}
	session.Set("user_id", userID)
	session.Save()
	
	if err := bcrypt.CompareHashAndPassword([]byte(stored_hashedPass), []byte(val_user.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "Invalid Credentials"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Success": "Login Successful"})
}
