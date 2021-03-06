package auth

import (
	"be/repository/auth"
	"errors"
	"net/http"

	"be/delivery/controllers/templates"
	"be/delivery/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type AuthController struct {
	repo auth.Auth
}

func New(repo auth.Auth) *AuthController {
	return &AuthController{
		repo: repo,
	}
}

func (ac *AuthController) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		Userlogin := Userlogin{}

		if err := c.Bind(&Userlogin); err != nil || Userlogin.UserName == "" || Userlogin.Password == "" {
			return c.JSON(http.StatusBadRequest, templates.BadRequest(nil, "invalid email or password", err))
		}

		checkedUser, err := ac.repo.Login(Userlogin.UserName, Userlogin.Password)

		if err != nil {
			log.Info(err)
			switch err.Error() {
			case "incorrect password":
				err = errors.New("incorrect password")
			case "record not found":
				err = errors.New("account is not found")
			default:
				err = errors.New("there's some problem is server")
			}

			return c.JSON(http.StatusInternalServerError, templates.InternalServerError(nil, err.Error(), nil))
		}

		// log.Info(checkedUser)
		token, err := middlewares.GenerateToken(checkedUser["data"].(string), checkedUser["type"].(string))

		if err != nil {
			log.Warn(err)
			err = errors.New("there's some problem is server")
			return c.JSON(http.StatusNotAcceptable, templates.BadRequest(http.StatusNotAcceptable, err.Error(), nil))
		}

		return c.JSON(http.StatusOK, templates.Success(nil, "success login", map[string]interface{}{
			"type":  checkedUser["type"].(string),
			"doctor_uid": checkedUser["doctor_uid"].(string),
			"token": token,
		}))
	}
}
