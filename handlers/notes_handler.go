package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Note struct {
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type NoteDetails struct {
	Id      int    `json:"id"`
	Title   string `json:"title" `
	Content string `json:"content" `
}

// CreateNotes - Add new note
func CreateNotes(c *gin.Context, db *sql.DB) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, please login"})
		return
	}

	var note Note
	if err := c.ShouldBindJSON(&note); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input, all fields required"})
		return
	}

	query := `INSERT INTO notes(title, content, user_id) VALUES($1, $2, $3)`
	_, err := db.Exec(query, note.Title, note.Content, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add note"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Note added successfully"})
}

// DisplayNotes - Get all notes for user
func DisplayNotes(c *gin.Context, db *sql.DB) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, please login"})
		return
	}

	query := `SELECT id, title, content FROM notes WHERE user_id = $1`
	rows, err := db.Query(query, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
		return
	}
	defer rows.Close()

	var notes []NoteDetails
	for rows.Next() {
		var n NoteDetails
		if err := rows.Scan(&n.Id, &n.Title, &n.Content); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read note"})
			return
		}
		notes = append(notes, n)
	}

	// Return empty array instead of 204 (frontend-friendly)
	c.JSON(http.StatusOK, gin.H{"notes": notes})
}

func DeleteNotes(c *gin.Context,db *sql.DB){
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, please login"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
    	return
	}
	query := `delete from notes where id = $1 and user_id = $2`
	_,err = db.Exec(query,id,userID)
	if err != nil {
		c.JSON(http.StatusNoContent,gin.H{"err":"No content to delete"})
	}
	c.JSON(http.StatusOK,gin.H{"Success":"Successfully Deleted the task"})
}

func UpdateNotes(c *gin.Context,db *sql.DB){
	type Input struct{
		Title string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	var input Input
	if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
	

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
    	return
	}
	query := `update notes set title = $1 , content = $2 where id = $3`
	_,err = db.Exec(query,input.Title,input.Content,id)
	if err != nil {
		c.JSON(http.StatusNoContent,gin.H{"err":"No task to update"})
	}
	c.JSON(http.StatusOK,gin.H{"Success":"Successfully Updated "})
}