package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"code.google.com/p/go.crypto/scrypt"

	"github.com/BurntSushi/toml"
	"github.com/acsellers/multitemplate"
	_ "github.com/acsellers/multitemplate/terse"
	"github.com/acsellers/zem/store"
	"github.com/gorilla/securecookie"
)

var (
	Setup    = flag.Bool("setup", false, "Migrates database to newest schema, and loads admin accounts if necessary")
	Conn     *store.Conn
	Cookie   *securecookie.SecureCookie
	Template *multitemplate.Template
)

var Config struct {
	DBType       string
	DBHost       string
	DBUser       string
	DBPass       string
	DBName       string
	DBPort       string
	HashSecret   string
	ResourcePath string
}

func init() {
	f, err := os.Open("zem.toml")
	if err != nil {
		log.Fatal("Couldn't open config file, where's zem.toml?")
	}

	_, err = toml.DecodeReader(f, &Config)
	if err != nil {
		log.Fatal("Couldn't decode config file, got error: ", err)
	}

	hashKey, err := scrypt.Key(
		[]byte("Hash Secret"),
		[]byte(Config.HashSecret),
		16384, 8, 1, 32,
	)
	if err != nil {
		log.Fatal("hash key", err)
	}
	hideKey, err := scrypt.Key(
		[]byte("Hide Secret"),
		[]byte(Config.HashSecret),
		16384, 8, 1, 32,
	)
	if err != nil {
		log.Fatal("encrypt key", err)
	}
	Cookie = securecookie.New(hashKey, hideKey)

	Template = CompileTemplates()
}

func main() {
	flag.Parse()
	if *Setup {
		log.Println("Running database setup")
		SetupDB()
	} else {
		log.Println("No setup required")
		ConnectDB()
	}

	log.Fatal(http.ListenAndServe(":8088", Router()))
}

func CompileTemplates() *multitemplate.Template {
	fmt.Println("Compiling Templates")
	t := multitemplate.New("")
	t.Base = filepath.Join(Config.ResourcePath, "html")
	t, err := t.ParseGlob(filepath.Join(Config.ResourcePath, "html", "*.html.*"))
	if err != nil {
		log.Fatal("Couldn't compile templates, error was:", err)
	}
	for _, tmpl := range t.Templates() {
		fmt.Println(tmpl.Name())
	}
	return t
}
