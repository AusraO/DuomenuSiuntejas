package main

import (
	// "encoding/json"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

type wetnessData struct {
	Hour        int64     `bson:"hour"`
	Wetness     int64     `bson:"wetness"`
	DateOfEntry time.Time `bson:"dateOfEntry"`
}

type SensorData struct {
	SensorName string        `bson:"sensorName"`
	Data       []wetnessData `bson:"data"`
}

var wg sync.WaitGroup

func main() {
	fmt.Println("Starting to send data", time.Now())
	serverAddress := "localhost:5225"
	connection, err := net.Dial("tcp", serverAddress)
	if err != nil {
		fmt.Println("Error connecting", err)
		return
	}
	defer connection.Close()

	wg.Add(1)

	go newDataGenerator(connection)

	wg.Wait()
	fmt.Println("Finished sending data", time.Now())
}
func newDataGenerator(connection net.Conn) {
	defer wg.Done()
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	var rawData []byte // To store the accumulated raw data
	var num int64 = 0

	var start_time uint64 = uint64(time.Time.UnixMicro(time.Now()))
	var data_in_bytes int = 0
	var single_packet_size = 0
	for i := 1; i <= 1000000; i++ {
		data := wetnessData{
			Hour:        int64((i-1)%24) + 1,
			Wetness:     int64(r.Intn(101)),
			DateOfEntry: time.Now(),
		}

		// Convert your wetnessData struct to a byte slice manually
		entry := fmt.Sprintf("Hour:%d Wetness:%d DateOfEntry:%s\n", data.Hour, data.Wetness, data.DateOfEntry)
		rawData = append(rawData, []byte(entry)...)
		data_in_bytes += len(rawData)

		// only for the first time
		if single_packet_size == 0 {
			single_packet_size = len(rawData)
		}

		// if i%1000 == 0 {
		_, err := connection.Write(rawData)
		if err != nil {
			fmt.Println("Failed to send data", err)
			return
		}
		num++
		fmt.Println("...sent...", num)
		rawData = nil // Reset rawData for the next batch
		// }
	}
	var end_time uint64 = uint64(time.Time.UnixMicro(time.Now()))

	var total_time = end_time - start_time

	fmt.Printf(" %f: MB sent\r\n", float64(float64(data_in_bytes)/1000000.0))
	fmt.Printf("Start time in %d us\r\n:", start_time)
	fmt.Printf("End time in %d us\r\n:", end_time)
	fmt.Printf("Single packet size:%d Bytes\r\n", single_packet_size)
	fmt.Println("Time spent on TCP packets sending in ms(miliseconds):", float64(float64(total_time)/1000.0000))
	fmt.Printf("Transfer speed:%f MB/s \r\n", float64(float64(data_in_bytes)/1000000.0)/float64(float64(total_time)/1000000.0000))
	// Send any remaining data that didn't reach 1000 iterations
	if len(rawData) > 0 {
		_, err := connection.Write(rawData)
		if err != nil {
			fmt.Println("Failed to send data", err)
			return
		}
	}
}
