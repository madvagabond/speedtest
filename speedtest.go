
package speedtest

import (
	"net/http"
	"encoding/json"
	"time"
	"fmt"
	"strings"


	"net/url"

	"io"
	"io/ioutil"
)

const serverList = "https://www.speedtest.net/api/js/servers?engine=js&limit=10"


type Server struct {
	Url string
	Country string 
	Host string 
	Distance int 
}




type Servers []Server 


func GetServerList(c *http.Client) (Servers, error) {
	var servers Servers
	resp, err := c.Get(serverList) 
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&servers)

	return servers, err
}


func isDistSorted(servers Servers) bool {
	first := servers[0] 
	j := len(servers) - 1
	last := servers[j]

	return last.Distance >= first.Distance
}



func SelectClosest(servers Servers) Server {
	if isDistSorted(servers) {
		return servers[0]
	}

	closest := servers[0]
	for _, x := range servers[1:] {
		if x.Distance < closest.Distance {

			closest = x
		}
	}


	return closest 
}






var dlSizes = [...]int{350, 500, 750, 1000, 1500, 2000, 2500, 3000, 3500, 4000}
var ulSizes = [...]int{100, 300, 500, 800, 1000, 1500, 2500, 3000, 3500, 4000} //kB




// size is an index for an array of sizes from 350x350 to 4Kx4K with range 0 to 9  
func Download(client *http.Client, server *Server, size int) error {

	prefix := strings.Split(server.Url, "upload.php")[0]
	dlSize := fmt.Sprint("%dxd%", dlSizes[size], dlSizes[size])
	dlUri :=  prefix + "/random" + dlSize 
	resp, e := client.Get(dlUri) 

	defer resp.Body.Close()
	if e != nil {
		return e 
	}


	_, err := io.Copy(ioutil.Discard, resp.Body)
	return err 
	
}


func calcSpeed(size int, dur time.Duration) float64 {
	reqMB := size * size * 2 / 1000 / 1000
	res := float64(reqMB) * 8 / float64(dur.Seconds())
	return res
}




func Upload(client *http.Client, server *Server, size int) error {

	size1 := ulSizes[size] 
	v := url.Values{}
	v.Add("content", strings.Repeat("0123456789", size1*100-51))


	body := strings.NewReader(v.Encode())
	
	resp, err := client.Post(server.Url, "application/x-www-form-urlencoded", body) 
	defer resp.Body.Close() 

	if err != nil {
		return err
	}
	_, err = io.Copy(ioutil.Discard, resp.Body)
	return err
}






func TestDownload(client *http.Client, server *Server, size int) (float64, error) {
	start := time.Now()
	err := Download(client, server, size) 
	end := time.Now()
	dur := end.Sub(start)

	dlSize := dlSizes[size]
	speed := calcSpeed(dlSize, dur)
	return speed, err

}


func TestPing(client *http.Client, server *Server) (time.Duration, error) {
	url := server.Url 
	pingURL := strings.Split(url, "/upload.php")[0] + "/latency.txt"

	start := time.Now()
	_, err := client.Get(pingURL) 
	end := time.Now() 

	latency :=  end.Sub(start) 
	return latency, err 

}


