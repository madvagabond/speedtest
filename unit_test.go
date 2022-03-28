package speedtest 


import (
	"testing"
	"net/http"
	"fmt"
)


func TestGetServers(t *testing.T) {
	var cli http.Client

	servers, err := GetServerList(&cli) 

	if err != nil {
		t.Fatal("Returned an error")
	}

	if len(servers) < 1 {
		t.Fatal("The serverlist is empty")
	}

} 


func getServer() Server {
	var cli http.Client
	servers, _ := GetServerList(&cli)

	return  SelectClosest(servers) 

}


func TestPing(t *testing.T) {
	var cli http.Client 

	srv :=  getServer()
	dur, err := PingLatency(&cli, &srv)
	fmt.Println( dur.Milliseconds())

	if err != nil {t.Fatal("ping returned error"  ) }
}


func TestDownload(t *testing.T) {
	var cli http.Client
	srv := getServer() 

	speed, err := DownloadSpeed(&cli, &srv, 0)
	fmt.Println(speed)
	if err != nil {t.Fatal("Download Test Failed")}
}


func TestUpload(t *testing.T) {

	var cli http.Client 
	srv := getServer() 
	speed, err := UploadSpeed(&cli, &srv, 0)
	fmt.Println(speed)
	if err != nil {t.Fatal("Speed Test Failed")}
}
