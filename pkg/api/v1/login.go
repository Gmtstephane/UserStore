package userStoreAPi

import (
	models "UserStore/pkg/models"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgconn"
	"github.com/labstack/echo-contrib/session"
	"gorm.io/gorm"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

//LoginHandler sign in the User and return a cookie session
// @Summary User Login
// @Description Get the auth cookie token
// @Accept  json
// @Produce  json
// @Param user body models.User true "models.Users{}"
// @Success 200
// @Header 200 {string} Auth-cookie ""
// @Failure 401 {string} Unauthorized
// @Failure default {object} pgconn.PgError
// @Router /api/v1/login [post]
func (u *Userstore) loginHandler(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")

	//Redis session
	session, _ := session.Get("Auth-cookie", c)

	var user models.User   //from request
	var userDB models.User //database Response
	var err error          //Error handling

	if err = json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		internalErr(c.Response().Writer, err)
		return nil
	}

	//Todo : Handle error types / NotFound, internal error etc...
	if tx := u.db.Where("email = ?", user.Email).First(&userDB); tx.Error != nil {

		switch tx.Error.(type) {
		case *pgconn.PgError:
			err := tx.Error.(*pgconn.PgError)
			switch err.Code {
			case "23505": // Duplicate data
				c.Response().WriteHeader(http.StatusConflict)
				json.NewEncoder(c.Response().Writer).Encode(tx.Error)
			default: //Default database error
				c.Response().WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(c.Response().Writer).Encode(tx.Error)
			}
		default:
			if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
				c.Response().WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(c.Response().Writer).Encode(gorm.ErrRecordNotFound.Error())
				fmt.Println(gorm.ErrRecordNotFound.Error())
				return tx.Error
			}
			c.Response().WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(c.Response().Writer).Encode(tx.Error)
		}
	}

	//Password hashed compare
	if err = bcrypt.CompareHashAndPassword([]byte(userDB.Password), []byte(user.Password)); err != nil {
		c.Response().WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(c.Response()).Encode("Wrong password for user " + user.Email)
	}

	// save user as authenticated with some parameters
	session.Values["authenticated"] = true
	session.Values["email"] = user.Email
	session.Values["id"] = user.Model.ID
	session.Save(c.Request(), c.Response())

	// json.NewEncoder(w).Encode(userDB)
	return nil
}
