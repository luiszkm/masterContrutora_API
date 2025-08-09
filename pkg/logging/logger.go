package logging

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

// LogLevel representa os níveis de log
type LogLevel string

const (
	LevelDebug LogLevel = "DEBUG"
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
	LevelFatal LogLevel = "FATAL"
)

// LogEntry representa uma entrada de log estruturada
type LogEntry struct {
	Timestamp   time.Time              `json:"timestamp"`
	Level       LogLevel               `json:"level"`
	Message     string                 `json:"message"`
	Service     string                 `json:"service"`
	Method      string                 `json:"method,omitempty"`
	Path        string                 `json:"path,omitempty"`
	UserID      string                 `json:"userId,omitempty"`
	RequestID   string                 `json:"requestId,omitempty"`
	Duration    *time.Duration         `json:"duration,omitempty"`
	StatusCode  *int                   `json:"statusCode,omitempty"`
	Error       string                 `json:"error,omitempty"`
	StackTrace  string                 `json:"stackTrace,omitempty"`
	Extra       map[string]interface{} `json:"extra,omitempty"`
	File        string                 `json:"file,omitempty"`
	Line        int                    `json:"line,omitempty"`
}

// AppLogger é nosso logger customizado
type AppLogger struct {
	*slog.Logger
	errorFile    *os.File
	auditFile    *os.File
	service      string
	logDirectory string
}

// Config para configuração do logger
type Config struct {
	Service      string
	LogDirectory string
	LogLevel     slog.Level
}

// NewAppLogger cria um novo logger da aplicação
func NewAppLogger(config Config) (*AppLogger, error) {
	if config.LogDirectory == "" {
		config.LogDirectory = "./logs"
	}

	// Criar diretórios se não existirem
	auditDir := filepath.Join(config.LogDirectory, "audit")
	if err := os.MkdirAll(auditDir, 0755); err != nil {
		return nil, fmt.Errorf("erro ao criar diretório de logs: %w", err)
	}

	// Arquivo de erros
	errorFileName := fmt.Sprintf("errors_%s.log", time.Now().Format("2006-01-02"))
	errorFilePath := filepath.Join(config.LogDirectory, errorFileName)
	errorFile, err := os.OpenFile(errorFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir arquivo de erro: %w", err)
	}

	// Arquivo de auditoria
	auditFileName := fmt.Sprintf("audit_%s.log", time.Now().Format("2006-01-02"))
	auditFilePath := filepath.Join(auditDir, auditFileName)
	auditFile, err := os.OpenFile(auditFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		errorFile.Close()
		return nil, fmt.Errorf("erro ao abrir arquivo de auditoria: %w", err)
	}

	// Logger estruturado que escreve para stdout e arquivo
	multiWriter := io.MultiWriter(os.Stdout, errorFile)
	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level: config.LogLevel,
	})

	logger := slog.New(handler)

	return &AppLogger{
		Logger:       logger,
		errorFile:    errorFile,
		auditFile:    auditFile,
		service:      config.Service,
		logDirectory: config.LogDirectory,
	}, nil
}

// Close fecha os arquivos de log
func (l *AppLogger) Close() error {
	var errs []error
	if l.errorFile != nil {
		if err := l.errorFile.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if l.auditFile != nil {
		if err := l.auditFile.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("erros ao fechar arquivos: %v", errs)
	}
	return nil
}

// LogError registra um erro com stack trace
func (l *AppLogger) LogError(ctx context.Context, message string, err error, extra map[string]interface{}) {
	entry := l.createLogEntry(LevelError, message)
	entry.Error = err.Error()
	entry.StackTrace = l.getStackTrace()
	entry.Extra = extra
	
	// Adicionar informações do contexto se disponíveis
	l.addContextInfo(ctx, entry)
	
	l.writeToFile(l.errorFile, entry)
	l.ErrorContext(ctx, message, "error", err, "extra", extra)
}

// LogAudit registra eventos de auditoria
func (l *AppLogger) LogAudit(ctx context.Context, action string, extra map[string]interface{}) {
	entry := l.createLogEntry(LevelInfo, action)
	entry.Extra = extra
	
	// Adicionar informações do contexto
	l.addContextInfo(ctx, entry)
	
	l.writeToFile(l.auditFile, entry)
	l.InfoContext(ctx, action, "extra", extra)
}

// LogHTTPRequest registra requisições HTTP para auditoria
func (l *AppLogger) LogHTTPRequest(ctx context.Context, method, path string, statusCode int, duration time.Duration, userID string, extra map[string]interface{}) {
	entry := l.createLogEntry(LevelInfo, "HTTP Request")
	entry.Method = method
	entry.Path = path
	entry.StatusCode = &statusCode
	entry.Duration = &duration
	entry.UserID = userID
	entry.Extra = extra
	
	l.addContextInfo(ctx, entry)
	
	l.writeToFile(l.auditFile, entry)
	
	// Log para stdout também
	l.InfoContext(ctx, "HTTP Request",
		"method", method,
		"path", path,
		"status", statusCode,
		"duration", duration,
		"userId", userID,
		"extra", extra,
	)
}

// createLogEntry cria uma entrada de log base
func (l *AppLogger) createLogEntry(level LogLevel, message string) *LogEntry {
	_, file, line, _ := runtime.Caller(2)
	
	return &LogEntry{
		Timestamp: time.Now().UTC(),
		Level:     level,
		Message:   message,
		Service:   l.service,
		File:      filepath.Base(file),
		Line:      line,
	}
}

// addContextInfo adiciona informações do contexto à entrada de log
func (l *AppLogger) addContextInfo(ctx context.Context, entry *LogEntry) {
	if ctx == nil {
		return
	}
	
	// Extrair informações do contexto se disponíveis
	if requestID := ctx.Value("requestId"); requestID != nil {
		if id, ok := requestID.(string); ok {
			entry.RequestID = id
		}
	}
	
	if userID := ctx.Value("userId"); userID != nil {
		if id, ok := userID.(string); ok {
			entry.UserID = id
		}
	}
}

// getStackTrace obtém o stack trace atual
func (l *AppLogger) getStackTrace() string {
	buf := make([]byte, 4096)
	n := runtime.Stack(buf, false)
	return string(buf[:n])
}

// writeToFile escreve a entrada de log no arquivo especificado
func (l *AppLogger) writeToFile(file *os.File, entry *LogEntry) {
	if file == nil {
		return
	}
	
	jsonData, err := json.Marshal(entry)
	if err != nil {
		// Fallback para log simples se falhar a serialização
		fmt.Fprintf(file, "%s [%s] %s: %s\n", 
			entry.Timestamp.Format(time.RFC3339), 
			entry.Level, 
			entry.Service, 
			entry.Message)
		return
	}
	
	fmt.Fprintln(file, string(jsonData))
}

// Convenience methods

// Debug registra mensagem de debug
func (l *AppLogger) Debug(ctx context.Context, message string, extra map[string]interface{}) {
	entry := l.createLogEntry(LevelDebug, message)
	entry.Extra = extra
	l.addContextInfo(ctx, entry)
	l.DebugContext(ctx, message, "extra", extra)
}

// Info registra mensagem informativa
func (l *AppLogger) Info(ctx context.Context, message string, extra map[string]interface{}) {
	entry := l.createLogEntry(LevelInfo, message)
	entry.Extra = extra
	l.addContextInfo(ctx, entry)
	l.InfoContext(ctx, message, "extra", extra)
}

// Warn registra aviso
func (l *AppLogger) Warn(ctx context.Context, message string, extra map[string]interface{}) {
	entry := l.createLogEntry(LevelWarn, message)
	entry.Extra = extra
	l.addContextInfo(ctx, entry)
	l.WarnContext(ctx, message, "extra", extra)
}

// Fatal registra erro fatal e termina a aplicação
func (l *AppLogger) Fatal(ctx context.Context, message string, err error, extra map[string]interface{}) {
	entry := l.createLogEntry(LevelFatal, message)
	entry.Error = err.Error()
	entry.StackTrace = l.getStackTrace()
	entry.Extra = extra
	l.addContextInfo(ctx, entry)
	
	l.writeToFile(l.errorFile, entry)
	l.ErrorContext(ctx, message, "error", err, "extra", extra)
	
	os.Exit(1)
}