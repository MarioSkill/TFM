package eidas

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/url"
	//"os"
	"strconv"
	"strings"
	"time"
)

type Block struct {
	Try     func()
	Catch   func(Exception)
	Finally func()
}

type Exception interface{}

func Throw(up Exception) {
	panic(up)
}

func (tcf Block) Do() {
	if tcf.Finally != nil {

		defer tcf.Finally()
	}
	if tcf.Catch != nil {
		defer func() {
			if r := recover(); r != nil {
				tcf.Catch(r)
			}
		}()
	}
	tcf.Try()

}

var eidasServer string
var eidasPort string
var resultServer string
var timeServer string

// = "http://163.117.148.105:8080"
//var resultsServer = "http://163.117.148.105:8081/eidasResult"
var Successfullogin = false
var start time.Time

//example  ./eidas 99 http://163.117.148.105 6060 http://163.117.148.105:8081/eidasResult timestampClientRegisted, serverTime, timetaskDeploy

func StartTest(clientID string, eidasS string, eidasP string, resultS string, timestamp string, serverTime string, timetaskDeploy string) {
	//func Start() {

	var result bool
	var time float64

	//var clientID = os.Args[1]
	eidasServer = eidasS   //os.Args[2] //"http://163.117.148.105:8080"
	eidasPort = eidasP     //os.Args[3]
	resultServer = resultS //os.Args[4] //"http://163.117.148.105:8081/eidasResult"
	timeServer = serverTime

	//fmt.Println(clientID, eidasServer, eidasPort, resultServer, timestamp)

	data := url.Values{}
	Block{
		Try: func() {
			result, time = runTest(clientID)

		},
		Catch: func(e Exception) {
			data.Set("error", e.(error).Error())
			data.Set("ID", clientID)
			data.Set("login", "false")
			data.Set("ClienteRegistrado", timestamp)
			data.Set("InicioEjecucion", start.Format("15:04:05"))
			data.Set("TareaPlanificada", timetaskDeploy)
			data.Set("LoginTime", "-1")
			//fmt.Printf("Error %v\n", e)

			sendResult(resultServer, data)
		}, Finally: func() {
			data.Set("error", "null")
			data.Set("ID", clientID)
			data.Set("login", strconv.FormatBool(result))
			data.Set("ClienteRegistrado", timestamp)
			data.Set("InicioEjecucion", start.Format("15:04:05"))
			data.Set("TareaPlanificada", timetaskDeploy)
			f := fmt.Sprint(time)
			data.Set("LoginTime", f)
			//fmt.Println("Login Client ", clientID, result)
			sendResult(resultServer, data)
		},
	}.Do()
	//return false
	//return true
	//os.Exit(0)
}

/*
func TiempoDeEspera(t time.Time, cid string, timeStart string, clienteRegistrado string) {
	data := url.Values{}
	data.Set("ID", cid)
	data.Set("InicioEjecucion", t.Format("15:04:05"))
	data.Set("ClienteRegistrado", clienteRegistrado)
	data.Set("TareaPlanificada", timeStart)
	sendResult(timeServer, data)
}*/
func runTest(ClientID string) (bool, float64) {
	//fmt.Println("eIDAS Exec Time Start", time.Since(start))

	start = time.Now()
	//go TiempoDeEspera(start, ClientID, timeStart, clienteRegistrado)
	//1 Peticion populateIndexPage ..
	doc := Request("/SP/populateIndexPage", url.Values{}, true) //populateIndexPage ..

	// ***************** 2 Peticion /SP/populateIndexPage *****************
	//en esta peticion debemos añadir todos los atributos estadar definidos por el CEF ..
	var next = GetDataById(doc, "formTab2", "action")
	data := url.Values{}

	var atributes Atributes
	json.Unmarshal(GetDataSetByName("Representative Attributes", eidasServer+":"+eidasPort), &atributes)

	for _, item := range atributes.GetAttributes.Attributes {
		data.Set(item.Name, item.Value)
	}
	doc = Request(next, data, true)
	// ***************** 2 Peticion /SP/populateIndexPage *****************

	// ***************** 3 Peticion SpecificConnector/ServiceProvider *****************
	next = GetDataById(doc, "countrySelector", "action")

	next = eidasServer + ":" + eidasPort + "/SpecificConnector/ServiceProvider"
	data = url.Values{}
	//fmt.Println(next)
	data.Set("SMSSPRequest", GetDataById(doc, "SMSSPRequest", "value"))
	doc = Request(next, data, false) //SpecificConnector/ServiceProvider
	// ***************** 3 Peticion SpecificConnector/ServiceProvider *****************

	// ***************** 4 Peticion EidasNode/SpecificConnectorRequest
	next = GetDataById(doc, "redirectForm", "action")
	data = url.Values{}
	data.Set("token", GetDataById(doc, "token", "value"))
	doc = Request(next, data, false) //SpecificConnector/ServiceProvider

	// ***************** 5 Peticion ColleagueRequest *****************
	//	fmt.Println("...")
	next = GetDataByTagName(doc, "redirectForm", "action")
	data = url.Values{}
	data.Set("relayState", GetDataById(doc, "relayState", "value"))
	data.Set("SAMLRequest", GetDataById(doc, "SAMLRequest", "value"))
	doc = Request(next, data, false) //ColleagueRequest

	// ***************** 6 Peticion SpecificProxyService/ProxyServiceRequest *****************
	next = GetDataById(doc, "redirectForm", "action")
	data = url.Values{}
	data.Set("token", GetDataById(doc, "token", "value"))
	doc = Request(next, data, false) //SpecificProxyService/ProxyServiceRequest

	// ***************** 7 Ppeticion /SpecificProxyService/AfterCitizenConsentRequest *****************
	next = "/SpecificProxyService/" + GetDataById(doc, "consentSelector", "action")
	data = url.Values{}
	var atributes2 Atributes
	json.Unmarshal(GetDataSetByName("Person Informacion", eidasServer+":"+eidasPort), &atributes2)

	for _, item := range atributes2.GetAttributes.Attributes {
		data.Set(item.Name, item.Value)
	}
	data.Set("binaryLightToken", GetDataById(doc, "binaryLightToken", "value"))
	doc = Request(next, data, true)
	// ***************** 7 Ppeticion /SpecificProxyService/AfterCitizenConsentRequest *****************

	// ***************** 8 Peticion /IdP/AuthenticateCitizen *****************
	next = GetDataByTagName(doc, "redirectForm", "action")
	data = url.Values{}
	data.Set("SMSSPRequest", GetDataById(doc, "SMSSPRequest", "value"))
	doc = Request(next, data, false)
	// ***************** 8 Peticion /IdP/AuthenticateCitizen *****************

	// ***************** 9 Peticion /IdP/Response *****************
	next = "/IdP/" + GetDataById(doc, "authenticationForm", "action")
	data = url.Values{}
	data.Set("username", "xavi")
	data.Set("password", "creus")
	data.Set("eidasloa", "A")
	data.Set("checkBoxIpAddress", "on")
	data.Set("doNotmodifyTheResponse", "on")
	data.Set("signAssertion", "null")
	data.Set("encryptAssertion", "null")
	data.Set("smsspToken", GetDataByTagName(doc, "smsspToken", "value"))
	//data.Set("smsspToken", GetDataByTagName(doc, "smsspToken", "value"))
	data.Add("username", GetDataByTagName(doc, "username", "value"))
	data.Set("callback", GetDataByTagName(doc, "callback", "value"))
	aux := getElementById(doc, "jSonRequestDecoded")
	data.Set("jSonRequestDecoded", aux.FirstChild.Data)
	//fmt.Println(data)
	doc = Request(next, data, true)
	// ***************** 9 Peticion /IdP/Response *****************

	// ***************** /SpecificProxyService/IdpResponse *****************
	next = GetDataById(doc, "callback", "value")
	data = url.Values{}
	data.Set("errorMessage", "null")
	data.Set("errorMessageTitle", "null")
	data.Set("SMSSPResponse", GetDataById(doc, "jSonResponseEncoded", "value"))
	doc = Request(next, data, false)
	// ***************** /SpecificProxyService/IdpResponse *****************

	// 10 ***************** SpecificProxyService/AfterCitizenConsentResponse *****************
	next = "/SpecificProxyService/" + GetDataById(doc, "consentSelector", "action")
	data = url.Values{}
	data.Set("binaryLightToken", GetDataById(doc, "binaryLightToken", "value"))
	doc = Request(next, data, true)
	//fmt.Println(data)
	// 10 ***************** SpecificProxyService/AfterCitizenConsentResponse *****************

	//11 ***************** /EidasNode/SpecificProxyServiceResponse *****************
	next = GetDataById(doc, "redirectForm", "action")
	data = url.Values{}
	data.Set("token", GetDataById(doc, "token", "value"))
	doc = Request(next, data, false)
	//11 ***************** /EidasNode/SpecificProxyServiceResponse *****************

	//12 ***************** /EidasNode/ColleagueResponse *****************
	next = GetDataById(doc, "ColleagueResponse", "action")
	data = url.Values{}
	data.Set("RelayState", GetDataById(doc, "relayState", "value"))
	data.Set("SAMLResponse", GetDataByTagName(doc, "SAMLResponse", "value"))
	doc = Request(next, data, false)
	//12 ***************** /EidasNode/ColleagueResponse *****************

	//13 ***************** /SpecificConnector/ConnectorResponse *****************
	next = GetDataById(doc, "redirectForm", "action")
	data = url.Values{}
	data.Set("token", GetDataById(doc, "token", "value"))
	doc = Request(next, data, false)
	//13 ***************** /SpecificConnector/ConnectorResponse *****************

	//14 ***************** /SP/ReturnPage *****************
	next = GetDataByTagName(doc, "redirectForm", "action")
	next = eidasServer + ":" + eidasPort + "/SP/ReturnPage"
	data = url.Values{}
	data.Set("SMSSPResponse", GetDataById(doc, "SMSSPResponse", "value"))
	doc = Request(next, data, false)
	//14 ***************** /SP/ReturnPage *****************

	//15 ***************** /SP/populateReturnPage *****************
	next = "/SP/" + GetDataByTagName(doc, "countrySelector", "action")

	data = url.Values{}
	data.Set("SMSSPResponse", GetDataById(doc, "SMSSPResponse", "value"))
	doc = Request(next, data, true)

	//15 ***************** /SP/populateReturnPage *****************
	elapsed := time.Since(start).Seconds() // / time.Second

	//fmt.Println("t: ", (elapsed))
	//fmt.Println("Login Client ", ClientID, Successfullogin)
	return Successfullogin, elapsed
}

func GetDataById(doc *html.Node, selector string, atribute string) string {

	r1 := getElementById(doc, selector)
	//fmt.Printf("%+v\n", r1.Attr)
	//if selector == "jSonRequestDecoded" {
	//fmt.Printf("%+v\n", r1.Attr)

	var element string
	//fmt.Println(r1)
	for _, b := range r1.Attr {

		if b.Key == atribute {
			element = b.Val
		}

	}
	return element
}

func GetDataByTagName(doc *html.Node, selector string, atribute string) string {

	r1 := getElementByTagname(doc, selector)

	var element string
	//fmt.Println(r1)
	for _, b := range r1.Attr {

		if b.Key == atribute {
			element = b.Val
		}

	}
	return element
}

func Request(resource string, aux interface{}, addResource bool) *html.Node {

	data := aux.(url.Values)
	var u *url.URL
	if addResource == true {
		u, _ = url.ParseRequestURI(eidasServer + ":" + eidasPort)
		u.Path = resource
	} else {
		u, _ = url.ParseRequestURI(resource)
	}

	urlStr := u.String() // 'https://api.com/user/'
	//fmt.Println(urlStr)
	urlStr = strings.Replace(urlStr, "8080", eidasPort, 1)
	//fmt.Println(urlStr)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	//r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	resp, _ := client.Do(r)
	body, _ := ioutil.ReadAll(resp.Body)
	//"/EidasNode/ColleagueRequest"
	//
	if strings.Contains(resource, "SP/populateReturnPage") {
		if strings.Contains(string(body), "Login Succeeded") {
			Successfullogin = true
		}
		//fmt.Println(string(body))
	}

	doc, err := html.Parse(strings.NewReader((string(body))))
	if err != nil {
		panic("Fail to parse!")
	}
	return doc
}

func sendResult(resource string, aux interface{}) {

	data := aux.(url.Values)
	//var u *url.URL
	//u, _ = url.ParseRequestURI(resource)
	urlStr := resource //u.String() // 'https://urlServer'

	//fmt.Println(urlStr, data)

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(data.Encode())) // URL-encoded payload
	//r.Header.Add("Authorization", "auth_token=\"XXXXXXX\"")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(data.Encode())))

	client.Do(r)
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println((resp))

}

func GetAttribute(n *html.Node, key string) (string, bool) {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}

/*
func checkId(n *html.Node, id string) bool {
	if n.Type == html.ElementNode {
		s, ok := GetAttribute(n, "id")
		if ok && s == id {
			return true
		}
	}
	return false
}


func traverse(n *html.Node, id string) *html.Node {
	if checkId(n, id) {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result := traverse(c, id)
		if result != nil {
			return result
		}
	}

	return nil
}
*/

/******************************/

func checkElement(n *html.Node, name string, element string) bool {
	if n.Type == html.ElementNode {
		s, ok := GetAttribute(n, element)
		if ok && s == name {
			return true
		}
	}
	return false
}

func traverse(n *html.Node, name string, element string) *html.Node {
	if checkElement(n, name, element) {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result := traverse(c, name, element)
		if result != nil {
			return result
		}
	}

	return nil
}

func getElementByTagname(n *html.Node, name string) *html.Node {
	return traverse(n, name, "name")
}

func getElementById(n *html.Node, id string) *html.Node {
	return traverse(n, id, "id")
}
