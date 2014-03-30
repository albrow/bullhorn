package app

import (
	"bullhorn/app/controllers"
	"bullhorn/app/models"
	"github.com/albrow/zoom"
	"github.com/garyburd/redigo/redis"
	"github.com/revel/revel"
	"log"
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.ActionInvoker,           // Invoke the action.
	}

	// register startup functions with OnAppStart
	// ( order dependent )
	// revel.OnAppStart(InitDB())
	// revel.OnAppStart(FillCache())
	revel.OnAppStart(func() {
		// get config vars
		host, found := revel.Config.String("zoom.host")
		if !found {
			log.Fatalf("Missing required config: zoom.host")
		}
		port, found := revel.Config.String("zoom.port")
		if !found {
			log.Fatalf("Missing required config: zoom.port")
		}

		// init zoom
		config := &zoom.Configuration{
			Address: (host + ":" + port),
		}
		zoom.Init(config)

		// test the connection
		conn := zoom.GetConn()
		defer conn.Close()
		if reply, err := redis.String(conn.Do("PING")); err != nil {
			log.Fatalf("Could not connect to database: %s", err)
		} else {
			if reply != "PONG" {
				log.Fatalf("Incorrect database reply. Expected PONG but got: %s", err)
			}
		}

		// register models
		if err := models.Init(); err != nil {
			log.Fatalf("Problem registering models: %s", err)
		}

		// init what we need in controllers
		controllers.Init()
	})
}

// TODO turn this into revel.HeaderFilter
// should probably also have a filter for CSRF
// not sure if it can go in the same filter or not
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	// Add some common security headers
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}
