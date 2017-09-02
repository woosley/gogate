package main

import (
	"flag"
	"fmt"
	"github.com/woosley/gogate/gate"
	"github.com/woosley/gogate/gate/types"
	"net/url"
	"os"
)

var options types.Opt = types.Opt{
	Listen:    1234,
	Help:      false,
	Is_master: false,
	Expire:    60,
	Key:       "ip",
}

func init() {
	flag.StringVar(&options.Master_addr, "master", "", "master url")
	flag.StringVar(&options.Master_addr, "m", "", "master url")
	flag.IntVar(&options.Listen, "port", options.Listen, "listen port")
	flag.IntVar(&options.Listen, "p", options.Listen, "listen port")
	flag.IntVar(&options.Listen, "expire", options.Listen, "expire time in seconds")
	flag.BoolVar(&options.Is_master, "is-master", options.Is_master, "is master or not")
	flag.BoolVar(&options.Help, "help", options.Help, "show help")
	flag.BoolVar(&options.Help, "h", options.Help, "show help")
	flag.StringVar(&options.Key, "k", options.Key, "the uniq key to the nodes")
	flag.StringVar(&options.Key, "key", options.Key, "the uniq key to the nodes")
	flag.Usage = print_help
}

func print_help() {
	help := `gotypes -- keep your machine state
	-h|--help show help 
	-p|--port listen port
	-m|--master master url
	-k|--key the uniq key to the node
	--is-master start node as master
	--expire expire time in seconds on master
	`
	fmt.Println(help)
}

func main() {
	flag.Parse()
	if options.Help {
		print_help()
		os.Exit(0)
	}

	if (!options.Is_master) && options.Master_addr == "" {
		fmt.Println("not master and no master address")
		print_help()
		os.Exit(1)
	}

	bad_url := false
	if uri, err := url.ParseRequestURI(options.Master_addr); err != nil {
		bad_url = true
	} else {
		if uri.Scheme != "http" && uri.Scheme != "https" {
			bad_url = true
		}
	}

	if !options.Is_master && bad_url {
		fmt.Println("bad master url address")
		print_help()
		os.Exit(1)
	}

	gate.App(options)
}
