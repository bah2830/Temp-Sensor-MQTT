package ds18b20

import (
	"errors"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

var path = "/sys/bus/w1/devices/"

type Sensor struct {
	Name        string
	LastReading SensorReading
}

type SensorReading struct {
	T          time.Time
	Celsius    float64
	Fahrenheit float64
}

func SetPath(p string) {
	path = p
}

// Sensors get all connected sensor IDs as array
func Sensors() (sensors []*Sensor, err error) {
	data, err := ioutil.ReadFile(path + "w1_bus_master1/w1_master_slaves")
	if err != nil {
		return sensors, err
	}

	sensorNames := strings.Split(string(data), "\n")
	if len(sensorNames) == 0 {
		return sensors, errors.New("No sensors found")
	}

	for _, name := range sensorNames {
		if name != "" {
			sensor := &Sensor{
				Name: name,
			}

			sensors = append(sensors, sensor)
		}
	}

	return sensors, err
}

// GetTemperature get the temperature of a given sensor
func (s *Sensor) GetTemperature() (SensorReading, error) {
	newReading := SensorReading{
		T: time.Now(),
	}

	data, err := ioutil.ReadFile(path + s.Name + "/w1_slave")
	if err != nil {
		return newReading, err
	}

	if strings.Contains(string(data), "YES") {
		arr := strings.SplitN(string(data), " ", 3)

		celcius := 0.0

		switch arr[1][0] {
		case 'f': //-0.5 ~ -55°C
			x, err := strconv.ParseInt(arr[1]+arr[0], 16, 32)
			if err != nil {
				return newReading, err
			}
			celcius = float64(^x+1) * 0.0625

		case '0': //0~125°C
			x, err := strconv.ParseInt(arr[1]+arr[0], 16, 32)
			if err != nil {
				return newReading, err
			}
			celcius = float64(x) * 0.0625
		}

		newReading.Celsius = celcius
		newReading.Fahrenheit = celcius*1.8 + 32

		s.LastReading = newReading
		return newReading, nil
	}

	return newReading, errors.New("Cannot read current temperature data for" + s.Name)
}
