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

func AuthMiddleware (next echo.HandlerFunc) echo.HandlerFunc{

	// By including Authorization header in Vary, the server ensures that the casching system does not cache responses meant for authorized users and serve them to unauthorized users and vice versa
	return func (c echo.Context) error{

		c.Response().Header().Add("Vary", "Authorization")
		// retrieving the authHeader from the Headers
		authHeader := c.Request().Header.Get("Authorization")

		fmt.Println("AuthHeader", authHeader)
		// the authHeader returns the value of the token in the headers together the Bearer keyword
		// next we check for this Bearer keyword if it exists

		if strings.HasPrefix(authHeader, "Bearer ") == false{

			// if the Bearer keyword does not exist, (false) return an error
			return common.Send

		}

		return next(c)
	}

}