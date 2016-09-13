package routes

import (
	"knowledge-base/webActions"
	"net/http"
)

//Route struct
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

//Routes - array of Route
type Routes []Route

var articlesWebAction webActions.ArticlesWebActions
var tagsWebAction webActions.TagsWebActions

//CreateRoutes return all routes
func CreateRoutes() Routes {
	return Routes{
		Route{
			"UrlRoot",
			"GET",
			"/",
			webActions.URLRoot,
		},
		Route{
			"UrlRoot",
			"GET",
			"/knowledge/",
			webActions.URLRoot,
		},
		Route{
			"Read",
			"GET",
			"/knowledge/articles",
			articlesWebAction.Read,
		},
		Route{
			"ReadOne",
			"GET",
			"/knowledge/articles/{id}",
			articlesWebAction.ReadOne,
		},
		Route{
			"Create",
			"POST",
			"/knowledge/articles",
			articlesWebAction.Create,
		},
		Route{
			"Update",
			"PUT",
			"/knowledge/articles/{id}",
			articlesWebAction.Update,
		},
		Route{
			"Update",
			"DELETE",
			"/knowledge/articles/{id}",
			articlesWebAction.Delete,
		},
		Route{
			"AllTags",
			"GET",
			"/knowledge/tags",
			tagsWebAction.All,
		},
	}
}
