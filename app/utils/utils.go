package utils

import (
	"net/http"
	"playground/public"
)

func StaticFileHandler(appEnv string) http.Handler {
	if appEnv == "production" {
		return staticProd()
	}
	return disableCache(staticDev())
}

func staticDev() http.Handler {
	return http.StripPrefix("/public/", http.FileServer(http.Dir("public")))
}

func staticProd() http.Handler {
	return http.StripPrefix("/public/", http.FileServer(http.FS(public.AssetsFS)))
}

func disableCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}

// GlobalFormState will store form-related errors and the form state.
type GlobalFormState struct {
	FormValues map[string]string // Key: field name, Value: field value
	FormErrors map[string]string // Key: field name, Value: error message
}

// NewGlobalFormState initializes and returns a new GlobalFormState.
func NewGlobalFormState() *GlobalFormState {
	return &GlobalFormState{
		FormErrors: make(map[string]string),
		FormValues: make(map[string]string),
	}
}

// AddError adds an error message for a specific field.
func (gfs *GlobalFormState) AddError(field, message string) {
	gfs.FormErrors[field] = message
}

// GetError retrieves the error message for a specific field.
func (gfs *GlobalFormState) GetError(field string) string {
	return gfs.FormErrors[field]
}

// HasError checks if there are any errors stored.
func (gfs *GlobalFormState) HasError(field string) bool {
	return len(gfs.FormErrors[field]) > 0
}

// HasErrors checks if there are any errors stored.
func (gfs *GlobalFormState) HasErrors() bool {
	return len(gfs.FormErrors) > 0
}

// SetFormValues sets the value for a specific form field.
func (gfs *GlobalFormState) SetFormValues(field, value string) {
	gfs.FormValues[field] = value
}

// GetFormValue retrieves the value for a specific form field.
func (gfs *GlobalFormState) GetFormValue(field string) string {
	return gfs.FormValues[field]
}

// Clear clears all errors in the GlobalFormState.
func (gfs *GlobalFormState) Clear() {
	gfs.FormValues = make(map[string]string)
	gfs.FormErrors = make(map[string]string)
}

// Redirect with HTMX support.
func Redirect(w http.ResponseWriter, r *http.Request, status int, url string) error {
	if len(r.Header.Get("HX-Request")) > 0 {
		w.Header().Set("HX-Redirect", url)
		w.WriteHeader(status)
		return nil
	}
	http.Redirect(w, r, url, status)
	return nil
}
