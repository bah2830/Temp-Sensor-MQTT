# MQTT Temp Sensor

Simple temperature sensor for the ds18b20 that sends it's results to a MQTT broker. It was written for the raspberry pi but it will likely run anywhere.

## Basic usage
`./tempMon -mqtt.host=192.168.1.10:1883 -mqtt.username=username -mqtt.password=password`

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
