# MQTT Temp Sensor

Simple temperature sensor for the ds18b20 that sends it's results to a MQTT broker. It was written for the raspberry pi but it will likely run anywhere.

## Requirements
This assumes a ds18b20 is wired properly to the raspberry pi and a working mqtt broker. Cert authentication with mqtt is currently not supported. I suggest only using this on an internal network for that reason.

| Pi | DS158B20 |
| --- | --- |
| vcc | 3.3v |
| gnd | gnd |
| sig | GPIO4 |

4.7K pullup resistor across vcc and sig. This is important to prevent inaccurate readings. If a longer line is needed (+1m) I suggest using a sheilded cable like eth5 or eth6 and reducing the pull up resistor to 2.2K.


## Basic usage
`./tempMon -mqtt.host=192.168.1.10:1883 -mqtt.username=username -mqtt.password=password`

```
Usage of tempMon:
  -bind_interface string
    	Interface the network is bound to (default "wlan0")
  -debug
    	Debugging using test_data.txt
  -mqtt.clientID string
    	MQTT client ID (default "room_temps")
  -mqtt.host string
    	MQTT Host (default "192.168.1.10:1883")
  -mqtt.password string
    	MQTT password (default "pass")
  -mqtt.username string
    	MQTT username (default "user")
```

If using multiple devices make sure `-mqtt.clientID` is unique for each one.

## Install

1. Download latest version from https://github.com/bah2830/Temp-Sensor-MQTT/releases or build yourself
2. Copy to any directory on the pi. (/home/pi/tempMon)
3. Update /etc/rc.local to start at boot. Add to end of file before `exit 0`
```
modprobe w1-gpio
modprobe w1-therm
/home/pi/tempMon -mqtt.host=192.168.1.10:1883 -mqtt.username=user -mqtt.password=pass -mqtt.clientID=XXX
```
4. Reboot

## MQTT Message

### Topic
`sensors/temp/[mac_address]/[sensor_name]`

### Message
```
{
  "mac_address" : "11:11:11:11:11:11",
  "ip_address" : "192.168.1.11",
  "ts" : 1514684628,
  "celcius" : 24.625,
  "fahrenheit" : 76.325
}
```

## Troubleshooting
1. Verify `/sys/bus/w1/devices/w1_bus_master1/w1_master_slaves` has has a device listed. Should be something like `28-XXXXXXXXXxXXX`
2. Verify `/sys/bus/w1/devices/28-XXXXXXXXXxXX/w1_slave` contains text data
```
79 01 4b 46 7f ff 0c 10 29 : crc=29 YES
79 01 4b 46 7f ff 0c 10 29 t=23562
```

If any of the above are not working google instructions for getting the ds18b20 working before using this.