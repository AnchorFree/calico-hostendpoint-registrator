package main

import (
	"context"
	"github.com/anchorfree/golang/pkg/jsonlog"
	"github.com/kelseyhightower/envconfig"
	"github.com/projectcalico/libcalico-go/lib/apis/v3"
	client "github.com/projectcalico/libcalico-go/lib/clientv3"
	"github.com/projectcalico/libcalico-go/lib/options"
	"os"
)

type Config struct {
	SetLabel      string `default:"" split_words:"true"`
	EndpointLabel string `default:"" split_words:"true"`
	Hostname      string
}

type App struct {
	config Config
	log    jsonlog.Logger
}

// NewApp initializes the logger and parses the configuration.
func NewApp() *App {

	log := &jsonlog.StdLogger{}
	log.Init("hep", false, false, nil)

	app := &App{config: Config{}}
	err := envconfig.Process("hep", &app.config)
	if err != nil {
		log.Fatal("can't initialize application", err)
	}
	app.log = log

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal("can't get node's hostname", err)
	}
	app.config.Hostname = hostname
	return app

}

// CreateHostEndpoint creates a new HostEndpoint resource.
func CreateHostEndpoint(name, label, iface string) error {

	cl, err := client.NewFromEnv()
	if err != nil {
		return err
	}

	HEPInterface := cl.HostEndpoints()

	HEP := v3.NewHostEndpoint()
	HEP.ObjectMeta.Name = name
	HEP.ObjectMeta.Labels = map[string]string{name: "true", label: "true"}
	Spec := v3.HostEndpointSpec{Node: name, InterfaceName: iface}
	HEP.Spec = Spec

	ctx := context.TODO()
	_, err = HEPInterface.Create(ctx, HEP, options.SetOptions{})
	return err

}

// CreateGNS creates a new GNS resource.
func CreateGNS(name, label, nodeIP string) error {

	cl, err := client.NewFromEnv()
	if err != nil {
		return err
	}

	GNSInterface := cl.GlobalNetworkSets()

	GNS := v3.NewGlobalNetworkSet()
	GNS.ObjectMeta.Name = name
	GNS.ObjectMeta.Labels = map[string]string{name: "true", label: "true"}
	HostList := v3.GlobalNetworkSetSpec{Nets: []string{nodeIP + "/32"}}
	GNS.Spec = HostList

	ctx := context.TODO()
	_, err = GNSInterface.Create(ctx, GNS, options.SetOptions{})
	return err

}

func main() {

	app := NewApp()

	iface, nodeIPWithCIDR := GetNetworkConfig()

	err := CreateHostEndpoint(app.config.Hostname, app.config.EndpointLabel, iface)
	if err != nil {
		app.log.Fatal("failed to create hostendpoint", err)
	}

	err = CreateGNS(app.config.Hostname, app.config.SetLabel, nodeIPWithCIDR[:len(nodeIPWithCIDR)-3])
	if err != nil {
		app.log.Fatal("failed to create GNS", err)
	}

	quit := make(chan bool, 1)
	_ = <-quit

}
