// Package handlers defines handle functions for endpoint listening and processing.
package handlers

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"github.com/kataras/go-sessions/v3"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"os"
	"time"

	"webServerGoLiftover/internal/api/errors"
	"webServerGoLiftover/internal/config"
	"webServerGoLiftover/internal/cookie"
	"webServerGoLiftover/internal/service/liftover"
	"webServerGoLiftover/internal/utils"
)

// URLHandler defines URLHandler object structure.
type URLHandler struct {
	converter    liftover.Converter
	serverConfig *config.ServerConfig
	app          *config.Application
	logger       *zerolog.Logger
}

// InitURLHandler initializes URLHandler object setting its attributes.
func InitURLHandler(serverConfig *config.ServerConfig, app *config.Application, converter liftover.Converter, logger *zerolog.Logger) *URLHandler {
	return &URLHandler{serverConfig: serverConfig, app: app, converter: converter, logger: logger}
}

// MainHandler38to19 handles file upload, processing and provision for hg38-to-hg19 scheme.
func (h *URLHandler) MainHandler38to19(w http.ResponseWriter, r *http.Request) {
	h.logger.Info().Msg("38-to-19 file upload endpoint hit")
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
		h.logger.Warn().Err(err).Msg("Multipart data parsing failed")
		http.Error(w, errors.MultipartParsingError, http.StatusInternalServerError)
		return
	}
	file, header, err := r.FormFile("uploadFile")
	defer file.Close()
	if err != nil {
		h.logger.Warn().Err(err).Msg("Multipart data retrieval failed")
		http.Error(w, errors.MultipartGetFileError, http.StatusInternalServerError)
		return
	}
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		h.logger.Warn().Err(err).Msg("Multipart data reading failed")
		http.Error(w, errors.MultipartReadingError, http.StatusInternalServerError)
		return
	}
	tempFile, err := os.CreateTemp(h.app.Path.UploadedDir, header.Filename+"_upload*")
	defer tempFile.Close()
	defer func(path string) {
		err := utils.RemoveFile(path)
		if err != nil {
			h.logger.Warn().Err(err).Msg("File removal failed")
		}
	}(tempFile.Name())
	if err != nil {
		h.logger.Warn().Err(err).Msg("Temp file opening failed")
		http.Error(w, errors.TempFileOpeningError, http.StatusInternalServerError)
		return
	}
	_, err = tempFile.Write(fileBytes)
	if err != nil {
		h.logger.Warn().Err(err).Msg("Temp file writing failed")
		http.Error(w, errors.TempFileWritingError, http.StatusInternalServerError)
		return
	}
	h.logger.Info().Msg(fmt.Sprintf("Client-provided file was written to", tempFile.Name()))
	fileConformityStatus := utils.CheckUploadedFileConformity(tempFile.Name())
	if fileConformityStatus != "ok" {
		cookieObj := cookie.SetConformityFailedCookie()
		http.SetCookie(w, &cookieObj)
		return
	}
	outputFile := "hg38toHg19." + utils.GetTempId(tempFile.Name()) + ".txt"
	err = h.converter.Convert38to19(h.app.Path.WD, tempFile.Name(), h.app.Path.ProcessedDir+outputFile)
	if err != nil {
		h.logger.Warn().Err(err).Msg("File conversion failed")
		http.Error(w, errors.PyliftoverNonzeroExitError, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+outputFile)
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", utils.GetFileSize(h.app.Path.ProcessedDir+outputFile))
	cookieObj := cookie.SetDownloadInitiatedCookie()
	http.SetCookie(w, &cookieObj)
	http.ServeFile(w, r, h.app.Path.ProcessedDir+outputFile)
	h.logger.Info().Msg(fmt.Sprintf("Successfully provided client with file %s\n", outputFile))
	err = utils.RemoveFile(h.app.Path.ProcessedDir + outputFile)
	if err != nil {
		h.logger.Warn().Err(err).Msg("File removal failed")
	}
}

// MainHandler19to38 handles file upload, processing and provision for hg19-to-hg38 scheme.
func (h *URLHandler) MainHandler19to38(w http.ResponseWriter, r *http.Request) {
	h.logger.Info().Msg("19-to-38 file upload endpoint hit")
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
		h.logger.Warn().Err(err).Msg("Multipart data parsing failed")
		http.Error(w, errors.MultipartParsingError, http.StatusInternalServerError)
		return
	}
	file, header, err := r.FormFile("uploadFile")
	defer file.Close()
	if err != nil {
		h.logger.Warn().Err(err).Msg("Multipart data retrieval failed")
		http.Error(w, errors.MultipartGetFileError, http.StatusInternalServerError)
		return
	}
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		h.logger.Warn().Err(err).Msg("Multipart data reading failed")
		http.Error(w, errors.MultipartReadingError, http.StatusInternalServerError)
		return
	}
	tempFile, err := os.CreateTemp(h.app.Path.UploadedDir, header.Filename+"_upload*")
	defer tempFile.Close()
	defer func(path string) {
		err := utils.RemoveFile(path)
		if err != nil {
			h.logger.Warn().Err(err).Msg("File removal failed")
		}
	}(tempFile.Name())
	if err != nil {
		h.logger.Warn().Err(err).Msg("Temp file opening failed")
		http.Error(w, errors.TempFileOpeningError, http.StatusInternalServerError)
		return
	}
	_, err = tempFile.Write(fileBytes)
	if err != nil {
		h.logger.Warn().Err(err).Msg("Temp file writing failed")
		http.Error(w, errors.TempFileWritingError, http.StatusInternalServerError)
		return
	}
	h.logger.Info().Msg(fmt.Sprintf("Client-provided file was written to", tempFile.Name()))
	fileConformityStatus := utils.CheckUploadedFileConformity(tempFile.Name())
	if fileConformityStatus != "ok" {
		cookieObj := cookie.SetConformityFailedCookie()
		http.SetCookie(w, &cookieObj)
		return
	}
	outputFile := "hg19toHg38." + utils.GetTempId(tempFile.Name()) + ".txt"
	err = h.converter.Convert19to38(h.app.Path.WD, tempFile.Name(), h.app.Path.ProcessedDir+outputFile)
	if err != nil {
		h.logger.Warn().Err(err).Msg("File conversion failed")
		http.Error(w, errors.PyliftoverNonzeroExitError, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+outputFile)
	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", utils.GetFileSize(h.app.Path.ProcessedDir+outputFile))
	cookieObj := cookie.SetDownloadInitiatedCookie()
	http.SetCookie(w, &cookieObj)
	http.ServeFile(w, r, h.app.Path.ProcessedDir+outputFile)
	h.logger.Info().Msg(fmt.Sprintf("Successfully provided client with file %s\n", outputFile))
	err = utils.RemoveFile(h.app.Path.ProcessedDir + outputFile)
	if err != nil {
		h.logger.Warn().Err(err).Msg("File removal failed")
	}
}

// Authorizer provides middleware authorization for file upload endpoints.
func (h *URLHandler) Authorizer(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, "Unauthorized access", http.StatusUnauthorized)
			h.logger.Warn().Msg("Unauthorized access detected")
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
		h.logger.Warn().Msg("Invalid login attempt detected")
	}
}
