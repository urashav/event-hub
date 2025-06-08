package httputils

import (
	"encoding/json"
	"log"
	"net/http"
)

// Response представляет собой обертку для HTTP-ответа
type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
}

// SendJSON отправляет JSON-ответ с указанным кодом состояния
func SendJSON(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// SendError отправляет JSON-ответ с ошибкой
func SendError(w http.ResponseWriter, message string, code int) {
	SendJSON(w, Response{Error: message}, code)
}

// SendSuccess отправляет успешный JSON-ответ
func SendSuccess(w http.ResponseWriter, data interface{}, code int) {
	SendJSON(w, Response{Data: data}, code)
}

// ErrorResponse возвращает стандартизированный ответ с ошибкой
var ErrorResponse = struct {
	BadRequest       func(w http.ResponseWriter, msg string)
	Unauthorized     func(w http.ResponseWriter, msg string)
	Forbidden        func(w http.ResponseWriter, msg string)
	NotFound         func(w http.ResponseWriter, msg string)
	InternalError    func(w http.ResponseWriter, msg string)
	MethodNotAllowed func(w http.ResponseWriter, msg string)
}{
	BadRequest: func(w http.ResponseWriter, msg string) {
		SendError(w, msg, http.StatusBadRequest)
	},
	Unauthorized: func(w http.ResponseWriter, msg string) {
		SendError(w, msg, http.StatusUnauthorized)
	},
	Forbidden: func(w http.ResponseWriter, msg string) {
		SendError(w, msg, http.StatusForbidden)
	},
	NotFound: func(w http.ResponseWriter, msg string) {
		SendError(w, msg, http.StatusNotFound)
	},
	InternalError: func(w http.ResponseWriter, msg string) {
		SendError(w, msg, http.StatusInternalServerError)
	},
	MethodNotAllowed: func(w http.ResponseWriter, msg string) {
		SendError(w, msg, http.StatusMethodNotAllowed)
	},
}

// Чтение JSON из тела запроса
func DecodeJSON(r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// Получение параметров запроса
func GetQueryParam(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}
