package routes

import (
	"encoding/json"
	"net/http"

	"github.com/JuanLopezAranzazu/go-restapi-authentication-jwt/db"
	"github.com/JuanLopezAranzazu/go-restapi-authentication-jwt/middlewares"
	"github.com/JuanLopezAranzazu/go-restapi-authentication-jwt/models"
	"github.com/JuanLopezAranzazu/go-restapi-authentication-jwt/utils"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validar campos
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email y contraseña son obligatorios", http.StatusBadRequest)
		return
	}

	// hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Error al generar la contraseña", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	// validar correo
	if err := db.DB.Create(&user).Error; err != nil {
		http.Error(w, "Email ya registrado", http.StatusBadRequest)
		return
	}

	// generar tokens
	accessToken, _ := utils.GenerateJWT(user.ID)
	refreshToken, _ := utils.GenerateRefreshToken(user.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(TokenResponse{AccessToken: accessToken, RefreshToken: refreshToken})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// validar correo
	var user models.User
	if err := db.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		http.Error(w, "Credenciales incorrectas", http.StatusUnauthorized)
		return
	}

	// validar contraseña
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		http.Error(w, "Credenciales incorrectas", http.StatusUnauthorized)
		return
	}

	// generar tokens
	accessToken, _ := utils.GenerateJWT(user.ID)
	refreshToken, _ := utils.GenerateRefreshToken(user.ID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(TokenResponse{AccessToken: accessToken, RefreshToken: refreshToken})
}

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, err := utils.ValidateRefreshToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "Refresh token inválido", http.StatusUnauthorized)
		return
	}

	// generar nuevo access token
	accessToken, _ := utils.GenerateJWT(claims.UserID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": accessToken,
	})
}

func MeHandler(w http.ResponseWriter, r *http.Request) {
	// obtener userID del contexto
	userID := r.Context().Value(middlewares.UserIDKey).(uint)

	var user models.User
	if err := db.DB.First(&user, userID).Error; err != nil {
		http.Error(w, "Usuario no encontrado", http.StatusNotFound)
		return
	}

	resp := struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
