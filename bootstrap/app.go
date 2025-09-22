// bootstrap/app.go
package bootstrap

import (
	"database/sql"
	"fmt"
	"log"

	_ "your_project/app/models"

	"your_project/app/providers"

	"github.com/aasoft24/golara/wpkg/cache"
	"github.com/aasoft24/golara/wpkg/configs"
	"github.com/aasoft24/golara/wpkg/database"
	"github.com/aasoft24/golara/wpkg/foundation"
	"github.com/aasoft24/golara/wpkg/gola"
	"github.com/aasoft24/golara/wpkg/routing"
	"github.com/aasoft24/golara/wpkg/session"
	"github.com/aasoft24/golara/wpkg/view"
	"gorm.io/gorm"
)

// Init returns router so main.go can use it
func Init() *routing.Router {
	// 1️⃣ Load config
	configs.Init()

	// 2️⃣ Init DB safely
	if err := database.InitDB(); err != nil {
		fmt.Println(err) // ❌ show error but continue
	}

	// 3️⃣ Init redis
	cache.InitRedis()

	// 4️⃣ Template engine
	templateEngine := view.NewTemplateEngine("resources/views", "app")
	ctx := &gola.Context{TemplateEngine: templateEngine}
	router := routing.NewRouter(ctx)

	// 5️⃣ Session middleware
	store := session.NewMemoryStore()
	sessionManager := session.NewManager(store, "go_session")

	router.Use(func(next func(ctx *gola.Context)) func(ctx *gola.Context) {
		return func(ctx *gola.Context) {
			sess, _ := sessionManager.Start(ctx.Writer, ctx.Request)
			ctx.Session = sess
			ctx.SessionManager = sessionManager
			next(ctx)
			_ = sess.Save()
		}
	})

	// 6️⃣ Logging middleware
	router.Use(func(next func(ctx *gola.Context)) func(ctx *gola.Context) {
		return func(ctx *gola.Context) {
			log.Printf("%s %s", ctx.Request.Method, ctx.Request.URL.Path)
			next(ctx)
		}
	})

	// 7️⃣ Cache
	appCache := cache.NewMemoryCache()

	// 8️⃣ Application container
	app := foundation.NewApplication()

	// Bind GORM *gorm.DB instead of *sql.DB
	app.Bind((*gorm.DB)(nil), database.DB)

	// If you still want *sql.DB:
	if database.DB != nil {
		if sqlDB, err := database.DB.DB(); err == nil {
			app.Bind((*sql.DB)(nil), sqlDB)
		}
	}

	app.Bind((*view.TemplateEngine)(nil), templateEngine)
	app.Bind((*cache.Cache)(nil), appCache)

	// 9️⃣ Register route service provider
	app.Register(providers.NewRouteServiceProvider(router, templateEngine))
	app.Boot()

	return router
}
