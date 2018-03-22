package logger

import (
	"os"

	"github.com/op/go-logging"
)

// Password is a redactable string type
type Password string

// Redacted is an method of Password to implement the Redactor interface
func (p Password) Redacted() interface{} {
	return logging.Redact(string(p))
}

// New returns a new Logger instance
func New() *logging.Logger {
	logger := logging.MustGetLogger("server")

	format := logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfunc} >> %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)

	stdout := logging.NewLogBackend(os.Stdout, "", 0)
	stdoutFormatted := logging.NewBackendFormatter(stdout, format)

	logging.SetBackend(stdoutFormatted)

	return logger
}
