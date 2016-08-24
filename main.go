/*Simple Go Task
 */

package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"time"
	)


type ResponseList struct {
	Responses []Response
}

//Type holds our parsed response
type Response struct{
	Time string `json:time`
	Data []FieldInfo `json:data`
}

//FieldInfo type 
type FieldInfo struct {
	Symbol string `json:symbol`
	Series string `json:series`
	OpenPrice string `json:openPrice`
	HighPrice string `json:highPrice`
	LowPrice string `json:lowPrice`
	Ltp string `json:ltp`
	PreviousPrice string `json:previousPrice`
	NetPrice string `json:netPrice`
	TradedQuantity string `json:tradedQuantity`
	TurnoverInLakhs string `json:turnoverInLakhs`
	LastCorpAnnouncementDate string `json:lastCorpAnnouncementDate`
	LastCorpAnnouncement string `json:lastCorpAnnouncement`
}


//Helper to build HTML page string
func make_html(d *Response) string {
	
	//Initial stuff
	initial := `<!Doctype html>
<html>
<head>
<title>SIMPLE GO TASK</title>
</head>
<body>` + `<div align="center">Snapshot Time: `+ d.Time + "</br></br></br>"


	//Initial table stuff
	table := `<table "width:100%">
	<tr>
	<th>Symbol</th>
	<th>Series</th>
	<th>OpenPrice</th>
	<th>HighPrice</th>
	<th>LowPrice</th>
	<th>Ltp</th>
	<th>PreviousPrice</th>
	<th>NetPrice</th>
	<th>TradedQunatity</th>
	<th>TurnoverInLakhs</th>
	<th>LastCorpAnnouncementDate</th>
	<th>LastCorpAnnouncement</th></tr>`
	
	var row string
	//Iterate over our Data and make appropriate rows for the table
	for _,v:= range d.Data {
	row_helper := `
	<tr>
	<td>%s</td>
	<td>%s</td>
	<td>%s</td>
	<td>%s</td>
	<td>%s</td>
	<td>%s</td>
	<td>%s</td>
	<td>%s</td>
	<td>%s</td>
	<td>%s</td>
	<td>%s</td>
	<td>%s</td></tr>`	
	
	row = row + "\n" + fmt.Sprintf(row_helper, v.Symbol,v.Series,v.OpenPrice, v.HighPrice, v.LowPrice, v.Ltp, v.PreviousPrice, v.NetPrice, v.TradedQuantity,v.TurnoverInLakhs, v.LastCorpAnnouncementDate, v.LastCorpAnnouncement)
	}
	
	//Close table
	row = row + "</table>"
	
	//Final stuff
	final := "</body>\n</html>"
	
	return initial + table + row + final
}


//get_stuff() handles the crawling. Returns type Response which has our data
func get_stuff(d *Response) {

	//get a Client type
	scraper := &http.Client{}
	
	//Build a request
	request, err := http.NewRequest("GET", "https://www.nseindia.com/live_market/dynaContent/live_analysis/gainers/niftyGainers1.json", nil)
	
	//Some Headers
	request.Header.Add("User-Agent","Mozilla/5.0(X11; Linux x86_64; rv:47.0) Gecko/20100101 Firefox/47.0")
	request.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	
	//Make the request
	response, err := scraper.Do(request)
	
	if err != nil {
		panic(err)
	}
	
	defer response.Body.Close()
	
	//Parse into byte array
	data,err := ioutil.ReadAll(response.Body)
	if err!= nil {
		panic(err)
	}
	
	//Parse into our struct Response
	err = json.Unmarshal(data, d)
	
	//if maintaining snapshot list then uncomment and change function definition
	//list.Responses = append(list.Responses, *d)
}

//Simple handler for displaying stuff	
func handler(w http.ResponseWriter, r *http.Request, stuff string) {
	fmt.Fprintf(w, fmt.Sprintf("%s",stuff), r.URL.Path[1:])
}


	
	
func main() {

	//ticker returns a type which has a channel with ticks for events after every d duration
	ticker := time.NewTicker(300 * time.Second) //300 seconds =  5mins
	quit := make(chan struct{}) //channel of struct{} to bail out of the ticker. Not necessary for perpetually running service
	
	//response_list := new(ResponseList)	
	d := new(Response)

	//Routine that triggers a crawl every 5 minutes. ticker returns a channel ticker.C for event every duration d
	go func () {
		
		for {
			select {
				case <- ticker.C :
					//spawns another routine to accomodate for quick ticks. Keeping it generic
					go get_stuff(d)
				
				case <- quit:
					ticker.Stop()
			}
		}
	}()

	//catch /snap URl and show current snapshot
	http.HandleFunc("/snap", func(w http.ResponseWriter, r *http.Request) {
							stuff:= make_html(d)
							handler(w,r,stuff)
							})
	
	
	//Maintains a list of all snapshots since program starts.
	//http.HandleFunc("/snapshots", func(w http.ResponseWriter, r *http.Request){
								
								//output, err := json.Marshal(response_list)
								//fmt.Println(response_list)
								//if err != nil {
									//panic(err)
								//}
								//handler(w,r,string(output))
	
	
					//})

	//serve
	fmt.Println("Starting server....")
	http.ListenAndServe(":5000", nil)
	
}

	
	
	


