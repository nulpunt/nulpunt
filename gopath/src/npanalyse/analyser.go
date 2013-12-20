package main

import (
	"github.com/GeertJohan/go.incremental"
	"labix.org/v2/mgo/bson"
	"log"
	"time"
)

var (
	analyserCount incremental.Uint
)

var (
	interval = 30 * time.Second
)

func initAnalysers(numAnalysers uint) {
	workChan := make(chan bson.ObjectId)
	doneChan := make(chan bool)

	processEndFuncs = append(processEndFuncs, func() {
		workLock.Lock()
		close(workChan)
	})

	// start some analysers
	for i := uint(0); i < numAnalysers; i++ {
		an := newAnalyser(workChan, doneChan)
		go an.work()
		processEndFuncs = append(processEndFuncs, func() {
			<-doneChan
		})
	}

	colUploads.Find(bson.M{""})
	//++ find if there are unprocessed documents
	//++ if not, sleep for a while, then retry
	//++ if yes, get workLock, lock the message, send the message on workChan, release workLock

}

// analyser handles uploads-analyse messages
type analyser struct {
	num      uint
	workChan chan bson.ObjectId
	jobCount incremental.Uint

	// when analyser is closing (workChan closed), should send a single bool on this chan.
	doneChan chan bool
}

func newAnalyser(workChan chan bson.ObjectId, doneChan chan bool) *analyser {
	return &analyser{
		num:      analyserCount.Next(),
		workChan: workChan,
		doneChan: doneChan,
	}
}

type uploadData struct {
	//++ stuff
}

func (an *analyser) work() {
	jobNum := an.jobCount.Next()
	for {
		uploadId, ok := <-an.workChan
		if !ok {
			log.Printf("workChan closed, worker %d stopped\n", an.num)
			an.doneChan <- true
			return
		}
		if flags.Verbose {
			log.Printf("Starting job %d-%d uploadId: %s\n", an.num, jobNum, uploadId)

		}

		upload := &uploadData{}
		colUploads.FindId(uploadId).One(upload)

		//++ process
	}
}
