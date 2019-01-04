package http

import (
	"bytes"
	"chatbot/utils"
	"context"
	"github.com/davecgh/go-spew/spew"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	CONTENT_TYPE_JSON = "application/json"
	ENV_DNS_FILE_PATH = "USE_DNS_FILE"
)

var ApiHttpClient *http.Client
var DefaultDnsFileResolver *FileDnsResolver

func init() {

	var transport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	ApiHttpClient = &http.Client{
		Transport: transport,
	}

	var zone = make(map[string]Addressing, 0)
	if dnsFile := utils.EnvVar(ENV_DNS_FILE_PATH); dnsFile != "" {
		if _, err := os.Stat(dnsFile); !os.IsNotExist(err) {
			if bs, err := ioutil.ReadFile(dnsFile); err == nil {
				fileContent := string(bs)
				lines := strings.Split(fileContent, "\n")
				for _, line := range lines {
					lp := strings.Split(line, "=")
					if len(lp) == 2 && lp[1] != "" {
						domain := strings.TrimSpace(lp[0])
						ips := strings.TrimSpace(lp[1])
						zone[domain] = strings.Split(ips, ",")
					}
				}
			}
		} else {
			log.Printf("disable dns file, failed to open dns file %s , %s \n", dnsFile, err)
		}
	} else {
		log.Println("disable dns file")
	}

	if len(zone) > 0 {
		DefaultDnsFileResolver = &FileDnsResolver{
			ZoneFile: zone,
		}

		log.Println("dns file is enabled")
		spew.Dump(zone)
	}

}

func NewApiClient() *http.Client {

	if DefaultDnsFileResolver != nil {
		var netDialer *net.Dialer = &net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}

		ApiHttpClient.Transport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			ss := strings.Split(addr, ":")
			domain := ss[0]
			port := ss[1]
			if ip := DefaultDnsFileResolver.Lookup(domain); ip != "" {
				return netDialer.DialContext(ctx, network, ip+":"+port)
			} else {
				return netDialer.DialContext(ctx, network, addr)
			}
		}
	}
	return ApiHttpClient
}

type FileDnsResolver struct {
	ZoneFile map[string]Addressing
}

type Addressing []string

func (addr Addressing) RR() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	idx := r.Intn(len(addr))
	return addr[idx]
}

func (fdr FileDnsResolver) Lookup(host string) (ip string) {
	if v, ok := fdr.ZoneFile[host]; ok {
		return v.RR()
	} else {
		return ""
	}
}

func PostJson(url string, body []byte) (jsonBytes []byte, err error) {
	resp, err := NewApiClient().Post(url, CONTENT_TYPE_JSON, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	} else {
		respBytes, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return nil, err
		} else {
			return respBytes, nil
		}
	}
}

func GetJson(url string) (jsonBytes []byte, err error) {
	resp, err := NewApiClient().Get(url)
	if err != nil {
		return nil, err
	} else {
		respBytes, err := ioutil.ReadAll(resp.Body)
		defer resp.Body.Close()
		if err != nil {
			return nil, err
		} else {
			return respBytes, nil
		}
	}
}
