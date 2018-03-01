package app

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type slicedresponse struct {
	Total       int         `json:"total"`
	CurrentPage int         `json:"current_page"`
	PerPage     int         `json:"per_page"`
	Members     interface{} `json:"members"`
	Count       int         `json:"count"`
}

func asJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func getSection(w http.ResponseWriter, r *http.Request, sl slicer) {
	pp, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	curr, _ := strconv.Atoi(r.URL.Query().Get("page"))

	if pp == 0 {
		pp = 10
	}

	if curr == 0 {
		curr = 1
	}

	count := sl.length()
	start := min((curr-1)*pp, count)
	end := min(start+pp, count)
	members := sl.slice(start, end)

	resp := slicedresponse{
		Total:       count,
		CurrentPage: curr,
		PerPage:     pp,
		Members:     members,
		Count:       members.length(),
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		errorJSON(w, "Unable to return JSON body: %s", err.Error())
		return
	}
}

func getSingle(w http.ResponseWriter, r *http.Request, name string, sl slicer) {
	id := r.Context().Value(ID_KEY).(int)
	total := sl.length()

	if id > total {
		errorJSON(w, "No %s found with ID: %v", name, id)
		return
	}

	if err := json.NewEncoder(w).Encode(sl.uniq(id - 1)); err != nil {
		errorJSON(w, "Unable to return JSON body: %s", err.Error())
		return
	}
}

func genhome(r *http.Request, resources ...string) map[string][]string {
	prefix := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		prefix += "s"
	}

	m := make(map[string][]string)

	for _, v := range resources {
		m[v] = []string{
			prefix + "://" + r.Host + "/f/" + v + "",
			prefix + "://" + r.Host + "/f/" + v + "/{id}",
		}
	}

	return m
}

func fakerest(w http.ResponseWriter, r *http.Request) {
	resp := genhome(r, "users", "domains", "products", "posts")

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		errorJSON(w, "Unable to return JSON body: %s", err.Error())
		return
	}
}

func fakeusers(w http.ResponseWriter, r *http.Request) {
	getSection(w, r, systemUsers)
}

func fakesingleuser(w http.ResponseWriter, r *http.Request) {
	getSingle(w, r, "user", systemUsers)
}

func fakeposts(w http.ResponseWriter, r *http.Request) {
	getSection(w, r, systemPosts)
}

func fakesinglepost(w http.ResponseWriter, r *http.Request) {
	getSingle(w, r, "post", systemPosts)
}

func fakedomains(w http.ResponseWriter, r *http.Request) {
	getSection(w, r, systemDomains)
}

func fakesingledomain(w http.ResponseWriter, r *http.Request) {
	getSingle(w, r, "domain", systemDomains)
}

func fakeproducts(w http.ResponseWriter, r *http.Request) {
	getSection(w, r, systemProducts)
}

func fakesingleproduct(w http.ResponseWriter, r *http.Request) {
	getSingle(w, r, "product", systemProducts)
}
