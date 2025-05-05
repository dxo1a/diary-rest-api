package routes

import (
	"api-gateway/pkg/config"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// мап префиксов /api/v1/... на нужные сервисы
func Register(r *gin.Engine, cfg *config.Config) {
	api := r.Group("/api/v1")

	// Calendar
	api.GET("/days/:date", ReverseProxy(cfg.CalendarService))
	api.PUT("/days/:date", ReverseProxy(cfg.CalendarService))
	api.GET("/days", ReverseProxy(cfg.CalendarService))

	// Action
	api.GET("/days/:date/actions", ReverseProxy(cfg.ActionService))
	api.POST("/days/:date/actions", ReverseProxy(cfg.ActionService))
	api.GET("/days/:date/actions/:id", ReverseProxy(cfg.ActionService))
	api.PUT("/days/:date/actions/:id", ReverseProxy(cfg.ActionService))
	api.DELETE("/days/:date/actions/:id", ReverseProxy(cfg.ActionService))

	// Habit
	api.GET("/habits", ReverseProxy(cfg.HabitService))
	api.POST("/habits", ReverseProxy(cfg.HabitService))
	api.PUT("/habits/:id", ReverseProxy(cfg.HabitService))
	api.DELETE("/habits/:id", ReverseProxy(cfg.HabitService))
	api.POST("/habits/:id/entries", ReverseProxy(cfg.HabitService))
	api.DELETE("/habits/:id/entries/:date", ReverseProxy(cfg.HabitService))

	// Metrics
	api.GET("/metrics", ReverseProxy(cfg.MetricsService))
	api.GET("/metrics/report", ReverseProxy(cfg.MetricsService))
}

// создает хандлер, который перенаправляет запрос в tagetURL
func ReverseProxy(targetURL string) gin.HandlerFunc {
	target, err := url.Parse(targetURL)
	if err != nil {
		panic("invalid service URL: " + targetURL)
	}
	return func(c *gin.Context) {
		// сборочка URL (база + исходный путь + запрос)
		dest := target.ResolveReference(&url.URL{
			Path:     c.Request.URL.Path,
			RawQuery: c.Request.URL.RawQuery,
		})

		// и наш новенький HTTP запросек
		req, err := http.NewRequest(c.Request.Method, dest.String(), c.Request.Body)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
			return
		}
		// копирование заголовков
		for k, v := range c.Request.Header {
			req.Header[k] = v
		}

		// выполняем запрос и копируем ответ
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		}
		defer resp.Body.Close()

		for k, v := range resp.Header {
			c.Writer.Header()[k] = v
		}
		c.Status(resp.StatusCode)
		io.Copy(c.Writer, resp.Body)
	}
}
