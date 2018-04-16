package main

// Create a separate file in the same directory as main called key.go with
// package main
// var KEY = "<key generated at http://datamine.mta.info/>"
//
// Then run with go run main.go key.go

import (
  "fmt"
  "io/ioutil"
  "log"
  "net/http"
  "sort"
  "time"

  "github.com/golang/protobuf/proto"

  "real_time_subway/transit_realtime"
)

var FEED = 21  // http://datamine.mta.info/list-of-feeds. 21 is the BDFM lines.

var TRAIN = "F"

// http://web.mta.info/developers/data/nyct/subway/Stations.csv  <-- GTFS Stop ID

// TODO: check north and south at once.
var STATION = "F21N"


// Go has no int64 sort, so this lets us specify our own.
type timeStamps []int64
func (x timeStamps) Len() int { return len(x) }
func (x timeStamps) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x timeStamps) Less(i, j int) bool { return x[i] < x[j] }


func main() {

  url := fmt.Sprintf(
  "http://datamine.mta.info/mta_esi.php?key=%s&feed_id=%d", KEY, FEED)

  client := &http.Client {}
  req, err := http.NewRequest("GET", url, nil)
  req.Header.Set("Content-Type", "application/octet-stream")


  resp, err := client.Get(url)
  if resp.StatusCode != 200 {
    log.Fatal("Bad response code from %s. Got %d: %s.", url, resp.StatusCode,
              resp.Status)
  }

  if err != nil {
    log.Fatal(err)
  }

  data, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    log.Fatal(err)
  }

  // Bad access still returns a 200, so add this kind of gross check.
  if string(data) == "You have no access." {
    log.Fatal("No access. Is your key correct?")
  }

  newFeed := &transit_realtime.FeedMessage{}
  err = proto.Unmarshal(data, newFeed)
  if err != nil {
    log.Fatal("Couldn't unmarshal the data feed:", err)
  }

  header := newFeed.GetHeader()
  fmt.Printf("Feed version: %s; ", *(header.GtfsRealtimeVersion))
  // ...this can't be the best way to unpack the time, can it?
  fmt.Printf("last updated: %v\n", time.Unix(int64(*(header.Timestamp)), 0))


  times := timeStamps{}
  entities := newFeed.GetEntity()
  for _, entity := range entities {
    update := entity.GetTripUpdate()
    if update != nil && update.GetTrip().GetRouteId() == TRAIN {
      for _, stop_time := range update.GetStopTimeUpdate() {
        if *stop_time.StopId == STATION {
          arrival := stop_time.GetArrival()
          times = append(times, arrival.GetTime())
        }
      }
    }
  }

  sort.Sort(times)
  for _, eventTime := range times {
    // TODO: print how many minutes.
    fmt.Printf("%s\n", time.Unix(eventTime, 0))
  }
}

