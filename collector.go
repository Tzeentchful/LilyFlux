package LilyFlux

import "fmt"
import "time"
import "github.com/influxdb/influxdb-go"
import fluxConnect "github.com/Tzeentchful/LilyFlux/connect"

func StartCollector(connect *fluxConnect.FluxConnect, addr *string, port *uint16, user *string, pass *string, database *string, done chan bool) {

	// Create a new object for the Influx DB instance
	client, err := influxdb.NewClient(&influxdb.ClientConfig{
		Host:     fmt.Sprintf("%s:%d", *addr, *port),
		Username: *user,
		Password: *pass,
		Database: *database,
	})

	// Werid stuff happens with JSON unmarshalling when compression is enabled
	// We aren't sending much data anyway
	client.DisableCompression()

	if err != nil {
		fmt.Println("Error creating Influx client: ", err)
		return
	}

	// Start logging the player count to influx every 5 seconds
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-ticker.C:
				if !connect.Client.Connected() {
					break
				}
				series := &influxdb.Series{
					Name:    "network",
					Columns: []string{"players"},
					Points: [][]interface{}{
						[]interface{}{connect.Players},
					},
				}

				// Write the data to the Influx DB with time precision in seconds
				if err := client.WriteSeriesWithTimePrecision([]*influxdb.Series{series}, "s"); err != nil {
					fmt.Println("Error while writing data:", err)
				}
			case <-done:
				ticker.Stop()
				return
			}
		}
	}()

}
