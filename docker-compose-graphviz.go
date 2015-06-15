package main

import (
	"fmt"
	"github.com/alexcesaro/log"
	"github.com/alexcesaro/log/stdlog"
	"github.com/awalterschulze/gographviz"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

var logger log.Logger

func abort(msg string) {
	logger.Critical(msg)
	os.Exit(1)
}

type service struct {
	Links       []string
	VolumesFrom []string "volumes_from"
	Ports       []string
}

func main() {
	var (
		bytes   []byte
		data    map[string]service
		err     error
		graph   *gographviz.Graph
		project string
	)
	logger = stdlog.GetFromFlags()
	project = ""

	// Load docker-compose.yml
	bytes, err = ioutil.ReadFile("docker-compose.yml")
	if err != nil {
		abort(err.Error())
	}

	// Parse it as YML
	data = make(map[string]service, 5)
	yaml.Unmarshal(bytes, &data)
	if err != nil {
		abort(err.Error())
	}

	// Create directed graph
	graph = gographviz.NewGraph()
	graph.SetName(project)
	graph.SetDir(true)

	// Add legend
	graph.AddSubGraph(project, "cluster_legend", map[string]string{"label": "Legend"})
	graph.AddNode("cluster_legend", "legend_service", map[string]string{"label": "service"})
	graph.AddNode("cluster_legend", "legend_service_with_ports", map[string]string{"label": "\"service with ports\\n80:80 443:443\""})
	graph.AddEdge("legend_service", "legend_service_with_ports", true, map[string]string{"label": "links"})
	graph.AddEdge("legend_service_with_ports", "legend_service", true, map[string]string{"label": "volumes_from", "style": "dashed"})

	// Round one: populate nodes
	for name, service := range data {
		var attrs = map[string]string{"label": name}
		if service.Ports != nil {
			attrs["label"] += "\\n" + strings.Join(service.Ports, " ")
			attrs["shape"] = "box"
		}
		attrs["label"] = fmt.Sprintf("\"%s\"", attrs["label"])
		graph.AddNode(project, name, attrs)
	}

	// Round two: populate connections
	for name, service := range data {
		// links
		if service.Links != nil {
			for _, linkTo := range service.Links {
				if strings.Contains(linkTo, ":") {
					linkTo = strings.Split(linkTo, ":")[0]
				}
				graph.AddEdge(name, linkTo, true, nil)
			}
		}
		// volumes_from
		if service.VolumesFrom != nil {
			for _, linkTo := range service.VolumesFrom {
				graph.AddEdge(name, linkTo, true, map[string]string{"style": "dotted"})
			}
		}
	}

	fmt.Print(graph)
}
