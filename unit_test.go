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
