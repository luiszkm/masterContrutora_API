package logging

import (
	"context"
	"fmt"
	"time"
)

// DashboardLogger fornece métodos específicos para logging do dashboard
type DashboardLogger struct {
	*AppLogger
}

// NewDashboardLogger cria um logger específico para o dashboard
func NewDashboardLogger(appLogger *AppLogger) *DashboardLogger {
	return &DashboardLogger{
		AppLogger: appLogger,
	}
}

// LogDashboardError registra erros específicos do dashboard
func (d *DashboardLogger) LogDashboardError(ctx context.Context, section string, operation string, err error, extra map[string]interface{}) {
	if extra == nil {
		extra = make(map[string]interface{})
	}

	extra["dashboardSection"] = section
	extra["dashboardOperation"] = operation
	extra["component"] = "dashboard"

	message := fmt.Sprintf("Dashboard error in section '%s' during operation '%s'", section, operation)
	d.LogError(ctx, message, err, extra)
}

// LogDashboardQuery registra queries do dashboard para debug
func (d *DashboardLogger) LogDashboardQuery(ctx context.Context, section string, queryName string, duration time.Duration, rowCount int, extra map[string]interface{}) {
	if extra == nil {
		extra = make(map[string]interface{})
	}

	extra["dashboardSection"] = section
	extra["queryName"] = queryName
	extra["queryDuration"] = duration.String()
	extra["rowCount"] = rowCount
	extra["component"] = "dashboard"
	extra["queryPerformance"] = true

	message := fmt.Sprintf("Dashboard query '%s' in section '%s' completed", queryName, section)
	d.Info(ctx, message, extra)
}

// LogDashboardServiceCall registra chamadas de serviço do dashboard
func (d *DashboardLogger) LogDashboardServiceCall(ctx context.Context, section string, method string, startTime time.Time, err error, extra map[string]interface{}) {
	duration := time.Since(startTime)

	if extra == nil {
		extra = make(map[string]interface{})
	}

	extra["dashboardSection"] = section
	extra["serviceMethod"] = method
	extra["serviceDuration"] = duration.String()
	extra["component"] = "dashboard"

	if err != nil {
		extra["serviceError"] = true
		message := fmt.Sprintf("Dashboard service call failed: %s.%s", section, method)
		d.LogError(ctx, message, err, extra)
	} else {
		message := fmt.Sprintf("Dashboard service call completed: %s.%s", section, method)
		d.Info(ctx, message, extra)
	}
}

// LogDashboardPerformance registra métricas de performance do dashboard
func (d *DashboardLogger) LogDashboardPerformance(ctx context.Context, section string, totalDuration time.Duration, queryCount int, extra map[string]interface{}) {
	if extra == nil {
		extra = make(map[string]interface{})
	}

	extra["dashboardSection"] = section
	extra["totalDuration"] = totalDuration.String()
	extra["queryCount"] = queryCount
	extra["avgQueryTime"] = (totalDuration / time.Duration(queryCount)).String()
	extra["component"] = "dashboard"
	extra["performance"] = true

	message := fmt.Sprintf("Dashboard section '%s' performance metrics", section)
	d.Info(ctx, message, extra)
}

// LogDashboardCache registra eventos de cache do dashboard
func (d *DashboardLogger) LogDashboardCache(ctx context.Context, section string, action string, hit bool, extra map[string]interface{}) {
	if extra == nil {
		extra = make(map[string]interface{})
	}

	extra["dashboardSection"] = section
	extra["cacheAction"] = action
	extra["cacheHit"] = hit
	extra["component"] = "dashboard"
	extra["cache"] = true

	message := fmt.Sprintf("Dashboard cache %s for section '%s' - hit: %v", action, section, hit)
	d.Info(ctx, message, extra)
}

// LogDashboardAuth registra eventos de autorização do dashboard
func (d *DashboardLogger) LogDashboardAuth(ctx context.Context, section string, userID string, permission string, allowed bool, extra map[string]interface{}) {
	if extra == nil {
		extra = make(map[string]interface{})
	}

	extra["dashboardSection"] = section
	extra["userId"] = userID
	extra["permission"] = permission
	extra["allowed"] = allowed
	extra["component"] = "dashboard"
	extra["authorization"] = true

	var message string
	if allowed {
		message = fmt.Sprintf("Dashboard access granted to section '%s' for user '%s'", section, userID)
		d.Info(ctx, message, extra)
	} else {
		message = fmt.Sprintf("Dashboard access denied to section '%s' for user '%s' - missing permission '%s'", section, userID, permission)
		d.Warn(ctx, message, extra)
	}
}

// LogDashboardValidation registra erros de validação de parâmetros
func (d *DashboardLogger) LogDashboardValidation(ctx context.Context, section string, parameter string, value interface{}, validationError string, extra map[string]interface{}) {
	if extra == nil {
		extra = make(map[string]interface{})
	}

	extra["dashboardSection"] = section
	extra["parameter"] = parameter
	extra["parameterValue"] = value
	extra["validationError"] = validationError
	extra["component"] = "dashboard"
	extra["validation"] = true

	message := fmt.Sprintf("Dashboard validation error in section '%s' for parameter '%s'", section, parameter)
	d.Warn(ctx, message, extra)
}

// LogDashboardData registra informações sobre os dados retornados
func (d *DashboardLogger) LogDashboardData(ctx context.Context, section string, dataType string, recordCount int, isEmpty bool, extra map[string]interface{}) {
	if extra == nil {
		extra = make(map[string]interface{})
	}

	extra["dashboardSection"] = section
	extra["dataType"] = dataType
	extra["recordCount"] = recordCount
	extra["isEmpty"] = isEmpty
	extra["component"] = "dashboard"
	extra["dataMetrics"] = true

	message := fmt.Sprintf("Dashboard data for section '%s' - type: %s, records: %d", section, dataType, recordCount)
	d.Debug(ctx, message, extra)
}
