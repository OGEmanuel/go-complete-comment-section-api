package main

import (
	"errors"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Image struct {
	PNG  string `json:"png"`
	WEBP string `json:"webp"`
}

type User struct {
	Image    Image  `json:"image"`
	Username string `json:"username"`
}

type Reply struct {
	ID         string `json:"id"`
	Content    string `json:"content"`
	CreatedAt  string `json:"createdAt"`
	Score      int    `json:"score"`
	ReplyingTo string `json:"replyingTo,omitempty"`
	User       User   `json:"user"`
}

type Comment struct {
	ID        string  `json:"id"`
	Content   string  `json:"content"`
	CreatedAt string  `json:"createdAt"`
	Score     int     `json:"score"`
	User      User    `json:"user"`
	Replies   []Reply `json:"replies"`
}

type CommentSection struct {
	CurrentUser User      `json:"currentUser"`
	Comments    []Comment `json:"comments"`
}

var emptyRely = []Reply{}

var imageAmyrobson = Image{
	PNG:  "https://i.ibb.co/wJ81pfW/image-amyrobson.png",
	WEBP: "https://i.ibb.co/wJ81pfW/image-amyrobson.png",
}

var userAmyrobson = User{
	Image:    imageAmyrobson,
	Username: "amyrobson",
}

var replyMaxblagun = []Reply{
	{
		ID:         "3",
		Content:    "If you're still new, I'd recommend focusing on the fundamentals of HTML, CSS, and JS before considering React. It's very tempting to jump ahead but lay a solid foundation first.",
		CreatedAt:  "1 week ago",
		Score:      4,
		ReplyingTo: "maxblagun",
		User:       userRamsesmiron,
	},
	{
		ID:         "4",
		Content:    "I couldn't agree more with this. Everything moves so fast and it always seems like everyone knows the newest library/framework. But the fundamentals are what stay constant.",
		CreatedAt:  "2 days ago",
		Score:      2,
		ReplyingTo: "ramsesmiron",
		User:       userJuliusomo,
	},
}

var imageMaxblagun = Image{
	PNG:  "https://i.ibb.co/tYLc7Jv/image-maxblagun.png",
	WEBP: "https://i.ibb.co/tYLc7Jv/image-maxblagun.png",
}

var userMaxblagun = User{
	Image:    imageMaxblagun,
	Username: "maxblagun",
}

var imageRamsesmiron = Image{
	PNG:  "https://i.ibb.co/Y28dxbN/image-ramsesmiron.png",
	WEBP: "https://i.ibb.co/Y28dxbN/image-ramsesmiron.png",
}

var userRamsesmiron = User{
	Image:    imageRamsesmiron,
	Username: "ramsesmiron",
}

var imageJuliusomo = Image{
	PNG:  "https://i.ibb.co/3hVx9Cw/image-juliusomo.png",
	WEBP: "https://i.ibb.co/3hVx9Cw/image-juliusomo.png",
}

var userJuliusomo = User{
	Image:    imageJuliusomo,
	Username: "juliusomo",
}

var commentsList = []Comment{
	{
		ID:        "1",
		Content:   "Impressive! Though it seems the drag feature could be improved. But overall it looks incredible. You've nailed the design and the responsiveness at various breakpoints works really well.",
		CreatedAt: "1 month ago",
		Score:     12,
		User:      userAmyrobson,
		Replies:   emptyRely,
	},
	{
		ID:        "2",
		Content:   "Woah, your project looks awesome! How long have you been coding for? I'm still new, but think I want to dive into React as well soon. Perhaps you can give me an insight on where I can learn React? Thanks!",
		CreatedAt: "2 weeks ago",
		Score:     5,
		User:      userMaxblagun,
		Replies:   replyMaxblagun,
	},
}

func deleteReply(c *gin.Context) {
	id := c.Param("id")
	repId, ok := c.GetQuery("repId")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing Id query parameter."})
		return
	}

	Comment, err := getCommentById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Comment not found"})
		return
	}

	var index int
	found := false

	for i, rep := range Comment.Replies {
		if repId == rep.ID {
			index = i
			found = true
			break
		}
	}

	if !found {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Reply not found"})
		return
	}

	Comment.Replies = append(Comment.Replies[:index], Comment.Replies[index+1:]...)
	c.IndentedJSON(http.StatusOK, Comment.Replies)
}

func deleteComment(c *gin.Context) {
	id := c.Param("id")

	var index int = -1
	for i, com := range commentsList {
		if id == com.ID {
			index = i
			break // Break the loop when the comment is found
		}
	}

	if index == -1 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Comment not found."})
		return
	}

	commentsList = append(commentsList[:index], commentsList[index+1:]...)
	c.IndentedJSON(http.StatusOK, commentsList)
}

func decrementScore(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing Id query parameter."})
		return
	}

	Comment, err := getCommentById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Comment not found"})
	}

	if Comment.Score == 0 {
		return
	}

	Comment.Score -= 1
	c.IndentedJSON(http.StatusOK, Comment)
}

func decrementReplyScore(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing Id query"})

		return
	}

	Reply, err := getReplyById(c, id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Reply not found"})
	}

	if Reply.Score == 0 {
		return
	}

	Reply.Score -= 1
	c.IndentedJSON(http.StatusOK, Reply)
}

func incrementReplyScore(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing Id query"})

		return
	}

	Reply, err := getReplyById(c, id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Reply not found"})
	}

	Reply.Score += 1
	c.IndentedJSON(http.StatusOK, Reply)
}

func incrementScore(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing Id query parameter."})
		return
	}

	Comment, err := getCommentById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Comment not found"})
	}

	Comment.Score += 1
	c.IndentedJSON(http.StatusOK, Comment)
}

func editComment(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Comment not found"})
		return
	}

	Comment, err := getCommentById(id)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Comment not found"})
		return
	}

	if err := c.BindJSON(&Comment); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload"})
		return
	}

	updatedContent := Comment.Content
	Comment.Content = updatedContent

	c.IndentedJSON(http.StatusOK, Comment)
}

func editReply(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Reply not found"})
		return
	}

	Reply, err := getReplyById(c, id)

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Reply not found"})
		return
	}

	if err := c.BindJSON(&Reply); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Invalid request payload"})
	}

	updatedReplyContent := Reply.Content
	Reply.Content = updatedReplyContent

	c.IndentedJSON(http.StatusOK, Reply)
}

func commentById(c *gin.Context) {
	id := c.Param("id")
	Comment, err := getCommentById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Comment not found."})
		return
	}

	c.IndentedJSON(http.StatusOK, Comment)
}

func replyById(c *gin.Context) {
	id := c.Query("id")
	Reply, err := getReplyById(c, id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Reply not found."})
		return
	}

	c.IndentedJSON(http.StatusOK, Reply)
}

func addReply(c *gin.Context) {
	id := c.Param("id")
	Comment, _ := getCommentById(id)

	var newReply Reply

	if err := c.BindJSON(&newReply); err != nil {
		return
	}

	Comment.Replies = append(Comment.Replies, newReply)
	c.IndentedJSON(http.StatusOK, Comment)
}

func getReplyById(c *gin.Context, repId string) (*Reply, error) {
	id := c.Param("id")
	Comment, _ := getCommentById(id)

	for i, rep := range Comment.Replies {
		if rep.ID == repId {
			return &Comment.Replies[i], nil
		}
	}

	return nil, errors.New("Reply not found")
}

func getCommentById(id string) (*Comment, error) {
	for i, com := range commentsList {
		if com.ID == id {
			return &commentsList[i], nil
		}
	}

	return nil, errors.New("Comment not found")
}

func addComment(c *gin.Context) {
	var newComment Comment

	// Unmarshal the request body into newComment
	if err := c.BindJSON(&newComment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Initialize commentsList if it's nil
	if commentsList == nil {
		commentsList = make([]Comment, 0)
	}

	// Append the newComment to commentsList
	commentsList = append(commentsList, newComment)

	c.IndentedJSON(http.StatusCreated, newComment)
}

func getCurrentUser(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, userJuliusomo)
}

func getAllComments(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, commentsList)
}

func main() {
	router := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	router.Use(cors.New(config))

	router.GET("/current-user", getCurrentUser)
	router.GET("/all-comments", getAllComments)
	router.GET("/comment/:id", commentById)
	router.GET("/reply/:id", replyById)
	router.POST("/add-comment", addComment)
	router.PATCH("/add-reply/:id", addReply)
	router.PATCH("/edit-comment", editComment)
	router.PATCH("/edit-reply/:id", editReply)
	router.PATCH("/increment-reply-score/:id", incrementReplyScore)
	router.PATCH("/decrement-reply-score/:id", decrementReplyScore)
	router.PATCH("/increment-score", incrementScore)
	router.PATCH("/decrement-score", decrementScore)
	router.DELETE("/delete-comment/:id", deleteComment)
	router.DELETE("/delete-reply/:id", deleteReply)
	router.Run("localhost:8080")

}
