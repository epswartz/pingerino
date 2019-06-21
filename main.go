package main

import (
	"bytes"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	ping "github.com/sparrc/go-ping"
	"github.com/wcharczuk/go-chart" //exposes "chart"
)

func main() {

	// Set up pinger

	pinger, err := ping.NewPinger("8.8.8.8")
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}
	pinger.Count = 10 // It does 10 packets

	seqs := []float64{}
	times := []float64{}
	xTicks := []chart.Tick{} // X axis ticks
	yTicks := []chart.Tick{}
	maxTime := time.Duration(0) // Tracks max time for drawing ticks on Y axis

	// On recieve, put data into graph
	pinger.OnRecv = func(pkt *ping.Packet) {
		seqs = append(seqs, float64(pkt.Seq))
		times = append(times, float64(pkt.Rtt/time.Millisecond))
		xTicks = append(xTicks, chart.Tick{
			Label: strconv.Itoa(pkt.Seq),
			Value: float64(pkt.Seq),
		})
		if pkt.Rtt > maxTime {
			maxTime = pkt.Rtt
		}
		fmt.Println("Recorded: ", pkt.Seq, pkt.Rtt)
	}

	pinger.Run() // blocks until it's done

	// Build Y axis ticks
	tickSize := math.Ceil(float64(maxTime/time.Millisecond) / 10)
	for i := float64(1); i <= 10; i++ {
		yTicks = append(yTicks, chart.Tick{
			Label: fmt.Sprintf("%.1f", tickSize*i),
			Value: tickSize * i,
		})
	}

	// Init graph object
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Style:     chart.StyleShow(),
			NameStyle: chart.StyleShow(),
			Name:      "icmp_seq",
			Ticks:     xTicks,
		},
		YAxis: chart.YAxis{
			Style:     chart.StyleShow(),
			NameStyle: chart.StyleShow(),
			Name:      "rtt (ms)",
			Ticks:     yTicks,
		},
		Series: []chart.Series{
			chart.ContinuousSeries{
				XValues: seqs,
				YValues: times,
			},
		},
	}

	buffer := bytes.NewBuffer([]byte{})
	err = graph.Render(chart.PNG, buffer)
	if err != nil {
		panic(err)
	}
	file, err := os.Create("Graph.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Write(buffer.Bytes())
}
