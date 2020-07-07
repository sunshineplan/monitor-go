package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type bank interface {
	formatter() string
	update()
}

func markup(old string, new float64) string {
	if old != "" {
		oldFloat, err := strconv.ParseFloat(strings.ReplaceAll(strings.ReplaceAll(old, "↑", ""), "↓", ""), 64)
		if err != nil {
			log.Printf("Parse Float Error: %v", err)
			return " "
		}
		if oldFloat > new {
			return "↓"
		} else if oldFloat == new {
			return " "
		}
	}
	return "↑"
}

type nbcb struct {
	id     string
	name   string
	income string
	limit  string
	mark   string
}

func (n nbcb) formatter() string {
	return fmt.Sprintf("%s(%s) %13s%s", n.name, n.income, n.limit, n.mark)
}

func (n *nbcb) update() {
	url := "https://i.nbcb.com.cn/zhongtai/finance/prds/p-onsale-channel?turnPageBeginPos=1&turnPageShowNum=6&prdOrigin=0&prdClassify=0"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Panicf("New Request Error: %v", err)
	}
	req.Header.Set("X-GW-APP-ID", "1101")
	req.Header.Set("X-GW-BACK-HTTP", "GET")
	client := &http.Client{Transport: &http.Transport{Proxy: nil}}
	resp, err := client.Do(req)
	if err != nil {
		log.Panicf("Do Request Error: %v", err)
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicf("ReadAll body Error: %v", err)
	}
	var jsonData interface{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		log.Panicf("Unmarshal json Error: %v", err)
	}
	list := jsonData.(map[string]interface{})["body"].(map[string]interface{})["list"]
	for _, i := range list.([]interface{}) {
		if i.(map[string]interface{})["prdCode"] == n.id {
			n.name = i.(map[string]interface{})["prdName"].(string)
			n.income = i.(map[string]interface{})["expectedRateShow"].(string)
			n.mark = markup(n.limit, i.(map[string]interface{})["perUseLimit"].(float64))
			n.limit = fmt.Sprintf("%.2f", (i.(map[string]interface{})["perUseLimit"].(float64)))
			break
		}
	}
}
