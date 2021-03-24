package userStoreAPi

import (
	models "UserStore/pkg/models"
	"encoding/json"
	"net/http"

	"github.com/jackc/pgconn"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

//LoginHandler sign in the User and return a cookie session
// ListAccounts godoc
// @Summary List accounts
// @Description get accounts
// @Accept  json
// @Produce  json
// @Param user body models.User true "models.Users{}"
// @Success 200
// @Header 200 {string} Auth-cookie ""
// @Failure 409,500 {object} pgconn.PgError
// @Failure default {object} pgconn.PgError
// @Router /signup [post]
func (u *Userstore) signUpHandler(c echo.Context) error {
	c.Response().Header().Set("Content-Type", "application/json")

	var err error   //Error handling
	var hash []byte //Encrypted password

	//Post Decode
	var user models.User
	if err = json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		internalErr(c.Response().Writer, err)

		return nil
	}

	//Check if user fields are correct
	ok, msg := user.Validate()
	if !ok {
		c.Response().WriteHeader(http.StatusBadRequest)
		json.NewEncoder(c.Response().Writer).Encode(msg)
		return nil
	}

	//Bcrypt hash the password
	if hash, err = bcrypt.GenerateFromPassword([]byte(user.Password), 10); err != nil {
		internalErr(c.Response().Writer, err)
		return nil
	}

	user.Password = string(hash)
	//Insert the password in database
	//TODO: Handle multiples types of error (user already exist, database error, unvalid param...)

	if tx := u.db.Create(&user); tx.Error != nil {
		err := tx.Error.(*pgconn.PgError)
		switch err.Code {
		case "23505":
			c.Response().WriteHeader(http.StatusConflict)
			json.NewEncoder(c.Response().Writer).Encode(tx.Error)
		default:
			c.Response().WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(c.Response().Writer).Encode(tx.Error)
		}
	}
	return nil
}
