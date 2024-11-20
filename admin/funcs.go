package admin

import (
	"log/slog"
	"net/http"

	urlpkg "net/url"

	"github.com/gouniverse/hb"
	"github.com/gouniverse/taskstore"
	"github.com/samber/lo"
	"github.com/spf13/cast"
)

func adminHeader(store taskstore.StoreInterface, logger *slog.Logger, r *http.Request) hb.TagInterface {
	linkHome := hb.NewHyperlink().
		HTML("Dashboard").
		Href(url(r, pathQueueManager, nil)).
		Class("nav-link")
	linkQueue := hb.Hyperlink().
		HTML("Queue").
		Href(url(r, pathQueueManager, nil)).
		Class("nav-link")
	linkTasks := hb.Hyperlink().
		HTML("Tasks").
		Href(url(r, pathTaskManager, nil)).
		Class("nav-link")

	queueCount, err := store.QueueCount(taskstore.QueueQuery())

	if err != nil {
		logger.Error(err.Error())
		queueCount = -1
	}

	taskCount, err := store.TaskCount(taskstore.TaskQuery())

	if err != nil {
		logger.Error(err.Error())
		taskCount = -1
	}

	ulNav := hb.NewUL().Class("nav  nav-pills justify-content-center")
	ulNav.AddChild(hb.NewLI().Class("nav-item").Child(linkHome))

	ulNav.Child(hb.LI().
		Class("nav-item").
		Child(linkQueue.
			Child(hb.Span().
				Class("badge bg-secondary").
				HTML(cast.ToString(queueCount)))))

	ulNav.Child(hb.LI().
		Class("nav-item").
		Child(linkTasks.
			Child(hb.Span().
				Class("badge bg-secondary").
				HTML(cast.ToString(taskCount)))))

	divCard := hb.NewDiv().Class("card card-default mt-3 mb-3")
	divCardBody := hb.NewDiv().Class("card-body").Style("padding: 2px;")
	return divCard.AddChild(divCardBody.AddChild(ulNav))
}

func breadcrumbs(r *http.Request, pageBreadcrumbs []Breadcrumb) hb.TagInterface {
	adminHomeURL := "/admin" //AdminHomeURL(r)
	//path := utils.Req(r, "path", "")

	adminHomeBreadcrumb := lo.
		If(adminHomeURL != "", Breadcrumb{
			Name: "Home",
			URL:  adminHomeURL,
		}).
		Else(Breadcrumb{})

	breadcrumbItems := []Breadcrumb{
		adminHomeBreadcrumb,
		{
			Name: "Uptime",
			URL:  url(r, "", nil),
		},
	}

	breadcrumbItems = append(breadcrumbItems, pageBreadcrumbs...)

	breadcrumbs := breadcrumbsUI(breadcrumbItems)

	return hb.Div().
		Child(breadcrumbs)
}

func redirect(w http.ResponseWriter, r *http.Request, url string) string {
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
	// http.Redirect(w, r, url, http.StatusSeeOther)
	return ""
}

func url(r *http.Request, path string, params map[string]string) string {
	endpoint := r.URL.Path

	if params == nil {
		params = map[string]string{}
	}

	params["controller"] = path

	url := endpoint + query(params)

	return url
}

func query(queryData map[string]string) string {
	queryString := ""

	if len(queryData) > 0 {
		v := urlpkg.Values{}
		for key, value := range queryData {
			v.Set(key, value)
		}
		queryString += "?" + httpBuildQuery(v)
	}

	return queryString
}

func httpBuildQuery(queryData urlpkg.Values) string {
	return queryData.Encode()
}

type Breadcrumb struct {
	Name string
	URL  string
}

func breadcrumbsUI(breadcrumbs []Breadcrumb) hb.TagInterface {

	ol := hb.OL().Attr("class", "breadcrumb")

	for _, breadcrumb := range breadcrumbs {

		link := hb.Hyperlink().
			HTML(breadcrumb.Name).
			Href(breadcrumb.URL)

		li := hb.LI().
			Class("breadcrumb-item").
			Child(link)

		ol.AddChild(li)
	}

	nav := hb.Nav().
		Class("d-inline-block").
		Attr("aria-label", "breadcrumb").
		Child(ol)

	return nav
}
