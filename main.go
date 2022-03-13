package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"github.com/kataras/go-sessions/v3"
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
	path struct {
		uploadedDir  string
		processedDir string
		cwd          string
	}
}

func main() {
	app := new(application)
	app.auth.username = os.Getenv("AUTH_USERNAME")
	app.auth.password = os.Getenv("AUTH_PASSWORD")
	app.auth.configpath = os.Getenv("CONFIG")
	yamlConfig := config.ParseConfig(app.auth.configpath)
	app.path.uploadedDir = yamlConfig.Constants.FileStorage + "/uploaded-files/"
	app.path.processedDir = yamlConfig.Constants.FileStorage + "/processed-files/"
	app.path.cwd = utils.GetCwd()

	if app.auth.username == "" {
		log.Fatal("### ERROR: Basic auth username must be provided")
	}
	if app.auth.password == "" {
		log.Fatal("### ERROR: Basic auth password must be provided")
	}

	mux := http.NewServeMux()
	mux.Handle("/index/", http.StripPrefix("/index/", http.FileServer(http.Dir("./public"))))
	mux.HandleFunc("/index/hg19to38/", app.basicAuth(app.mainHandler19to38))
	mux.HandleFunc("/index/hg38to19/", app.basicAuth(app.mainHandler38to19))

	srv := &http.Server{
		Addr:         fmt.Sprintf("%v:%v", yamlConfig.Constants.ServerIP, yamlConfig.Constants.ServerPort),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	utils.MakeDir(yamlConfig.Constants.FileStorage)
	utils.MakeDir(app.path.uploadedDir)
	utils.MakeDir(app.path.processedDir)
	log.Printf("### INFO: Using CWD %s", app.path.cwd)
	log.Printf("### INFO: Starting server on %s", srv.Addr)
	err := srv.ListenAndServeTLS(yamlConfig.Constants.CertFile, yamlConfig.Constants.KeyFile)
	log.Fatal("### ERROR: ", err)
}

func (app *application) mainHandler38to19(w http.ResponseWriter, r *http.Request) {
	sess := sessions.New(sessions.Config{
		Cookie:                      "sessionCookie",
		Expires:                     time.Minute * 1,
		DisableSubdomainPersistence: false,
	})
	s := sess.Start(w, r)
	s.Set("name", "authorizedClientSession")
	fmt.Println("### INFO: File Upload Endpoint Hit")
	err := r.ParseMultipartForm(1)
	if err != nil {
		fmt.Println("### ERROR", err)
		http.Error(w, "Oops! Could not parse multipart form data", http.StatusInternalServerError)
		s.Destroy()
		return
	}
	file, handler, err := r.FormFile("uploadFile")
	if err != nil {
		fmt.Println("### ERROR", err)
		http.Error(w, "Oops! Could not get a file from the provided multipart form data", http.StatusInternalServerError)
		s.Destroy()
		return
	}
	defer file.Close()
	fmt.Printf("--- INFO: Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("--- INFO: File Size: %+v\n", handler.Size)
	fmt.Printf("--- INFO: MIME Header: %+v\n", handler.Header)
	tempFile, err := ioutil.TempFile(app.path.uploadedDir, handler.Filename+"_upload*")
	if err != nil {
		fmt.Println("### ERROR:", err)
		http.Error(w, "Oops! Could not open a temporary file to store upload", http.StatusInternalServerError)
		s.Destroy()
		return
	}
	defer tempFile.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("### ERROR:", err)
		http.Error(w, "Oops! Could not read multipart form data", http.StatusInternalServerError)
		s.Destroy()
		return
	}
	_, err1 := tempFile.Write(fileBytes)
	if err1 != nil {
		fmt.Println("### ERROR:", err1)
		http.Error(w, "Oops! Could not dump uploaded multipart form data into a temporary file", http.StatusInternalServerError)
		s.Destroy()
		return
	}
	fmt.Println("### INFO: Client-provided file was written to", tempFile.Name())
	fileConformityStatus := utils.CheckUploadedFileConformity(tempFile.Name())
	if fileConformityStatus != "ok" {
		fmt.Println("### WARNING:", fileConformityStatus)
		cookieObj := cookie.SetConformityFailedCookie()
		http.SetCookie(w, &cookieObj)
		utils.RemoveFile(tempFile.Name())
		s.Destroy()
		return
	}
	tempId := utils.GetTempId(tempFile.Name())
	outputFile := "hg38toHg19." + tempId + ".txt"
	executableCmd := utils.MakeCmdStruct(app.path.cwd, tempFile.Name(), app.path.processedDir+outputFile, "hg38")
	fmt.Println("### INFO: Compiled shell command", executableCmd.String())
	errCmdRun := executableCmd.Run()
	if errCmdRun != nil {
		fmt.Println("### ERROR", errCmdRun)
		http.Error(w, "Oops! Pyliftover returned a non-zero exit code", http.StatusInternalServerError)
		utils.RemoveFile(tempFile.Name())
		s.Destroy()
		return
	} else {
		w.Header().Set("Content-Disposition", "attachment; filename="+outputFile)
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", utils.GetFileSize(app.path.processedDir+outputFile))
		cookieObj := cookie.SetDownloadInitiatedCookie()
		http.SetCookie(w, &cookieObj)
		http.ServeFile(w, r, app.path.processedDir+outputFile)
		fmt.Printf("Successfully provided client with file %s\n", outputFile)
		utils.RemoveFile(tempFile.Name())
		utils.RemoveFile(app.path.processedDir + outputFile)
		s.Destroy()
	}
}

func (app *application) mainHandler19to38(w http.ResponseWriter, r *http.Request) {
	sess := sessions.New(sessions.Config{
		Cookie:                      "sessionCookie",
		Expires:                     time.Minute * 1,
		DisableSubdomainPersistence: false,
	})
	s := sess.Start(w, r)
	s.Set("name", "authorizedClientSession")
	fmt.Println("### INFO: File Upload Endpoint Hit")
	err := r.ParseMultipartForm(1)
	if err != nil {
		fmt.Println("### ERROR", err)
		http.Error(w, "Oops! Could not parse multipart form data", http.StatusInternalServerError)
		s.Destroy()
		return
	}
	file, handler, err := r.FormFile("uploadFile")
	if err != nil {
		fmt.Println("### ERROR", err)
		http.Error(w, "Oops! Could not get a file from the provided multipart form data", http.StatusInternalServerError)
		s.Destroy()
		return
	}
	defer file.Close()
	fmt.Printf("--- INFO: Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("--- INFO: File Size: %+v\n", handler.Size)
	fmt.Printf("--- INFO: MIME Header: %+v\n", handler.Header)
	tempFile, err := ioutil.TempFile(app.path.uploadedDir, handler.Filename+"_upload*")
	if err != nil {
		fmt.Println("### ERROR:", err)
		http.Error(w, "Oops! Could not open a temporary file to store upload", http.StatusInternalServerError)
		s.Destroy()
		return
	}
	defer tempFile.Close()
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("### ERROR:", err)
		http.Error(w, "Oops! Could not read multipart form data", http.StatusInternalServerError)
		s.Destroy()
		return
	}
	_, err1 := tempFile.Write(fileBytes)
	if err1 != nil {
		fmt.Println("### ERROR:", err1)
		http.Error(w, "Oops! Could not dump uploaded multipart form data into a temporary file", http.StatusInternalServerError)
		s.Destroy()
		return
	}
	fmt.Println("### INFO: Client-provided file was written to", tempFile.Name())
	fileConformityStatus := utils.CheckUploadedFileConformity(tempFile.Name())
	if fileConformityStatus != "ok" {
		fmt.Println("### WARNING:", fileConformityStatus)
		cookieObj := cookie.SetConformityFailedCookie()
		http.SetCookie(w, &cookieObj)
		utils.RemoveFile(tempFile.Name())
		s.Destroy()
		return
	}
	tempId := utils.GetTempId(tempFile.Name())
	outputFile := "hg19toHg38." + tempId + ".txt"
	executableCmd := utils.MakeCmdStruct(app.path.cwd, tempFile.Name(), app.path.processedDir+outputFile, "hg19")
	fmt.Println("### INFO: Compiled shell command", executableCmd.String())
	errCmdRun := executableCmd.Run()
	if errCmdRun != nil {
		fmt.Println("### ERROR", errCmdRun)
		http.Error(w, "Oops! Pyliftover returned a non-zero exit code", http.StatusInternalServerError)
		utils.RemoveFile(tempFile.Name())
		s.Destroy()
		return
	} else {
		w.Header().Set("Content-Disposition", "attachment; filename="+outputFile)
		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Content-Length", utils.GetFileSize(app.path.processedDir+outputFile))
		cookieObj := cookie.SetDownloadInitiatedCookie()
		http.SetCookie(w, &cookieObj)
		http.ServeFile(w, r, app.path.processedDir+outputFile)
		fmt.Printf("Successfully provided client with file %s\n", outputFile)
		utils.RemoveFile(tempFile.Name())
		utils.RemoveFile(app.path.processedDir + outputFile)
		s.Destroy()
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
