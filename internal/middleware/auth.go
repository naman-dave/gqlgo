package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/naman-dave/gqlgo/internal/jwt"
	models "github.com/naman-dave/gqlgo/internal/model"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		header := c.Request.Header["Authorization"]

		// Allow unauthenticated users in
		if len(header) == 0 || header[0] == "" {
			c.Next()
			return
		}

		//validate jwt token
		tokenStr := header[0]
		mobilenumber, err := jwt.ParseToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{
				"data":    "",
				"error":   "not authenticated",
				"message": "not authenticated",
			})
			c.Abort()
			return
		}

		id, token, err := models.GetUserIdByUsername(mobilenumber)
		if err != nil {
			c.JSON(http.StatusBadRequest, map[string]string{
				"data":    "",
				"error":   "not authenticated",
				"message": "not authenticated",
			})
			c.Abort()
			return
		}

		if token != tokenStr {
			c.JSON(http.StatusBadRequest, map[string]string{
				"data":    "",
				"error":   "not authenticated",
				"message": "not authenticated",
			})
			c.Abort()
			return
		}

		ctx := context.WithValue(c.Request.Context(), "user_id", id)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func CheckIfLoggedIn(c context.Context) (int, error) {
	userID, ok := c.Value("user_id").(int)
	if !ok {
		return 0, fmt.Errorf("not logged")
	}

	return userID, nil
}
