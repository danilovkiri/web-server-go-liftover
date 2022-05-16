package errors

const (
	MultipartParsingError      = "Oops! Could not parse multipart form data"
	MultipartGetFileError      = "Oops! Could not get a file from the provided multipart form data"
	MultipartReadingError      = "Oops! Could not read multipart form data"
	TempFileOpeningError       = "Oops! Could not open a temporary file to store upload"
	TempFileWritingError       = "Oops! Could not dump uploaded multipart form data into a temporary file"
	PyliftoverNonzeroExitError = "Oops! Pyliftover returned a non-zero exit code"
)
