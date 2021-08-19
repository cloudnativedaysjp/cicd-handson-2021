package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"syscall/js"

	"local.packages/common"
)

var document = js.Global().Get("document")
var app = js.Global().Get("document").Call("getElementById", "app")

type Project common.Project
type Projects []Project

func main() {
	c := make(chan bool)
	go load()
	<-c
}

func load() {
	projects, err := getProjects()
	if err != nil {
		landscapeDiv := document.Call("createElement", "div")
		landscapeDiv.Get("classList").Call("add", "landscape")
		app.Call("appendChild", landscapeDiv)

		errorDiv := document.Call("createElement", "div")
		errorDiv.Set("textContent", "error")
		landscapeDiv.Call("appendChild", errorDiv)
		return
	}

	loadByRelation("CNCF Graduated Projects", filterByRelation(projects, "graduated"))
	loadByRelation("CNCF Incubating Projects", filterByRelation(projects, "incubating"))
	loadByRelation("CNCF Sandbox Projects", filterByRelation(projects, "sandbox"))
	loadByRelation("CNCF Member Products/Projects", filterByRelation(projects, "member"))
	loadByRelation("Non-CNCF Member Products/Projects", filterByRelation(projects, ""))
}

func loadByRelation(rel string, projects Projects) {
	if len(projects) == 0 {
		return
	}

	relationDiv := document.Call("createElement", "div")
	relationDiv.Get("classList").Call("add", "relation")
	relationDiv.Set("textContent", rel+" ("+strconv.Itoa(len(projects))+")")
	app.Call("appendChild", relationDiv)

	landscapeDiv := document.Call("createElement", "div")
	landscapeDiv.Get("classList").Call("add", "landscape")
	app.Call("appendChild", landscapeDiv)

	for _, proj := range projects {
		landscapeDiv.Call("appendChild", createItem(proj))
	}
}

func getProjects() (Projects, error) {
	resp, err := http.Get("/landscape")
	if err != nil {
		return Projects{}, err
	}
	defer resp.Body.Close()

	projects := Projects{}
	err = json.NewDecoder(resp.Body).Decode(&projects)
	return projects, err
}

func filterByRelation(projects Projects, relation string) Projects {
	var filtered Projects
	for _, proj := range projects {
		if proj.Project == relation {
			filtered = append(filtered, proj)
		}
	}
	return filtered
}

func createItem(proj Project) js.Value {
	item := document.Call("createElement", "div")
	item.Get("classList").Call("add", "item")

	itemLink := createHrefContent(proj)
	item.Call("appendChild", itemLink)

	if proj.Description == "" {
		itemLink.Call("appendChild", createTextContent(getTitle(proj), "item-title-only"))
	} else {
		itemLink.Call("appendChild", createTextContent(getTitle(proj), "item-title"))
		itemLink.Call("appendChild", createTextContent(proj.Description, "item-description"))
	}

	return item
}

func createTextContent(txt, cls string) js.Value {
	item := document.Call("createElement", "div")
	item.Set("textContent", fmt.Sprintf("%s", txt))
	item.Get("classList").Call("add", cls)
	return item
}

func createHrefContent(proj Project) js.Value {
	url := proj.RepoUrl
	if url == "" {
		url = proj.HomepageUrl
	}
	item := document.Call("createElement", "a")
	item.Set("href", url)
	return item
}

func getTitle(proj Project) string {
	if proj.StarCount == 0 {
		return proj.Name
	}
	return fmt.Sprintf("%s (☆%d)", proj.Name, proj.StarCount)
}
