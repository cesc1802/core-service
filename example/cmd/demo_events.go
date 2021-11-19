package cmd

import (
	"fmt"
	"github.com/cesc1802/core-service/events"
	"github.com/spf13/cobra"
	"sync"
)

var demoEvtCmd = &cobra.Command{
	Use: "event",
	Run: func(cmd *cobra.Command, args []string) {
		stream, _ := events.NewStream()
		var wg sync.WaitGroup

		wg.Add(1)
		go func() {
			evt, _ := stream.Consume("test")

			data := <-evt
			fmt.Println(data.Topic)
			fmt.Println(data.ID)
		}()

		stream.Publish("test", map[string]interface{}{
			"username": "thuocnv",
			"password": "123456",
		})

		wg.Done()
	},
}
