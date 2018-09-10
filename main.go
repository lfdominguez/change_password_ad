package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/micro/go-config"
	"github.com/micro/go-config/source/file"
	zxcvbn "github.com/nbutton23/zxcvbn-go"
	"github.com/paleg/libadclient"
)

var box = packr.NewBox("./template")

func showIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	str, _ := box.MustString("index.html")

	t := template.Must(template.New("index").Parse(str))

	t.Execute(w, map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
	})
}

func changePassword(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(0)

	user := r.PostFormValue("user")
	old_pass := r.PostFormValue("old_pass")
	new_pass := r.PostFormValue("new_pass")

	adclient.New()
	defer adclient.Delete()

	params := adclient.DefaultADConnParams()

	var conf = config.Map()

	params.Domain = conf["ad_domain"].(string)
	params.Binddn = user + "@" + params.Domain
	params.Bindpw = old_pass

	params.Secured = false

	params.Timelimit = 60
	params.Nettimeout = 60

	type Response struct {
		Ok    bool
		Error string
	}

	resp := Response{}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := adclient.Login(params); err != nil {
		log.Println("Failed to AD login:", err)
		if strings.Contains(err.Error(), "resolving") {
			resp.Ok = false
			resp.Error = "AD ADDRESS NOT RESOLVED"
		} else {
			resp.Ok = false
			resp.Error = "USER AUTH FAILED"
		}

		respJson, _ := json.Marshal(resp)
		w.Write(respJson)

		return
	}

	params.Binddn = conf["ad_admin_user"].(string) + "@" + params.Domain
	params.Bindpw = conf["ad_admin_pass"].(string)

	if err := adclient.Login(params); err != nil {
		log.Println("Failed to AD login:", err)
		if strings.Contains(err.Error(), "resolving") {
			resp.Ok = false
			resp.Error = "AD ADDRESS NOT RESOLVED"
		} else {
			resp.Ok = false
			resp.Error = "ADMIN AUTH FAILED"
		}

		respJson, _ := json.Marshal(resp)
		log.Println("JSON:", respJson)
		w.Write(respJson)

		return
	}

	if err := adclient.SetUserPassword(user, new_pass); err != nil {
		log.Println("Failed to Change password:", err)
		if strings.Contains(err.Error(), "Constraint violation") {
			resp.Ok = false
			resp.Error = "CONTRAINS"
		} else {
			resp.Ok = false
			resp.Error = "USER AUTH FAILED"
		}

		respJson, _ := json.Marshal(resp)
		w.Write(respJson)

		return
	}

	resp.Ok = true

	respJson, _ := json.Marshal(resp)
	w.Write(respJson)
}

func checkPassword(w http.ResponseWriter, r *http.Request) {
	type Password struct {
		Password string `json:"-"`
	}

	password := Password{}

	json.NewDecoder(r.Body).Decode(&password)

	strength := zxcvbn.PasswordStrength(string(password.Password), nil)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	respJson, _ := json.Marshal(strength)
	w.Write(respJson)
}

func main() {
	if _, err := os.Stat("config.yml"); os.IsNotExist(err) {
		log.Fatal("Configuration file not found")
	}

	config.Load(file.NewSource(
		file.WithPath("config.yml"),
	))

	r := mux.NewRouter()

	r.HandleFunc("/", showIndex)
	r.HandleFunc("/changePassword", changePassword).Methods("POST")

	r.Handle("/logo.png", http.FileServer(box))
	r.Handle("/style.css", http.FileServer(box))
	r.Handle("/juration.js", http.FileServer(box))
	r.Handle("/zxcvbn.js", http.FileServer(box))

	// CSRF := csrf.Protect([]byte(config.Map()['csrf_key']))

	err := http.ListenAndServe(":9090", r /*CSRF(r)*/)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
