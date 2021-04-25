package main

import (
	"net/url"

	"github.com/yscsky/yu"
)

func main() {
	hc := yu.NewHttpClient()
	var resp yu.Resp
	if err := hc.GetJSON("http://127.0.0.1:8080/api/getjson", &resp); err != nil {
		yu.LogErr(err, "GetJSON")
		return
	}
	yu.Logf("code: %d, data: %s", resp.Code, resp.Data)
	if err := hc.GetJSONAuth("http://127.0.0.1:8080/api/getjsonauth", "admin", "123456", &resp); err != nil {
		yu.LogErr(err, "GetJSON")
		return
	}
	yu.Logf("code: %d, data: %s", resp.Code, resp.Data)
	if err := hc.PostJSON("http://127.0.0.1:8080/api/postjson", &yu.Resp{Data: "hello"}, &resp); err != nil {
		yu.LogErr(err, "GetJSON")
		return
	}
	yu.Logf("code: %d, data: %s", resp.Code, resp.Data)
	if err := hc.PostJSONAuth("http://127.0.0.1:8080/api/postjsonauth", "admin", "123456", &yu.Resp{Data: "hello"}, &resp); err != nil {
		yu.LogErr(err, "GetJSON")
		return
	}
	yu.Logf("code: %d, data: %s", resp.Code, resp.Data)
	if err := hc.PostFormJSON("http://127.0.0.1:8080/api/postform", url.Values{"data": []string{"hello"}}, &resp); err != nil {
		yu.LogErr(err, "GetJSON")
		return
	}
	yu.Logf("code: %d, data: %s", resp.Code, resp.Data)
	if err := hc.PostFormJSONAuth("http://127.0.0.1:8080/api/postformauth", "admin", "123456", url.Values{"data": []string{"hello"}}, &resp); err != nil {
		yu.LogErr(err, "GetJSON")
		return
	}
	yu.Logf("code: %d, data: %s", resp.Code, resp.Data)
}
