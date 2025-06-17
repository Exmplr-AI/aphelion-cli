package api

import "time"

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Service struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Spec        map[string]interface{} `json:"spec"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type ServicesResponse struct {
	Services []Service `json:"services"`
	Total    int       `json:"total"`
}

type Memory struct {
	ID          string                 `json:"id"`
	SessionID   string                 `json:"session_id"`
	Summary     string                 `json:"summary"`
	Content     map[string]interface{} `json:"content"`
	CreatedAt   time.Time              `json:"created_at"`
	Similarity  float64                `json:"similarity,omitempty"`
}

type MemoriesResponse struct {
	Memories []Memory `json:"memories"`
	Total    int      `json:"total"`
	Cursor   string   `json:"cursor,omitempty"`
}

type MemoryStats struct {
	TotalMemories    int     `json:"total_memories"`
	TotalSessions    int     `json:"total_sessions"`
	AveragePerDay    float64 `json:"average_per_day"`
	OldestMemory     string  `json:"oldest_memory"`
	MostRecentMemory string  `json:"most_recent_memory"`
}

type Analytics struct {
	RequestMetrics RequestMetrics `json:"request_metrics"`
	UserMetrics    UserMetrics    `json:"user_metrics"`
	ToolMetrics    ToolMetrics    `json:"tool_metrics"`
	SessionMetrics SessionMetrics `json:"session_metrics"`
}

type RequestMetrics struct {
	TotalRequests    int     `json:"total_requests"`
	SuccessfulCount  int     `json:"successful_count"`
	ErrorCount       int     `json:"error_count"`
	AverageTime      float64 `json:"average_time"`
	SuccessRate      float64 `json:"success_rate"`
}

type UserMetrics struct {
	UniqueUsers   int `json:"unique_users"`
	ActiveUsers   int `json:"active_users"`
	NewUsers      int `json:"new_users"`
	ReturningUsers int `json:"returning_users"`
}

type ToolMetrics struct {
	TotalExecutions int         `json:"total_executions"`
	UniqueTools     int         `json:"unique_tools"`
	PopularTools    []ToolUsage `json:"popular_tools"`
}

type ToolUsage struct {
	Tool        string  `json:"tool"`
	Count       int     `json:"count"`
	SuccessRate float64 `json:"success_rate"`
}

type SessionMetrics struct {
	TotalSessions     int     `json:"total_sessions"`
	ActiveSessions    int     `json:"active_sessions"`
	AverageActivities float64 `json:"average_activities"`
	AverageDuration   float64 `json:"average_duration"`
}