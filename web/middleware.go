package web

import (
	b64 "encoding/base64"

	"fmt"
	"io"
	"net/http"
	"strings"
)

func authMiddleware(next http.Handler, passphrase string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			w.Header().Add("WWW-Authenticate", "Basic realm=\"401\"")
			w.WriteHeader(http.StatusUnauthorized)
			io.WriteString(w, "Authentication required.")
			return
		}
		println("Unparse auth ", auth)
		auth = strings.Split(auth, " ")[1] // in "Basic <encoded>" only keep the encoded part
		dec, err := b64.StdEncoding.DecodeString(auth)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, "Error parsing the authentication")
			return
		}
		auth = string(dec)
		auth = strings.Split(auth, ":")[1] // in "user:pass" only save the pass

		fmt.Printf("Tried to access using %s", auth)
		if auth != passphrase {
			w.Header().Add("WWW-Authenticate", "Basic realm=\"401\"")
			w.WriteHeader(http.StatusUnauthorized)
			io.WriteString(w, "Authentication required.")
			return
		}

		next.ServeHTTP(w, r)

	})
}
