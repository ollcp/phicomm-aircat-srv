package aircat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type writer interface {
	write(mac string, json string)
}

//example for line procotol:
//curl -i -XPOST 'http://localhost:8086/write?db=mydb' --data-binary 'cpu_load_short,host=server01,region=us-west value=0.64 1434055562000000000'
//we use as :
//aircat,mac=xxx temperature=1,humidity=2,value=3,hcho=4
type influxdb struct {
	db    string
	addr  string
	token string
}

func (s *influxdb) write(mac string, json string) {
	if s.addr == "" {
		println(mac, json)
		return
	}
	if line := formatLineProtocol(mac, json); line != "" {
		//we ignore error
		go func() {
			client := &http.Client{}
			url := fmt.Sprintf("http://%s/write?db=%s", s.addr, s.db)
			req, err := http.NewRequest("POST", url, strings.NewReader(line))
			if err != nil {
				return
			}
			req.Header.Set("Authorization", s.token)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			resp, err := client.Do(req)
			resp.Body.Close()
		}()
	}

}
func formatLineProtocol(mac string, js string) string {
	var air AirMeasure
	if err := json.Unmarshal([]byte(js), &air); err != nil {
		return ""
	}
	return fmt.Sprintf("aircat,mac=\"%s\" humidity=%s,temperature=%s,value=%s,hcho=%s", mac, air.Humidity, air.Temperature, air.Value, air.Hcho)
}

//AirMeasure reported from device
type AirMeasure struct {
	Humidity    string
	Temperature string
	Value       string
	Hcho        string
}
