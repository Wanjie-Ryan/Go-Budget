// supply jwt
// middleware intercepts and validates jwt
// if jwt is invalid, returns 401
// middleware attaches current user with current context
package middlewares

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Wanjie-Ryan/Go-Budget/common"
	"github.com/Wanjie-Ryan/Go-Budget/internal/models"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// to have access to the instantiated DB
type AppMiddleware struct {
	DB *gorm.DB
}

func (appMiddleware *AppMiddleware) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {

	// By including Authorization header in Vary, the server ensures that the casching system does not cache responses meant for authorized users and serve them to unauthorized users and vice versa
	return func(c echo.Context) error {

		c.Response().Header().Add("Vary", "Authorization")
		// retrieving the authHeader from the Headers
		authHeader := c.Request().Header.Get("Authorization")

		fmt.Println("AuthHeader", authHeader)
		// the authHeader returns the value of the token in the headers together the Bearer keyword
		// next we check for this Bearer keyword if it exists

		if !strings.HasPrefix(authHeader, "Bearer ") {

			// if the Bearer keyword does not exist, (false) return an error
			return common.SendUnauthorizedResponse(c, "Invalid Authorization")

		}

		// if the authHeader is available, then get the token by splitting the authHeader at the space
		authHeaderSplit := strings.Split(authHeader, " ")
		accessToken := authHeaderSplit[1]

		claims, err := common.ParseJWT(accessToken)

		if err != nil {
			return common.SendUnauthorizedResponse(c, err.Error())
		}
		fmt.Println("extracted claims", claims)
		// if claim is available, then check the claim has not yet expired.

		if common.IsClaimExpired(claims) {

			// if the claim is expired true, then return token expired
			return common.SendUnauthorizedResponse(c, "Token Expired")
		}

		// using the id that is in the claims, we'll try to get that particular user from the DB
		// user, err := common.GetUserByID(appMiddleware.DB, claims.ID)
		// usermodel instance, which will act as an interface for storing the data that has been retrieved with the query
		var user models.UserModel
		// the code below sort of reads, select * from users where id = claims.ID
		result := appMiddleware.DB.First(&user, claims.ID)
		// there is a possibility that the record may not be found, hence use gorm.ErrRecordNotFound
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return common.SendUnauthorizedResponse(c, "Invalid Token")
		} else if result.Error != nil {
			return common.SendUnauthorizedResponse(c, "Invalid Token")
		}
		// middleware attaches current user with current context
		// setting the authenticated user to the context
		// c.set is a function that lets you attach any data to the current request context so that other handlers or middlewares can access it later
		c.Set("user", user)

		return next(c)
	}

}
