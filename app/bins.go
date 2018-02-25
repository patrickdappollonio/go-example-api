package app

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"google.golang.org/appengine"
	"google.golang.org/appengine/user"
)

const BIN_PREFIX = "requests_"

func bincreate(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{"ID": key()}

	if err := tmplCreate.Execute(w, data); err != nil {
		http.Error(w, "Unable to print page contents: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func loadPrevious(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(ID_KEY).(string)
		key := BIN_PREFIX + id

		bindata, err := getContent(r, key)
		if err != nil && err != notfound {
			log.Println("Unable to get contents for ID", id, "due to error:", err.Error())
			http.Error(w, "Error while fetching information from the store: "+err.Error(), http.StatusInternalServerError)
			return
		}

		responses := make([]debugresponse, 0)
		if s := len(bindata); s > 0 {
			responses = fromstring(bindata)
		}

		ctx := context.WithValue(r.Context(), BIN_PREFIX, responses)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func binget(w http.ResponseWriter, r *http.Request) {
	responses := r.Context().Value(BIN_PREFIX).([]debugresponse)
	id := r.Context().Value(ID_KEY).(string)

	heading := fmt.Sprintf("%d request", len(responses))
	if len(responses) != 1 {
		heading += "s"
	}

	data := map[string]interface{}{
		"ID":      id,
		"Records": responses,
		"Heading": heading,
	}

	if err := tmplRead.Execute(w, data); err != nil {
		http.Error(w, "Unable to print page contents: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func binsave(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	var (
		responses = r.Context().Value(BIN_PREFIX).([]debugresponse)
		id        = r.Context().Value(ID_KEY).(string)
		key       = BIN_PREFIX + id

		current = debugresponse{
			Headers: cleanup(r.Header),
			Method:  r.Method,
			Path:    r.URL.Path,
			Time:    time.Now(),
		}

		buf bytes.Buffer
	)

	if _, err := io.Copy(&buf, r.Body); err != nil {
		log.Println("Unable to read request body for ID", id, "- Error:", err.Error())
		http.Error(w, "Unable to read request body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	current.Body = buf.String()

	if d := cleanup(r.URL.Query()); len(d) > 0 {
		current.Query = d
	}

	responses = append([]debugresponse{current}, responses...)

	if _, err := saveContentWithKey(r, key, tostring(responses)); err != nil {
		log.Println("Unable to save slice of requests to storage:", err.Error(), "\n", "Contents:", fmt.Sprintf("%#v", responses))
		http.Error(w, "Unable to save request to the store: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Saved! Now check /inspector/%s to inspect it.", id)
}

func bindelete(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value(ID_KEY).(string)
	key := BIN_PREFIX + id
	ctx := appengine.NewContext(r)

	u := user.Current(ctx)
	if u == nil {
		dest, _ := user.LoginURL(ctx, r.URL.Host)
		http.Redirect(w, r, dest, http.StatusFound)
		return
	}

	if !u.Admin {
		dest, _ := user.LogoutURL(ctx, "/")
		fmt.Fprintf(w, "I'm sorry, but your email %q is not allowed as an administrator of this site.\nLogout: %s", u.Email, dest)
		return
	}

	deleteContent(r, key)
	fmt.Fprintln(w, "Inspector data for ID:", id, "was deleted from Memcache")
}
