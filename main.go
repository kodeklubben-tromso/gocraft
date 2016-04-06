package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"

	"github.com/gorilla/mux"
)

var scriptdir string

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	f, err := ioutil.ReadFile("index.html")
	if err != nil {
		http.Error(w, "Could not find index.html", 500)
	} else {
		w.Write(f)
	}
}
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	modname := r.FormValue("modname")
	mod := r.FormValue("mod")

	if badfilename(modname) {
		http.Error(w, "Hei! Navnet på modden din skal kun ha små bokstaver. Ingen mellomrom eller rare tegn som f.eks ?/.,;;", 500)
		return
	}

	userdir := scriptdir + username
	err := os.MkdirAll(userdir, 0700)
	if err != nil {
		fmt.Println("could not create user script dir:", err)
		http.Error(w, "Could not create user script directory", 500)
		return
	}

	err = ioutil.WriteFile(userdir+"/"+modname+".js", []byte(mod), 0700)
	if err != nil {
		fmt.Println("Could not write mod", err)
		http.Error(w, "Could not write mod", 500)
		return
	}

	fmt.Println("user", username, "uploaded", modname)
	fmt.Fprintf(w, "Hei "+username+"! Vi har mottatt modden din "+modname+" og installert den! Sjekk om du får kjørt den på mc.cs.uit.no!")
}

func badfilename(filename string) bool {
	r, _ := regexp.Compile(`[[:punct:]]`)
	if r.MatchString(filename) {
		return true
	}

	r, _ = regexp.Compile(`[[:blank:]]`)
	if r.MatchString(filename) {
		return true
	}
	return false
}

func main() {
	dir := flag.String("scriptdir", ".", "scriptcraft server directory")
	flag.Parse()
	scriptdir = *dir + "/scripts/"

	fmt.Println("Storing mods in", scriptdir)

	fmt.Println("Starting server at :8080")
	r := mux.NewRouter()

	r.HandleFunc("/", IndexHandler)

	r.HandleFunc("/upload", UploadHandler)
	http.Handle("/", r)
	fmt.Println(http.ListenAndServe(":8080", nil))

}
