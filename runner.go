package blast

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"
)

func Run(fName string) {
	if !strings.HasSuffix(fName, ".blast") {
		fName += ".blast"
	}

	data, err := ioutil.ReadFile(fName)

	if err != nil {
		log.Fatalf("Could not read file %v", fName)
	}

	Init()
	Parse(string(data))
}

// CLI runs the command
// line interface for the
// Blast interpreter
func CLI() {
	Init()
	cr := bufio.NewReader(os.Stdin)

	for {
		fmt.Print(">")
		input, err := cr.ReadString('\n')

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(ParseBasicLine(input[:len(input)-1]))
	}
}

// Playground runs a webserver
// with a playground for the
// Blast interpreter
func Playground() {
	http.HandleFunc("/", handleIndex)
	http.ListenAndServe("localhost:9000", nil)
}

// handleIndex returns the index view for the playground
func handleIndex(w http.ResponseWriter, r *http.Request) {
	fp := path.Join("templates", "index.html")
	tmpl, err := template.ParseFiles(fp)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
