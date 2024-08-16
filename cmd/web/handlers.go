package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"ssnipp.com/internal/models"
	"ssnipp.com/internal/validator"
)

type snippetCreateForm struct {
	Content             string `form:"content"`
	Language            string `form:"language"`
	validator.Validator `form:"-"`
}

type userSignupForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

// Home page handler
func (app *application) home(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	// Load available languages
	data.Languages = getLanguages()

	// Initialize form with default values
	data.Form = snippetCreateForm{
		Language: "plaintext",
	}

	app.render(w, r, http.StatusOK, "home.html", data)
}

// View snippet handler
func (app *application) snippetView(w http.ResponseWriter, r *http.Request) {
	// Get the ID of the snippet from the URL parameter
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	// Retrieve the snippet from the database
	snippet, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			http.NotFound(w, r)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// Prepare template data
	data := app.newTemplateData(r)
	data.Snippet = snippet

	app.render(w, r, http.StatusOK, "view.html", data)
}

// Create snippet handler (POST)
func (app *application) snippetCreatePost(w http.ResponseWriter, r *http.Request) {
	var form snippetCreateForm

	// Decode the form data
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form contents
	form.CheckField(validator.NotBlank(form.Content), "content", "This field cannot be blank")
	form.CheckField(validator.PermittedValue(form.Language, getLanguageKeys()), "language", "Choose a valid language")

	// If there are any validation errors, re-display the form
	if !form.Valid() {
		data := app.newTemplateData(r)

		data.Languages = getLanguages()

		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "home.html", data)
		return
	}

	// Insert the snippet into the database
	id, err := app.snippets.Insert(form.Content, form.Language)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Add a flash message to the session
	app.sessionManager.Put(r.Context(), "flash", "Snippet successfully created!")

	// Redirect to the snippet view page
	http.Redirect(w, r, fmt.Sprintf("/view/%d", id), http.StatusSeeOther)
}

// User signup page handler
func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userSignupForm{}
	app.render(w, r, http.StatusOK, "signup.html", data)
}

// User signup handler (POST)
func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	var form userSignupForm

	// Decode the form data
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form contents.
	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "This field must be at least 8 characters long")

	// If there are any validation errors, re-display the signup form
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	// Try to create a new user record in the database
	err = app.users.Insert(form.Name, form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "Email address is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
		} else {
			app.serverError(w, r, err)
		}

		return
	}

	// Add a flash message to the session
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	// Redirect the user to the login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// User login page handler
func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)
	data.Form = userLoginForm{}
	app.render(w, r, http.StatusOK, "login.html", data)
}

// User login handler (POST)
func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	// Decode the form data
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Validate the form contents
	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	// If there are any validation errors, re-display the login form
	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form

		app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	// Try to authenticate the user
	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			form.AddNonFieldError("Email or password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form

			app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// Renew the session token
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Store the user's ID in the session
	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	// Redirect to the original destination or the home page
	path := app.sessionManager.PopString(r.Context(), "redirectPathAfterLogin")
	if path != "" {
		http.Redirect(w, r, path, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// User logout handler (POST)
func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	// Renew the session token
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Remove the user's ID from the session
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	// Add a flash message to the session
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	// Redirect to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Ping handler for health checks
func ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
