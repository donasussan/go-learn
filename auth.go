package handlers

import (
	"fmt"
	"net/http"
	"new/models"
	"strconv"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func handleDatabaseError(c *gin.Context, err error) bool {
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return true
	}
	return false
}
func ShowPage(c *gin.Context) {
	c.HTML(http.StatusOK, "web.html", nil)
}

func ShowLoginForm(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	db, err := models.InitDB()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer db.Close()
	var storedPassword string
	err = db.QueryRow("SELECT password FROM users WHERE username=?", username).Scan(&storedPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if password != storedPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	session := sessions.Default(c)
	session.Set("username", username)
	session.Save()
	c.Redirect(http.StatusSeeOther, "/dashboard")
}

func ShowSignupForm(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", nil)
}
func Dashboard(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	if username != nil {
		c.HTML(http.StatusOK, "web.html", gin.H{
			"username": username,
			"loggedIn": true,
		})
	} else {
		c.Redirect(http.StatusSeeOther, "/login")
	}
}
func Signup(c *gin.Context) {
	var person models.Data
	if err := c.ShouldBind(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db, err := models.InitDB()
	if handleDatabaseError(c, err) {
		return
	}
	defer db.Close()
	var existingUser string
	err = db.QueryRow("SELECT username FROM users WHERE username=?", person.Username).Scan(&existingUser)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}
	_, err = db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", person.Username, person.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var newUser models.User
	err = db.QueryRow("SELECT id, username FROM users WHERE username=?", person.Username).Scan(&newUser.ID, &newUser.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.HTML(http.StatusOK, "login.html", gin.H{
		"ID":       newUser.ID,
		"Username": newUser.Username,
	})

	c.Redirect(http.StatusSeeOther, "/login")

}
func ShowPagination(c *gin.Context) {
	c.HTML(http.StatusOK, "pagination.html", nil)
}
func getalluser(page, pageSize int) ([]models.Person, int, error) {
	offset := (page - 1) * pageSize
	db, err := models.InitDB()
	if err != nil {
		fmt.Println("Error opening database connection:", err)
		return nil, 0, err
	}
	defer db.Close()
	var totalCount int
	err = db.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalCount)
	if err != nil {
		fmt.Println("Error getting total user count:", err)
		return nil, 0, err
	}
	query := "SELECT id, username FROM users LIMIT ? OFFSET ?"
	rows, err := db.Query(query, pageSize, offset)
	if err != nil {
		fmt.Println("Error executing query:", err)
		return nil, 0, err
	}
	defer rows.Close()

	var result []models.Person
	for rows.Next() {
		var data models.Person
		if err := rows.Scan(&data.ID, &data.Username); err != nil {
			fmt.Println("Error scanning row:", err)
			return nil, 0, err
		}
		result = append(result, data)
	}
	return result, totalCount, nil
}

func GetDataHandler(c *gin.Context) {
	page := c.DefaultQuery("page", "1")          // Default to page 1 if not provided
	pageSize := c.DefaultQuery("pageSize", "10") // Default page size to 10 if not provided
	pageInt, _ := strconv.Atoi(page)
	pageSizeInt, _ := strconv.Atoi(pageSize)
	data, count, err := getalluser(pageInt, pageSizeInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := struct {
		Data  []models.Person `json:"data"`
		Count int             `json:"count"`
	}{
		Data:  data,
		Count: count,
	}

	c.JSON(http.StatusOK, response)
}

func DeleteUser(c *gin.Context) {
	userID := c.Query("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	db, err := models.InitDB()
	if handleDatabaseError(c, err) {
		return
	}
	defer db.Close()
	_, err = db.Exec("DELETE FROM users WHERE id=?", userID)

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
