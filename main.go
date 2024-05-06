package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"strings"

	"golang.org/x/net/dict"
	pluralize "github.com/gertd/go-pluralize"
)

func main() {
	dictServer := flag.String("dict-server", "dict.org:2628", "Specify the dict server address and port")
	listenAddr := flag.String("listen-addr", ":3010", "Specify the local listen address, prefix with unix: for a Unix domain socket")
	flag.Parse()

	var listener net.Listener
	var err error
	if strings.HasPrefix(*listenAddr, "unix:") {
		listener, err = net.Listen("unix", strings.TrimPrefix(*listenAddr, "unix:"))
	} else {
		listener, err = net.Listen("tcp", *listenAddr)
	}
	if err != nil {
		log.Fatalf("Error listening on %s: %v", *listenAddr, err)
	}

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/define", defineHandler(*dictServer))

	log.Printf("Listening on %s", *listenAddr)
	if err := http.Serve(listener, nil); err != nil {
		log.Fatal(err)
	}
}

func defineHandler(dictServer string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Printf("error parsing form: %v", err)
			http.Error(w, "error parsing form", http.StatusBadRequest)
			return
		}
		word := r.Form.Get("word")
		var params TmplParams
		if len(word) > 0 {
			defs, err := getDefinitions(dictServer, word)
			if err != nil {
				log.Printf("error defining %s: %v", word, err)
				http.Error(w, fmt.Sprintf("error getting definition for %s: %v", word, err), http.StatusInternalServerError)
				return
			}
			params.ShowResults = true
			params.SearchTerm = word
			params.Definitions = defs
		}
		renderTemplate(w, home, params)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, home, TmplParams{})
}

func renderTemplate(w http.ResponseWriter, t *template.Template, params TmplParams) {
	err := t.Execute(w, params)
	if err != nil {
		log.Printf("Error rendering template %s: %v", t.Name(), err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

var pl *pluralize.Client

func init() {
	pl = pluralize.NewClient()
	pl.AddSingularRule("(rect)a", "$1um")
}

func getDefinitions(serverAddr string, word string) ([]*dict.Defn, error) {
	cli, err := dict.Dial("tcp", serverAddr)
	if err != nil {
		return nil, fmt.Errorf("error connectiong to %s: %v", serverAddr, err)
	}
	defer cli.Close()
	defs, err := cli.Define("*", word)
	if err != nil {
		if strings.Contains(err.Error(), "552 no match") {
			word_singular := pl.Singular(word)
			if word_singular != word {
				return getDefinitions(serverAddr, word_singular)
			}
			return nil, nil
		}
		return nil, fmt.Errorf("error from dict server: %v", err)
	}
	return defs, nil
}

type TmplParams struct {
	SearchTerm  string
	ShowResults bool
	Definitions []*dict.Defn
}

var home = template.Must(template.New("home").Parse(`
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <link rel="stylesheet" href="https://cdn.simplecss.org/simple.min.css">
    <style>
		div.center {
			text-align: center;
		}

		div#results {
			height: 60dvh;
		}
    </style>
    <title>Dictionary Server</title>
</head>
<body>
  <header>
    <h1>Dictionary Server</h1>
  </header>

  <main>
	<form action="/define" method="get">
	  <div class="center">
	    <input id="word" type="text" name="word" required autofocus></input>
	    <button name="search">Search</button>
	  </div>
    </form>
    {{if .ShowResults}}
		<div id="results">
		{{with .Definitions}}
      		{{range .}}
    	  		<h3>From {{.Dict.Desc}}</h3>
      	  		<pre>{{.Text|printf "%s"}}</pre>
    		{{end}}
    	{{else}}
    		<p>
	  			Found no definitions for {{.SearchTerm}}
    		</p>
    	{{end}}
		</div>
    {{end}}
  </main>

</body>
</html>
`))
