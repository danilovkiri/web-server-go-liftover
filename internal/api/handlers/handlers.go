package handlers

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"github.com/kataras/go-sessions/v3"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"webServerGoLiftover/internal/api/errors"
	"webServerGoLiftover/internal/config"
	"webServerGoLiftover/internal/cookie"
	"webServerGoLiftover/internal/utils"
)

type URLHandler struct {
	serverConfig *config.ServerConfig
	app          *config.Application
}

func InitURLHandler(serverConfig *config.ServerConfig, app *config.Application) *URLHandler {
	return &URLHandler{serverConfig: serverConfig, app: app}
}

func (h *URLHandler) MainHandler38to19(w http.ResponseWriter, r *http.Request) {
	log.Println("### INFO: File Upload Endpoint Hit")
	sess := sessions.New(sessions.Config{
		Cookie:                      "sessionCookie",
		Expires:                     time.Minute * 1,
		DisableSubdomainPersistence: false,
	})
	s := sess.Start(w, r)
	s.Set("name", "authorizedClientSession")
	defer s.Destroy()
	err := r.ParseMultipartForm(1)
	if err != nil {
		log.Println(err)
		http.Error(w, errors.MultipartParsingError, http.StatusInternalServerError)
		return
	}
	file, header, err := r.FormFile("uploadFile")
	defer file.Close()
	if err != nil {
		log.Println(err)
		http.Error(w, errors.MultipartGetFileError, http.StatusInternalServerError)
		return
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		http.Error(w, errors.MultipartReadingError, http.StatusInternalServerError)
		return
	}
	tempFile, err := ioutil.TempFile(h.app.Path.UploadedDir, header.Filename+"_upload*")
	defer tempFile.Close()
	defer utils.RemoveFile(tempFile.Name())
	if err != nil {
		log.Println(err)
		http.Error(w, errors.TempFileOpeningError, http.StatusInternalServerError)
		return
	}
	_, err = tempFile.Write(fileBytes)
	if err != nil {
		log.Println(err)
		http.Error(w, errors.TempFileWritingError, http.StatusInternalServerError)
		return
	}
	log.Println("### INFO: Client-provided file was written to", tempFile.Name())
	fileConformityStatus := utils.CheckUploadedFileConformity(tempFile.Name())
	if fileConformityStatus != "ok" {
		log.Println(fileConformityStatus)
		cookieObj := cookie.SetConformityFailedCookie()
		http.SetCookie(w, &cookieObj)
		return
	}
	tempId := utils.GetTempId(tempFile.Name())
	outputFile := "hg38toHg19." + tempId + ".txt"
	executableCmd := utils.MakeCmdStruct(h.app.Path.Cwd, tempFile.Name(), h.app.Path.ProcessedDir+outputFile, "hg38")
	log.Println("Compiled shell command:", executableCmd.String())
	err = executableCmd.Run()
	if err != nil {
		log.Println(err)
		http.Error(w, errors.PyliftoverNonzeroExitError, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+outputFile)
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", utils.GetFileSize(h.app.Path.ProcessedDir+outputFile))
	cookieObj := cookie.SetDownloadInitiatedCookie()
	http.SetCookie(w, &cookieObj)
	http.ServeFile(w, r, h.app.Path.ProcessedDir+outputFile)
	fmt.Printf("Successfully provided client with file %s\n", outputFile)
	utils.RemoveFile(h.app.Path.ProcessedDir + outputFile)
}

func (h *URLHandler) MainHandler19to38(w http.ResponseWriter, r *http.Request) {
	log.Println("### INFO: File Upload Endpoint Hit")
	sess := sessions.New(sessions.Config{
		Cookie:                      "sessionCookie",
		Expires:                     time.Minute * 1,
		DisableSubdomainPersistence: false,
	})
	s := sess.Start(w, r)
	s.Set("name", "authorizedClientSession")
	defer s.Destroy()
	err := r.ParseMultipartForm(1)
	if err != nil {
		log.Println(err)
		http.Error(w, errors.MultipartParsingError, http.StatusInternalServerError)
		return
	}
	file, header, err := r.FormFile("uploadFile")
	defer file.Close()
	if err != nil {
		log.Println(err)
		http.Error(w, errors.MultipartGetFileError, http.StatusInternalServerError)
		return
	}
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		http.Error(w, errors.MultipartReadingError, http.StatusInternalServerError)
		return
	}
	tempFile, err := ioutil.TempFile(h.app.Path.UploadedDir, header.Filename+"_upload*")
	defer tempFile.Close()
	defer utils.RemoveFile(tempFile.Name())
	if err != nil {
		log.Println(err)
		http.Error(w, errors.TempFileOpeningError, http.StatusInternalServerError)
		return
	}
	_, err = tempFile.Write(fileBytes)
	if err != nil {
		log.Println(err)
		http.Error(w, errors.TempFileWritingError, http.StatusInternalServerError)
		return
	}
	log.Println("### INFO: Client-provided file was written to", tempFile.Name())
	fileConformityStatus := utils.CheckUploadedFileConformity(tempFile.Name())
	if fileConformityStatus != "ok" {
		log.Println(fileConformityStatus)
		cookieObj := cookie.SetConformityFailedCookie()
		http.SetCookie(w, &cookieObj)
		return
	}
	tempId := utils.GetTempId(tempFile.Name())
	outputFile := "hg19toHg38." + tempId + ".txt"
	executableCmd := utils.MakeCmdStruct(h.app.Path.Cwd, tempFile.Name(), h.app.Path.ProcessedDir+outputFile, "hg19")
	log.Println("Compiled shell command:", executableCmd.String())
	err = executableCmd.Run()
	if err != nil {
		log.Println(err)
		http.Error(w, errors.PyliftoverNonzeroExitError, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+outputFile)
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", utils.GetFileSize(h.app.Path.ProcessedDir+outputFile))
	cookieObj := cookie.SetDownloadInitiatedCookie()
	http.SetCookie(w, &cookieObj)
	http.ServeFile(w, r, h.app.Path.ProcessedDir+outputFile)
	fmt.Printf("Successfully provided client with file %s\n", outputFile)
	utils.RemoveFile(h.app.Path.ProcessedDir + outputFile)
}

func (h *URLHandler) Authorizer(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		}
		usernameHash := sha256.Sum256([]byte(username))
		passwordHash := sha256.Sum256([]byte(password))
		expectedUsernameHash := sha256.Sum256([]byte(h.app.Auth.Username))
		expectedPasswordHash := sha256.Sum256([]byte(h.app.Auth.Password))
		usernameMatch := subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1
		passwordMatch := subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1
		if usernameMatch && passwordMatch {
			next.ServeHTTP(w, r)
			return
		}
	}
}
