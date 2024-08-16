package main

type contextKey string

// isAuthenticatedContextKey is a context key for the isAuthenticated value.
const isAuthenticatedContextKey = contextKey("isAuthenticated")
