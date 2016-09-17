//i changed the idea behind the jobs scheduler, instead of waiting for a main routine to assign the ID for the "job/subroutine"
//the main go routine will launch total_processors go rountines into the wild and let them do their own logic in retrieving
//an event ID and assign it to the shell script.
//if the id returned is 0 which is converted to -1 it means the list doesn't contain any IDs and the job is done.

package main
import (
  "github.com/mediocregopher/radix.v2/pool"
  "fmt"
  "sync"
  "time"
  "strconv"
  "os/exec"
)

const listname = "testlist"
const total_processors = 15

//logging success and failure per routine
type Log struct{
  success []string
  fail []string
}

func (l *Log) printSuccess(){
  fmt.Println(l.success)
}

func (l *Log) printFail(){
  fmt.Println(l.fail)
}


//struct to hold the connection pool - as recommended by the library for concurrency purposes-
//and wait group that will signal if all the routines are done.
type EventsCon struct{
  p *pool.Pool
  wg sync.WaitGroup
}

//creating a connection to redis
func (e *EventsCon) createNewConnection(){//this connection will be shared with all running threads
  var err error
  e.p, err = pool.New("tcp", "localhost:6379", 10)
  if err!=nil{
    fmt.Println("error construction a connection")
    return
  }
}

// retrieving next ID, this function makes sure to have a thread safe connection to redis
//since it uses pool.Cmd - based on radix.v2 documentation
func (e *EventsCon) getNextId() (int,error){


  response,err := e.p.Cmd("LPOP",listname).Int()
  if err!=nil{
    //fmt.Println("ERROR in RESPONSE")
    return -1,err
  }

  return response,nil
}


//the go routine that will be launched, it gets the EventsCon pointer and a jobLog pointer and it's own ID
//if the id of the event is -1 then no more events in the Redis list and it can call wg.Done() to signal
//that it's over.
func processJob(e *EventsCon, jobLog *Log, myId int){
  id,err:=e.getNextId()
  if err!=nil{
    if id==-1 {
      jobLog.success = append(jobLog.success,"no more IDs")
      fmt.Println("Routine "+strconv.Itoa(myId)+" says no more IDs")
    }else{
      jobLog.fail = append(jobLog.fail,"error getting the ID: "+err.Error()+" returned ID="+strconv.Itoa(id))
    }
  }

  for id!=-1 {

    fmt.Println("response is "+ strconv.Itoa(id) + " for rountine " +strconv.Itoa(myId))
    start:=time.Now()
    out, err := exec.Command("./foo.sh", strconv.Itoa(id)).Output()
    if err != nil {
      // TODO: capture this in a log file
      // log.Fatal(err)
      jobLog.fail = append(jobLog.fail, "event "+strconv.Itoa(id)+" has failed")
    } else {
      fmt.Printf("= %s\n", out)
      elapsed := time.Since(start)
      jobLog.success = append(jobLog.success, "event "+strconv.Itoa(id)+" has succeded in "+elapsed.String())
    }
    id,err=e.getNextId()
    if err!=nil{
      if id==-1 {
        jobLog.success = append(jobLog.success,"no more IDs")
        fmt.Println("Routine "+strconv.Itoa(myId)+" says no more IDs")
      }else{
        jobLog.fail = append(jobLog.fail,"error getting the ID: "+err.Error()+" returned ID="+strconv.Itoa(id))
      }
    }
  }
  e.wg.Done()
}

//launch total_processors go rountines then wait for all done signal close the connection
//then print the log
func main() {
  var logs [total_processors]Log
  events := EventsCon{}
  events.createNewConnection()
  events.wg.Add(total_processors)
  for i:=0; i<total_processors; i++ {
    logs[i] = Log{make([]string,0),make([]string,0)}
    go processJob(&events,&logs[i],i)
  }
  events.wg.Wait()
  events.p.Empty()
  for i:=0; i<total_processors; i++ {
    fmt.Println("routine "+strconv.Itoa(i))
    logs[i].printSuccess()
    logs[i].printFail()
  }
}
