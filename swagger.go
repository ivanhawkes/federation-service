package main

import (
	"appengine"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
)

func gaeUrl() string {
	if appengine.IsDevAppServer() {
		return "http://localhost:8080"
	} else {
		// Include your URL on App Engine here.
		// I found no way to get AppID without appengine.Context and this always
		// based on a http.Request.
		return "http://federation.shatteredscreens.com/"
	}
}

func init() {
	// Optionally, you can install the Swagger Service which provides a nice Web UI on your REST API
	// You need to download the Swagger HTML5 assets and change the FilePath location in the config below.
	// Open <your_app_id>.appspot.com/apidocs and enter
	// http://<your_app_id>.appspot.com/apidocs.json in the api input field.
	config := swagger.Config{
		// You control what services are visible
		WebServices:    restful.RegisteredWebServices(),
		WebServicesUrl: gaeUrl(),
		ApiPath:        "/apidocs.json",

		// Optionally, specifiy where the UI is located
		SwaggerPath: "/apidocs/",

		// GAE support static content which is configured in your app.yaml.
		// This example expect the swagger-ui in static/swagger so you should place it there :)
		SwaggerFilePath: "static/swagger"}
	swagger.InstallSwaggerService(config)
}
