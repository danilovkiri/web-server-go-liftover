package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"webServerBasicAuth/modules/config"
	"webServerBasicAuth/modules/cookie"
	"webServerBasicAuth/modules/utils"
)

type application struct {
	auth struct {
		username   string
		password   string
		configpath string
	}
}

func main() {
	app := new(application)
	app.auth.username = os.Getenv("AUTH_USERNAME")
	app.auth.password = os.Getenv("AUTH_PASSWORD")
	app.auth.configpath = os.Getenv("CONFIG")
	yamlConfig := config.ParseConfig(app.auth.configpath)

	if app.auth.username == "" {
		log.Fatal("Basic auth username must be provided")
	}
	if app.auth.password == "" {
		log.Fatal("Basic auth password must be provided")
	}

	mux := http.NewServeMux()
	mux.Handle("/index/", http.StripPrefix("/index/", http.FileServer(http.Dir("./public"))))
	mux.HandleFunc("/index/hg19to38/", app.basicAuth(app.mainHandler19to38))
	mux.HandleFunc("/index/hg38to19/", app.basicAuth(app.mainHandler38to19))

	srv := &http.Server{
		Addr:         fmt.Sprintf("%v:%v", yamlConfig.Constants.ServerIP, yamlConfig.Constants.ServerPort),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("starting server on %s", srv.Addr)
	err := srv.ListenAndServeTLS(yamlConfig.Constants.CertFile, yamlConfig.Constants.KeyFile)
	log.Fatal(err)
}

func (app *application) mainHandler38to19(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")
	r.ParseMultipartForm(1)
	file, handler, err := r.FormFile("uploadFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)
	tempFile, err := ioutil.TempFile("temp-files", handler.Filename+"_upload*")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	tempFile.Write(fileBytes)

	fmt.Println("### INFO: Temp file written to", tempFile.Name())
	cwd := utils.GetCwd() + "/"
	tempId := utils.GetTempId(tempFile.Name())
	outputFile := "output-files/hg38tohg19." + tempId + ".txt"
	genomeBuild := "hg38"
	executableCmd := utils.MakeCmdStruct(cwd, cwd+tempFile.Name(), cwd+outputFile, genomeBuild)
	fmt.Println("### INFO: Compiled shell command", executableCmd.String())

	errCmdRun := executableCmd.Run()
	if errCmdRun != nil {
		fmt.Println("### ERROR", errCmdRun)
		fmt.Fprintln(w, "Oops! Something went wrong on the server side.")
	} else {
		w.Header().Set("Content-Disposition", "attachment; filename="+outputFile)
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.Header().Set("Content-Length", r.Header.Get("Content-Length"))
		cookieObj := cookie.SetDownloadInitiatedCookie()
		http.SetCookie(w, &cookieObj)
		http.ServeFile(w, r, outputFile)
		fmt.Printf("Providing file %s to client...\n", outputFile)
	}
}

func (app *application) mainHandler19to38(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")
	r.ParseMultipartForm(1)
	file, handler, err := r.FormFile("uploadFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)
	tempFile, err := ioutil.TempFile("temp-files", handler.Filename+"_upload*")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	tempFile.Write(fileBytes)

	fmt.Println("### INFO: Temp file written to", tempFile.Name())
	cwd := utils.GetCwd() + "/"
	tempId := utils.GetTempId(tempFile.Name())
	outputFile := "output-files/hg19tohg38." + tempId + ".txt"
	genomeBuild := "hg19"
	executableCmd := utils.MakeCmdStruct(cwd, cwd+tempFile.Name(), cwd+outputFile, genomeBuild)
	fmt.Println("### INFO: Compiled shell command", executableCmd.String())

	errCmdRun := executableCmd.Run()
	if errCmdRun != nil {
		fmt.Println("### ERROR", errCmdRun)
		fmt.Fprintln(w, "Oops! Something went wrong on the server side.")
	} else {
		w.Header().Set("Content-Disposition", "attachment; filename="+outputFile)
		w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
		w.Header().Set("Content-Length", r.Header.Get("Content-Length"))
		cookieObj := cookie.SetDownloadInitiatedCookie()
		http.SetCookie(w, &cookieObj)
		http.ServeFile(w, r, outputFile)
		fmt.Printf("Providing file %s to client...\n", outputFile)
	}
}

func (app *application) basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(app.auth.username))
			expectedPasswordHash := sha256.Sum256([]byte(app.auth.password))
			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)
			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
	})
}
