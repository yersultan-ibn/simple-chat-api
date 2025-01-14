package auth

import (
	"database/sql"
	"encoding/json"
	"ibn-salamat/simple-chat-api/database"
	notFound "ibn-salamat/simple-chat-api/handlers/not-found"
	"ibn-salamat/simple-chat-api/types"
	"log"
	"net/http"
	"strings"
)

func CheckConfirmCodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		notFound.NotFound(w, r)
		return
	}

	type Data struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	var data Data

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&data)

	data.Email = strings.Trim(data.Email, " ")
	data.Code = strings.Trim(data.Code, " ")

	if err != nil || data.Email == "" || data.Code == "" {
		jsonBody, _ := json.Marshal(types.ResponseMap{
			"errorMessage": "Incorrect json. Required fields: email, code.",
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonBody)

		return
	}

	var confirmed bool
	var leftTriesCount int
	var confirmationCode string

	err = database.DB.QueryRow(
		`SELECT confirmed, left_tries_count, confirmation_code
	FROM users_confirmation 
	WHERE email = $1 
	`, data.Email,
	).Scan(&confirmed, &leftTriesCount, &confirmationCode)

	if err != nil {
		errorMessage := err.Error()

		if err == sql.ErrNoRows {
			errorMessage = "Time for confirmation is up or you tried all tries. Please, start from the beginning."
		}

		jsonBody, _ := json.Marshal(types.ResponseMap{
			"errorMessage": errorMessage,
		})

		w.WriteHeader(http.StatusBadRequest)
		w.Write(jsonBody)
		return
	}

	if confirmed {
		jsonBody, _ := json.Marshal(types.ResponseMap{
			"message": "Already confirmed",
		})

		w.WriteHeader(http.StatusOK)
		w.Write(jsonBody)
		return
	}

	if data.Code != confirmationCode {
		if leftTriesCount > 1 {
			jsonBody, _ := json.Marshal(types.ResponseMap{
				"errorMessage":     "Incorrect code confirmation",
				"left_tries_count": leftTriesCount - 1,
			})

			w.WriteHeader(http.StatusBadRequest)
			w.Write(jsonBody)

			_, err = database.DB.Exec(`
				UPDATE users_confirmation SET left_tries_count = $2 WHERE email = $1
			`, data.Email, leftTriesCount-1)

			if err != nil {
				log.Println(err)
			}
			return
		} else if leftTriesCount == 1 {
			jsonBody, _ := json.Marshal(types.ResponseMap{
				"errorMessage":     "Try with new code.",
				"left_tries_count": leftTriesCount - 1,
			})

			w.WriteHeader(http.StatusBadRequest)
			w.Write(jsonBody)

			_, err = database.DB.Exec(`
				DELETE FROM users_confirmation
				WHERE email = $1
			`, data.Email)

			if err != nil {
				log.Println(err)
			}

			return
		}

	}

	if data.Code == confirmationCode {
		jsonBody, _ := json.Marshal(types.ResponseMap{
			"message": "Success",
		})

		w.WriteHeader(http.StatusOK)
		w.Write(jsonBody)

		_, err = database.DB.Exec(`
		UPDATE users_confirmation confirmed SET confirmed = TRUE WHERE email = $1
		`, data.Email)

		if err != nil {
			log.Println(err)
		}
		return
	}

	jsonBody, _ := json.Marshal(types.ResponseMap{
		"errorMessage": "Internal server error. Strange case",
	})

	w.WriteHeader(http.StatusInternalServerError)
	w.Write(jsonBody)
}
