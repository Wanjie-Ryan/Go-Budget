// supply jwt
// middleware intercepts and validates jwt
// if jwt is invalid, returns 401
// middleware attaches current user with current context
package middlewares

import (
	"fmt"
	"strings"

	"github.com/Wanjie-Ryan/Go-Budget/common"
	"github.com/labstack/echo/v4"
)

func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {

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

		return next(c)
	}

}
