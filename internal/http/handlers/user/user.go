package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/aman-void/go-http-server/internal/storage"
	"github.com/aman-void/go-http-server/internal/types"
	"github.com/aman-void/go-http-server/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

// Create User
func New(storage storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user types.User
		// deserialize request body and validate
		err := json.NewDecoder(r.Body).Decode(&user)

		if errors.Is(err, io.EOF) {
			response.WriteJson(
				w,
				http.StatusBadRequest,
				response.NewError(errors.New("fields are empty")),
			)

			return
		}

		if err != nil {
			response.WriteJson(
				w,
				http.StatusBadRequest,
				response.NewError(err),
			)
			return
		}

		// request validation
		err = validator.New().Struct(&user)

		if err != nil {
			// type assertion
			validationErrors, ok := err.(validator.ValidationErrors)

			if !ok {
				response.WriteJson(
					w,
					http.StatusInternalServerError,
					response.NewError(err),
				)
				return
			}

			var validationErrMsg []error

			for _, validationErr := range validationErrors {
				switch validationErr.ActualTag() {
				case "required":
					validationErrMsg = append(
						validationErrMsg,
						fmt.Errorf("field %s is required", validationErr.Field()),
					)

				default:
					validationErrMsg = append(
						validationErrMsg,
						fmt.Errorf("invalid value for field %s", validationErr.Field()),
					)
				}
			}

			response.WriteJson(
				w,
				http.StatusBadRequest,
				response.NewError(errors.Join(validationErrMsg...)),
			)

			return

		}

		// create user in database
		lastId, err := storage.CreateUser(
			user.Name,
			user.Email,
			user.Age,
		)

		if err != nil {
			response.WriteJson(
				w,
				http.StatusInternalServerError,
				response.NewError(err),
			)
			return
		}

		slog.Info(
			"user created",
			slog.Int64("lastId", lastId),
		)

		response.WriteJson(
			w,
			http.StatusCreated,
			response.NewSuccess(
				"user created successfully",
				map[string]int64{"lastId": lastId},
			),
		)

	}
}
