package main

import "os"
import "fmt"
import "strings"
import "runtime"
import "bufio"
import "github.com/Tzeentchful/LilyFlux/config"
import "github.com/Tzeentchful/LilyFlux"
import fluxConnect "github.com/Tzeentchful/LilyFlux/connect"

var VERSION string = "1.0"
var CONFIG_PATH string = "lilyflux.yml"

func main() {
	// Load the config file
	cfg, err := config.LoadConfig(CONFIG_PATH)
	if err != nil {
		cfg = config.DefaultConfig()
		err = config.SaveConfig(CONFIG_PATH, cfg)
		if err != nil {
			fmt.Println("Error while saving config", err)
			return
		}
	}

	// Start the lilypad connect client
	connectDone := make(chan bool)
	connect := fluxConnect.NewFluxConnect(&cfg.Connect.Address, &cfg.Connect.Credentials.Username, &cfg.Connect.Credentials.Password, connectDone)

	// Start the collector task
	collectorDone := make(chan bool)
	LilyFlux.StartCollector(connect, &cfg.Influx.Address, &cfg.Influx.Port, &cfg.Influx.Username, &cfg.Influx.Password, &cfg.Influx.Database, collectorDone)

	closeAll := func() {
		close(connectDone)
		close(collectorDone)
		os.Stdin.Close()
	}

	// Initialize vars for reading input
	serverErr := make(chan error, 1)
	stdinString := make(chan string, 1)
	stdinErr := make(chan error, 1)

	// Start getting input from user
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			str, err := reader.ReadString('\n')
			if err != nil {
				stdinErr <- err
			}
			stdinString <- str
		}
	}()

	fmt.Println("LilyFlux started, vesrion:", VERSION)

	// Interpret user's input
	for {
		select {
		case str := <-stdinString:
			str = strings.TrimSpace(str)
			if str == "reload" {
				fmt.Println("Reloading config...")
				newCfg, err := config.LoadConfig(CONFIG_PATH)
				if err != nil {
					fmt.Println("Error during reloading config", err)
					continue
				}
				*cfg = *newCfg
			} else if str == "debug" {
				fmt.Println("runtime.NumCPU:", runtime.NumCPU())
				fmt.Println("runtime.NumGoroutine:", runtime.NumGoroutine())
				memStats := &runtime.MemStats{}
				runtime.ReadMemStats(memStats)
				fmt.Println("runtime.MemStats.Alloc:", memStats.Alloc, "bytes")
				fmt.Println("runtime.MemStats.TotalAlloc:", memStats.TotalAlloc, "bytes")
			} else if str == "exit" || str == "stop" || str == "halt" {
				fmt.Println("Stopping...")
				closeAll()
				return
			} else if str == "help" {
				fmt.Println("LilyPad Proxy - Help")
				fmt.Println("reload - Reloads the proxy.yml")
				fmt.Println("debug  - Prints out CPU, Memory, and Routine stats")
				fmt.Println("stop   - Stops the process. (Aliases: 'exit', 'halt')")
			} else {
				fmt.Println("Command not found. Use \"help\" to view available commands.")
			}
		case err := <-stdinErr:
			fmt.Println("Error during stdin", err)
			closeAll()
			return
		case err := <-serverErr:
			fmt.Println("Error during listen", err)
			closeAll()
			return
		}
	}

}
