package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"

	c "app/controller"
	m "app/route/middleware"
	"app/shared/database"
	"app/shared/jsonconfig"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
)

func init() {
	// Verbose logging with file name and line number
	log.SetFlags(log.Lshortfile)

	// Use all CPU cores
	runtime.GOMAXPROCS(runtime.NumCPU())
}

// UpdateUser editting a user by username. Not allowed changing username
func Test(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	fmt.Printf("jajal")
}

func main() {
	// Load the configuration file
	jsonconfig.Load("config"+string(os.PathSeparator)+"config.json", config)

	// Connect to database
	database.Connect(config.Database)

	// router
	r := httprouter.New()

	// user api
	r.GET("/test", m.AuthUser(Test))
	r.POST("/api/v1/users", m.AnonUser(c.CreateUser))
	r.GET("/api/v1/users/:username", m.AnonUser(c.RetrieveUser))
	r.POST("/api/v1/users/:username", m.AuthUser(c.UpdateUser))
	// r.DELETE("/api/v1/users/:username", m.AuthUser(c.DeleteUser))

	// post api for posting on userpage / feedpage
	r.POST("/api/v1/feeds/:feed/posts", m.AuthUser(c.CreatePostFeed))
	r.GET("/api/v1/feeds/:feed/posts", m.AnonUser(c.GetPostByFeed))
	// r.POST("/api/v1/posts", m.AuthUser(c.CreatePostFeed)) // at default feed (home)
	r.GET("/api/v1/posts/:id", m.AnonUser(c.GetPostByID))
	r.POST("/api/v1/posts/:id", m.AuthUser(c.UpdatePost))
	r.DELETE("/api/v1/posts/:id", m.AuthUser(c.DeletePost))

	// entrypoint
	r.POST("/api/v1/login", m.AnonUser(c.LoginHandler))

	common := negroni.New()
	common.Use(negroni.NewLogger())
	common.Use(negroni.NewRecovery())
	common.UseHandler(r)

	// set cors
	handler := cors.Default().Handler(common)
	err := http.ListenAndServe(":3000", handler)
	if err != nil {
		log.Fatal(err)
	}
}

// *****************************************************************************
// Application Settings
// *****************************************************************************

// config the settings variable
var config = &configuration{}

// configuration contains the application settings
type configuration struct {
	Database database.MongoDB `json:"MongoDB"`
}

// ParseJSON unmarshals bytes to structs
func (c *configuration) ParseJSON(b []byte) error {
	return json.Unmarshal(b, &c)
}
