package logger

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// Logger handles logging to separate files for different metric types
type Logger struct {
	logDir   string
	format   string
	loggers  map[string]*logrus.Logger
	fileHandles map[string]*os.File
	csvWriters map[string]*csv.Writer
	mu       sync.Mutex
}

// NewLogger creates a new logger instance
func NewLogger(logDir, format string) (*Logger, error) {
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create log directory: %w", err)
	}

	l := &Logger{
		logDir:      logDir,
		format:      format,
		loggers:     make(map[string]*logrus.Logger),
		fileHandles: make(map[string]*os.File),
		csvWriters:  make(map[string]*csv.Writer),
	}

	// Initialize loggers for each metric type
	metricTypes := []string{"cpu", "memory", "disk", "network", "game_performance", "steam"}
	for _, metricType := range metricTypes {
		if err := l.initLogger(metricType); err != nil {
			return nil, fmt.Errorf("failed to initialize logger for %s: %w", metricType, err)
		}
	}

	return l, nil
}

func (l *Logger) initLogger(metricType string) error {
	logPath := filepath.Join(l.logDir, fmt.Sprintf("%s.log", metricType))

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	logger := logrus.New()
	logger.SetOutput(file)
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339Nano,
	})

	l.fileHandles[metricType] = file
	l.loggers[metricType] = logger

	if l.format == "csv" {
		writer := csv.NewWriter(file)
		l.csvWriters[metricType] = writer
	}

	return nil
}

// LogCPU logs CPU metrics
func (l *Logger) LogCPU(data interface{}) error {
	return l.log("cpu", data)
}

// LogMemory logs memory metrics
func (l *Logger) LogMemory(data interface{}) error {
	return l.log("memory", data)
}

// LogDisk logs disk metrics
func (l *Logger) LogDisk(data interface{}) error {
	return l.log("disk", data)
}

// LogNetwork logs network metrics
func (l *Logger) LogNetwork(data interface{}) error {
	return l.log("network", data)
}

// LogGamePerformance logs game performance metrics
func (l *Logger) LogGamePerformance(data interface{}) error {
	return l.log("game_performance", data)
}

// LogSteam logs Steam metrics
func (l *Logger) LogSteam(data interface{}) error {
	return l.log("steam", data)
}

func (l *Logger) log(metricType string, data interface{}) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	logger, exists := l.loggers[metricType]
	if !exists {
		return fmt.Errorf("logger for %s not initialized", metricType)
	}

	if l.format == "csv" {
		return l.logCSV(metricType, data)
	}

	// JSON format
	logger.WithFields(logrus.Fields{
		"metric": data,
	}).Info("metric")

	return nil
}

func (l *Logger) logCSV(metricType string, data interface{}) error {
	writer, exists := l.csvWriters[metricType]
	if !exists {
		return fmt.Errorf("CSV writer for %s not initialized", metricType)
	}

	// Convert data to JSON first, then parse for CSV
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	var dataMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &dataMap); err != nil {
		return fmt.Errorf("failed to unmarshal data: %w", err)
	}

	// Extract timestamp
	timestamp := time.Now().Format(time.RFC3339Nano)
	if ts, ok := dataMap["timestamp"]; ok {
		if tsStr, ok := ts.(string); ok {
			timestamp = tsStr
		}
	}

	// Create CSV row
	row := []string{timestamp}
	for key, value := range dataMap {
		if key != "timestamp" {
			row = append(row, fmt.Sprintf("%v", value))
		}
	}

	if err := writer.Write(row); err != nil {
		return fmt.Errorf("failed to write CSV row: %w", err)
	}

	writer.Flush()
	return writer.Error()
}

// Close closes all log file handles
func (l *Logger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	for metricType, file := range l.fileHandles {
		if writer, exists := l.csvWriters[metricType]; exists {
			writer.Flush()
		}
		if err := file.Close(); err != nil {
			return fmt.Errorf("failed to close log file for %s: %w", metricType, err)
		}
	}

	return nil
}

