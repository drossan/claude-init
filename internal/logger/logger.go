// Package logger proporciona un sistema de logging estructurado con niveles.
//
// Este paquete implementa un logger simple y eficiente que soporta
// diferentes niveles de verbosidad: DEBUG, INFO, WARN, ERROR.
package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

// Level representa el nivel de logging.
type Level int

const (
	// DEBUGLevel muestra mensajes de depuración detallados.
	DEBUGLevel Level = iota

	// INFOLevel muestra mensajes informativos generales.
	INFOLevel

	// WARNLevel muestra advertencias.
	WARNLevel

	// ERRORLevel muestra solo errores.
	ERRORLevel
)

// String retorna la representación en string del nivel.
func (l Level) String() string {
	switch l {
	case DEBUGLevel:
		return "DEBUG"
	case INFOLevel:
		return "INFO"
	case WARNLevel:
		return "WARN"
	case ERRORLevel:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger es un logger con soporte para niveles.
type Logger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	level       Level
	out         io.Writer
}

// New crea un nuevo Logger con el nivel especificado.
//
// El logger escribe todos los mensajes a out, pero solo muestra
// aquellos que tengan un nivel mayor o igual al nivel configurado.
func New(out io.Writer, level Level) *Logger {
	if out == nil {
		out = os.Stderr
	}

	flags := log.LstdFlags

	return &Logger{
		debugLogger: log.New(out, "[DEBUG] ", flags),
		infoLogger:  log.New(out, "[INFO] ", flags),
		warnLogger:  log.New(out, "[WARN] ", flags),
		errorLogger: log.New(out, "[ERROR] ", flags),
		level:       level,
		out:         out,
	}
}

// NewDefault crea un nuevo Logger con nivel INFO y salida a stderr.
func NewDefault() *Logger {
	return New(os.Stderr, INFOLevel)
}

// SetLevel cambia el nivel de logging.
func (l *Logger) SetLevel(level Level) {
	l.level = level
}

// Level retorna el nivel de logging actual.
func (l *Logger) Level() Level {
	return l.level
}

// Debug registra un mensaje de nivel DEBUG.
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level <= DEBUGLevel {
		l.debugLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Info registra un mensaje de nivel INFO.
func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= INFOLevel {
		l.infoLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Warn registra un mensaje de nivel WARN.
func (l *Logger) Warn(format string, v ...interface{}) {
	if l.level <= WARNLevel {
		l.warnLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Error registra un mensaje de nivel ERROR.
func (l *Logger) Error(format string, v ...interface{}) {
	if l.level <= ERRORLevel {
		l.errorLogger.Output(2, fmt.Sprintf(format, v...))
	}
}

// Debugf registra un mensaje de nivel DEBUG (alias para Debug).
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Debug(format, v...)
}

// Infof registra un mensaje de nivel INFO (alias para Info).
func (l *Logger) Infof(format string, v ...interface{}) {
	l.Info(format, v...)
}

// Warnf registra un mensaje de nivel WARN (alias para Warn).
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.Warn(format, v...)
}

// Errorf registra un mensaje de nivel ERROR (alias para Error).
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Error(format, v...)
}

// WithLevel retorna un nuevo logger con el nivel especificado.
func (l *Logger) WithLevel(level Level) *Logger {
	return &Logger{
		debugLogger: l.debugLogger,
		infoLogger:  l.infoLogger,
		warnLogger:  l.warnLogger,
		errorLogger: l.errorLogger,
		level:       level,
		out:         l.out,
	}
}

// StdLogger retorna un logger estándar de la librería log que escribe
// a este logger con nivel INFO.
func (l *Logger) StdLogger() *log.Logger {
	return l.infoLogger
}
