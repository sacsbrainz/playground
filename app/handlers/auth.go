package handlers

import (
	"database/sql"
	"log/slog"
	"net/http"
	"playground/app/cache"
	"playground/app/db"
	model "playground/app/types"
	"playground/app/utils"
	"playground/app/views/auth"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	// urlRegex   = regexp.MustCompile(`^(https?:\/\/)?(www\.)?([a-zA-Z0-9\-]+\.)+[a-zA-Z]{2,}(\/[a-zA-Z0-9\-._~:\/?#\[\]@!$&'()*+,;=]*)?$`)
)

func SignUpEmail(w http.ResponseWriter, r *http.Request) {
	// Initialize the error state
	formState := utils.NewGlobalFormState()

	// get form values
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	confirmPassword := r.FormValue("confirm_password")

	// Store form data
	formState.SetFormValues("first_name", firstName)
	formState.SetFormValues("last_name", lastName)
	formState.SetFormValues("email", email)
	formState.SetFormValues("password", password)
	formState.SetFormValues("confirm_password", confirmPassword)

	if email == "" {
		formState.AddError("email", "Email is required")
	}
	if !emailRegex.MatchString(email) {
		formState.AddError("email", "Invalid email address")
	}
	if password == "" {
		formState.AddError("password", "Password is required")
	}

	if password != confirmPassword {
		formState.AddError("confirm_password", "Password and Confirm password do not match")

	}
	if formState.HasErrors() {
		component := auth.SignUpForm(formState)
		w.Header().Set("content-type", "text/html; charset=utf-8")
		component.Render(r.Context(), w)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("hashing password", err)
		utils.Redirect(w, r, http.StatusSeeOther, "/500")
		return
	}

	query := `INSERT INTO users (id, first_name, last_name, email, password_hash) VALUES (?, ?, ?, ?, ?)`
	userID := strings.Replace(uuid.New().String(), "-", "", -1)

	_, dbErr := db.GetDb().Exec(query, userID, firstName, lastName, email, string(hash))

	if dbErr != nil {
		slog.Error("dbErr", dbErr)
		if dbErr.Error() == "UNIQUE constraint failed: users.email" {
			formState.AddError("email", "Email already in use")
			component := auth.SignUpForm(formState)
			w.Header().Set("content-type", "text/html; charset=utf-8")
			component.Render(r.Context(), w)
			return
		}
		utils.Redirect(w, r, http.StatusSeeOther, "/500")
		return
	}
	sessionID := strings.Replace(uuid.New().String(), "-", "", -1)

	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:     "token",
		Value:    sessionID,
		Expires:  expiration,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
	sessionName := "auth:" + sessionID

	// userId, ip, ua, createdAt, expiresAt
	session := userID + "::" + r.RemoteAddr + "::" + r.UserAgent() + "::" + time.Now().String() + "::" + time.Now().String()

	err = cache.GetCache().Str().SetExpires(sessionName, session, time.Since(expiration))
	if err != nil {
		slog.Error("cache", err)
		utils.Redirect(w, r, http.StatusSeeOther, "/500")
		return
	}

	utils.Redirect(w, r, http.StatusSeeOther, "/dashboard")
}
func LogInEmail(w http.ResponseWriter, r *http.Request) {
	// Initialize the error state
	formState := utils.NewGlobalFormState()

	// get form values
	email := r.FormValue("email")
	password := r.FormValue("password")

	// Store form data
	formState.SetFormValues("email", email)
	formState.SetFormValues("password", password)

	if email == "" {
		formState.AddError("error", "Email or password is incorrect")
	}
	if !emailRegex.MatchString(email) {
		formState.AddError("error", "Email or password is incorrect")
	}
	if password == "" {
		formState.AddError("error", "Email or password is incorrect")
	}

	if formState.HasErrors() {
		component := auth.LogInForm(formState)
		w.Header().Set("content-type", "text/html; charset=utf-8")
		component.Render(r.Context(), w)
		return
	}

	var user model.User
	query := `SELECT id, email, password_hash FROM users WHERE email = ? AND deleted_at IS NULL`

	dbErr := db.GetDb().QueryRow(query, email).Scan(&user.Id, &user.Email, &user.PasswordHash)

	if dbErr != nil {
		slog.Error("dbErr", dbErr)
		if dbErr == sql.ErrNoRows {

			formState.AddError("error", "Email or password is incorrect")

			component := auth.LogInForm(formState)
			w.Header().Set("content-type", "text/html; charset=utf-8")
			component.Render(r.Context(), w)
			return
		}

		utils.Redirect(w, r, http.StatusSeeOther, "/500")
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		formState.AddError("error", "Email or password is incorrect")

		component := auth.LogInForm(formState)
		w.Header().Set("content-type", "text/html; charset=utf-8")
		component.Render(r.Context(), w)
		return
	}

	sessionID := strings.Replace(uuid.New().String(), "-", "", -1)

	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{
		Name:     "token",
		Value:    sessionID,
		Expires:  expiration,
		Path:     "/",
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, &cookie)
	sessionName := "auth:" + sessionID

	// userId, ip, ua, createdAt, expiresAt
	session := user.Id.String() + "::" + r.RemoteAddr + "::" + r.UserAgent() + "::" + time.Now().String() + "::" + time.Now().String()

	err = cache.GetCache().Str().SetExpires(sessionName, session, time.Since(expiration))
	if err != nil {
		slog.Error("cache", err)
		utils.Redirect(w, r, http.StatusSeeOther, "/500")
		return
	}

	utils.Redirect(w, r, http.StatusSeeOther, "/dashboard")
}
