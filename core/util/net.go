package util

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
)

func GetIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", nil
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", nil
	// name, err := os.Hostname()
	// if err != nil {
	// 	return "", err
	// }

	// addrs, err := net.LookupHost(name)
	// if err != nil {
	// 	return "", err
	// }
	// return addrs[0], nil
}

func CreateListener() (l net.Listener, close func()) {
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}
	return l, func() {
		_ = l.Close()
	}
}

func MakeHttpRequest[T any](method string, url string, data interface{}, returnValue T) error {
	var body io.Reader = nil
	if data != nil {
		bytesData, err := json.Marshal(data)
		if err != nil {
			return err
		}
		body = bytes.NewReader(bytesData)
	}
	req, _ := http.NewRequest(method, url, body)
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return nil
	}
	responseBody, _ := io.ReadAll(resp.Body)
	if err = json.Unmarshal(responseBody, &returnValue); err != nil {
		return err
	}
	return nil
}
