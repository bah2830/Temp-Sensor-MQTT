package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/bah2830/Temp-Sensor-MQTT/ds18b20"
	"github.com/bah2830/Temp-Sensor-MQTT/mqtt"
	"github.com/prometheus/common/log"
)

var (
	debug         = flag.Bool("debug", false, "Debugging using test_data.txt")
	mqttHost      = flag.String("mqtt.host", "192.168.1.10:1883", "MQTT Host")
	mqttUsername  = flag.String("mqtt.username", "user", "MQTT username")
	mqttPassword  = flag.String("mqtt.password", "pass", "MQTT password")
	mqttClientID  = flag.String("mqtt.clientID", "room_temps", "MQTT client ID")
	bindInterface = flag.String("bind_interface", "wlan0", "Interface the network is bound to")
	ipAddress     = ""
	macAddress    = ""
	sensors       = make([]*ds18b20.Sensor, 0)
)

type messagePayload struct {
	MacAddress string  `json:"mac_address"`
	IPAddress  string  `json:"ip_address"`
	TimeStamp  int64   `json:"ts"`
	Celcius    float64 `json:"celcius"`
	Fahrenheit float64 `json:"fahrenheit"`
}

func main() {
	flag.Parse()
	setIPAndMacAddress()
	getSensors()

	mqtt := mqttConnect()
	defer mqtt.Disconnect()

	go startMonitoring(mqtt)

	k := make(chan os.Signal, 2)
	<-k
}

func mqttConnect() *mqtt.Client {
	mqtt, err := mqtt.Connect(mqtt.ConnectionOptions{
		Broker:   "tcp://" + *mqttHost,
		Username: *mqttUsername,
		Password: *mqttPassword,
		ClientID: *mqttClientID,
	})
	if err != nil {
		panic(err)
	}

	return mqtt
}

func startMonitoring(mqtt *mqtt.Client) {
	t := time.NewTicker(10 * time.Second)
	var loopCheck = func() {
		for _, sensor := range sensors {
			if _, err := sensor.GetTemperature(); err != nil {
				log.Errorf("Error getting reading from %s: %s", sensor.Name, err)
				continue
			}

			message, err := getMessagePayload(sensor)
			if err != nil {
				fmt.Println(err)
				continue
			}

			err = mqtt.SendMessage("sensors/temp/"+macAddress+"/"+sensor.Name, message)
			if err != nil {
				fmt.Println(err)
				continue
			}

		}
	}

	loopCheck()
	for range t.C {
		loopCheck()
	}
}

func getMessagePayload(sensor *ds18b20.Sensor) (string, error) {
	message := messagePayload{
		IPAddress:  ipAddress,
		MacAddress: macAddress,
		TimeStamp:  sensor.LastReading.T.Unix(),
		Celcius:    sensor.LastReading.Celsius,
		Fahrenheit: sensor.LastReading.Fahrenheit,
	}

	messageString, err := json.Marshal(message)
	if err != nil {
		return "", err
	}

	return string(messageString), nil
}

func getSensorsDebug() []*ds18b20.Sensor {
	ds18b20.SetPath("test_data/")
	sensor := &ds18b20.Sensor{
		Name: "debug_sensor",
	}

	sensors = append(sensors, sensor)

	return sensors
}

func getSensors() []*ds18b20.Sensor {
	if *debug {
		return getSensorsDebug()
	}

	var err error
	sensors, err = ds18b20.Sensors()
	if err != nil {
		panic(err)
	}

	if len(sensors) == 0 {
		panic("no sensors could be found")
	}
	return sensors
}

func setIPAndMacAddress() {
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			panic(err)
		}

		if i.Name == *bindInterface {
			for _, addr := range addrs {
				var ip net.IP
				switch v := addr.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				}

				if strings.Contains(ip.String(), "192.") {
					ipAddress = ip.String()
					macAddress = i.HardwareAddr.String()
				}
			}
		}
	}
}
