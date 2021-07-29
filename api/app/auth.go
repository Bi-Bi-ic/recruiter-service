package app

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/rgrs-x/service/api/models"
	u "github.com/rgrs-x/service/api/utils"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// APIAuthentication ...
func APIAuthentication(auths ...string) gin.HandlerFunc {

	return func(c *gin.Context) {
		response := make(map[string]interface{})

		xAPI := os.Getenv("X_API_KEY")

		sh1 := os.Getenv("X_SH1_FINGERPRINT")

		xAPIH := c.Request.Header.Get("x-api-key")

		sh1H := c.Request.Header.Get("x-sha1-fingerprint")

		if xAPI != xAPIH && sh1 != sh1H {
			response = u.Message(false, "Missing auth token config")
			c.Writer.Header().Set("Content-Type", "application/json")
			c.JSON(http.StatusForbidden, response)
			c.Abort()
			return
		}
		c.Next()
	}

}

// GeneralAuthentication ...
func GeneralAuthentication(auths ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		notAuth := []string{"/api/partner/sign_up", "/api/partner/sign_in"}

		tokenHeader := c.Request.Header.Get("Authorization")

		noneAuth(notAuth, c)

		ok := checkAuthrization(tokenHeader, c)
		if !ok {
			c.Abort()
			return
		}

		ok = checkFormatToken(tokenHeader, c)
		if !ok {
			c.Abort()
			return
		}

		tk, jwtToken, err := pareToken(tokenHeader)

		ok = checkErr(err, c)
		if !ok {
			c.Abort()
			return
		}

		checkTypeToken(tk, c)

		validToken(jwtToken, c)

		// pass all thing above
		fmt.Sprintf("User %", tk.UserId)

		c.Writer.Header().Set("user", tk.UserId.String())

		c.Next()
	}
}

// UserAuthentication ...
func UserAuthentication(auths ...string) gin.HandlerFunc {

	return func(c *gin.Context) {
		notAuth := []string{"/api/partner/sign_up", "/api/partner/sign_in"}

		tokenHeader := c.Request.Header.Get("Authorization")

		noneAuth(notAuth, c)

		ok := checkAuthrization(tokenHeader, c)
		if !ok {
			c.Abort()
			return
		}

		ok = checkFormatToken(tokenHeader, c)
		if !ok {
			c.Abort()
			return
		}

		tk, jwtToken, err := pareToken(tokenHeader)

		ok = checkErr(err, c)
		if !ok {
			c.Abort()
			return
		}

		checkTypeToken(tk, c)

		checkUserType(models.UserNormal, tk, c)

		validToken(jwtToken, c)

		// pass all thing above
		fmt.Sprintf("User %", tk.UserId)

		c.Writer.Header().Set("user", tk.UserId.String())

		c.Next()
	}
}

// PartnerAuthentication ...
func PartnerAuthentication(mode models.UserMode) gin.HandlerFunc {
	return func(c *gin.Context) {
		notAuth := []string{"/api/partner/sign_up", "/api/partner/sign_in"}

		tokenHeader := c.Request.Header.Get("Authorization")

		noneAuth(notAuth, c)

		ok := checkAuthrization(tokenHeader, c)
		if !ok {
			c.Abort()
			return
		}

		ok = checkFormatToken(tokenHeader, c)
		if !ok {
			c.Abort()
			return
		}

		tk, jwtToken, err := pareToken(tokenHeader)

		ok = checkErr(err, c)
		if !ok {
			c.Abort()
			return
		}

		checkTypeToken(tk, c)

		checkUserType(mode, tk, c)

		validToken(jwtToken, c)

		// pass all thing above
		fmt.Sprintf("User %", tk.UserId)

		c.Writer.Header().Set("user", tk.UserId.String())

		c.Next()
	}
}

// AdminAuthentication ...
func AdminAuthentication(mode models.UserMode) gin.HandlerFunc {
	return func(c *gin.Context) {
		notAuth := []string{"/api/admin/sign_up"}

		tokenHeader := c.Request.Header.Get("Authorization")

		noneAuth(notAuth, c)

		ok := checkAuthrization(tokenHeader, c)
		if !ok {
			c.Abort()
			return
		}

		ok = checkFormatToken(tokenHeader, c)
		if !ok {
			c.Abort()
			return
		}

		tk, jwtToken, err := pareToken(tokenHeader)

		ok = checkErr(err, c)
		if !ok {
			c.Abort()
			return
		}

		checkTypeToken(tk, c)

		checkUserType(mode, tk, c)

		validToken(jwtToken, c)

		// pass all thing above
		fmt.Sprintf("User %", tk.UserId)

		c.Writer.Header().Set("user", tk.UserId.String())

		c.Next()
	}
}

func noneAuth(paths []string, c *gin.Context) {
	requestPath := c.Request.URL.String()

	for _, value := range paths {

		if value == requestPath {
			c.Next()
			return
		}
	}
}

func checkAuthrization(tokenHeader string, c *gin.Context) bool {
	response := make(map[string]interface{})

	if tokenHeader == "" {
		response = u.Message(false, "Missing auth token")
		c.Writer.Header().Set("Content-Type", "application/json")
		c.JSON(http.StatusForbidden, response)
		return false
	}
	return true
}

func checkFormatToken(tokenHeader string, c *gin.Context) bool {
	response := make(map[string]interface{})
	splitted := strings.Split(tokenHeader, " ")
	if len(splitted) != 2 {
		response = u.Message(false, "Invalid/Malformed auth token")
		c.Writer.Header().Set("Content-Type", "application/json")
		c.JSON(http.StatusForbidden, response)
		return false
	}
	return true
}

func checkErr(err error, c *gin.Context) bool {
	response := make(map[string]interface{})
	if err != nil {
		response = u.Message(false, "Malformed authentication token")
		c.Writer.Header().Set("Content-Type", "application/json")
		c.JSON(http.StatusForbidden, response)
		return false
	}
	return true
}

func pareToken(tokenHeader string) (*models.Token, *jwt.Token, error) {

	splitted := strings.Split(tokenHeader, " ")
	tokenPart := splitted[1]
	tk := &models.Token{}

	jwtToken, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_password")), nil
	})

	if err != nil {
		return nil, nil, err
	}

	return tk, jwtToken, nil
}

func checkTypeToken(token *models.Token, c *gin.Context) {
	response := make(map[string]interface{})
	if token.Type != models.Now {
		response = u.Message(false, "Token is not valid.")
		c.Writer.Header().Set("Content-Type", "application/json")
		c.JSON(http.StatusForbidden, response)
		c.Abort()
		return
	}
}

func validToken(token *jwt.Token, c *gin.Context) {
	response := make(map[string]interface{})
	if !token.Valid {
		response = u.Message(false, "Token is not valid.")
		c.Writer.Header().Set("Content-Type", "application/json")
		c.JSON(http.StatusForbidden, response)
		c.Abort()
		return
	}
}

func checkUserType(mode models.UserMode, token *models.Token, c *gin.Context) {
	response := make(map[string]interface{})

	if mode != token.UserType {
		response = u.Message(false, "Your has reject by system")
		c.Writer.Header().Set("Content-Type", "application/json")
		c.JSON(http.StatusForbidden, response)
		c.Abort()
		return
	}
}
