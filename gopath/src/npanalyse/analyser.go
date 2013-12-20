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

	// find work
	go func() {
		for {
			updateID := &struct {
				ID bson.ObjectId `bson:"_id"`
			}{}
			err := colUploads.Find(bson.M{"analyseState": ""}).Select(bson.M{"_id": 1}).One(updateID)
			if err != nil {
				log.Printf("error searching for non-analysed update: %d\n", err)
				goto Sleep
			}
			if updateID.ID != "" {
				workLock.Lock()
				err := colUploads.UpdateId(updateID.ID, bson.M{"$set": bson.M{"analyseState": "started"}})
				if err != nil {
					log.Printf("error setting analyseState for upload %s to 'started'\n", updateID.ID)
					workLock.Unlock()
					continue
				}
				workChan <- updateID.ID
				workLock.Unlock()
				// try to find next job right away
				continue
			}
		Sleep:
			time.Sleep(30 * time.Second)
		}
	}()
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
