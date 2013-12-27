package main

import (
	"bytes"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"github.com/GeertJohan/go.incremental"
	"github.com/GeertJohan/go.leptonica"
	"github.com/GeertJohan/go.tesseract"
	"github.com/nfnt/resize"
	"image"
	"image/png"
	"io"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const docviewerWidth = 780
const thumbnailWidth = 100

var (
	pollInterval   = 2 * time.Second
	instanceUnique = "" // filled by init()

	analyserCount incremental.Uint

	regexpOutputFileName = regexp.MustCompile(`^output-[0]*([0-9]+).png$`)
)

func init() {
	// setup instanceUnique
	unixNano := time.Now().UnixNano()
	varint := make([]byte, 10)
	n := binary.PutVarint(varint, unixNano)
	base32Unique := base32.StdEncoding.EncodeToString(varint[:n])
	instanceUnique = strings.Replace(base32Unique, "=", "", -1)
	log.Printf("generated unique instance code %s\n", instanceUnique)
}

func initAnalysers(numAnalysers uint) {
	workChan := make(chan bson.ObjectId)
	var workLock sync.Mutex
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
			documentIDHolder := &struct {
				ID bson.ObjectId `bson:"_id"`
			}{}
			err := colDocuments.Find(bson.M{"analyseState": "uploaded"}).Select(bson.M{"_id": 1}).One(documentIDHolder)
			if err != nil {
				if err != mgo.ErrNotFound {
					log.Printf("error searching for non-analysed update: %s\n", err)
				}
				goto Sleep
			}
			if documentIDHolder.ID != "" {
				workLock.Lock()
				err := colDocuments.UpdateId(documentIDHolder.ID, bson.M{"$set": bson.M{"analyseState": "started"}})
				if err != nil {
					log.Printf("error setting analyseState for upload %s to 'started'\n", documentIDHolder.ID)
					workLock.Unlock()
					continue
				}
				workChan <- documentIDHolder.ID
				workLock.Unlock()
				// try to find next job right away
				continue
			}
		Sleep:
			time.Sleep(pollInterval)
		}
	}()
}

// analyser usernames uploads-analyse messages
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

type documentData struct {
	UploadFilename     string    `bson:"uploadFilename"`
	UploadGridFilename string    `bson:"uploadGridFilename"`
	UploadDate         time.Time `bson:"uploadDate"`
	UploaderUsername   string    `bson:"uploaderUsername"`
	Language           string    `bson:"language"`
	Title              string    `bson:"title"`
	PageCount          int       `bson:"pageCount"`
	AnalyseState       string    `bson:"analyseState"`
}

type pageData struct {
	ID            bson.ObjectId  `bson:"_id"`
	DocumentID    bson.ObjectId  `bson:"documentId"` // refers to `documents._id`)
	PageNumber    uint           `bson:"pageNumber"` // page number
	Lines         []*[]*charData `bson:"lines"`
	Text          string         `bson:"text"` //  the text in the same order as the lines-attribute, use for search/sharing. Contains ocr-errors
	HighresWidth  int            `bson:"highresWidth"`
	HighresHeight int            `bson:"highresHeight"`
}

type charData struct {
	X1 float32 `bson:"x1"` // offset-left in pixels
	Y1 float32 `bson:"y1"` // offset-top in pixels
	X2 float32 `bson:"x2"` // offset-bottom in pixels
	Y2 float32 `bson:"y2"` // offset-right in pixels
	C  string  `bson:"c"`  // character
}

func (an *analyser) work() {
	for {
		documentID, ok := <-an.workChan
		jobNum := an.jobCount.Next()
		if !ok {
			log.Printf("workChan closed, worker %d stopped\n", an.num)
			an.doneChan <- true
			return
		}

		an.Logf("starting job %d-%d documentID: %s", an.num, jobNum, documentID.Hex())

		document := &documentData{}
		err := colDocuments.FindId(documentID).One(document)
		if err != nil {
			log.Printf("error analysing doc %s: %s\n", documentID, err)
			continue
		}

		var tessLanguage string
		switch document.Language {
		case "nl_NL", "":
			tessLanguage = "nld"
		case "en_EN":
			tessLanguage = "eng"
		default:
			log.Printf("error invalid language '%s' for document %s\n", document.Language, documentID.Hex())
			continue
		}
		an.Logf("tesseract language: %s", tessLanguage)

		func() {
			an.Logf("docID: %s", documentID.Hex())
			//++ defer a function that checks if this func was successfull (update with updateId has analyseState "completed")
			//++ when was not successfull, set state to error, remove any pages with documentId

			tmpDirName := fmt.Sprintf("/tmp/npanalyse-%s-%d-%d", instanceUnique, an.num, jobNum)
			err = os.Mkdir(tmpDirName, 0774)
			if err != nil {
				log.Printf("failed to create tmp dir '%s': %s\n", tmpDirName, err)
				return
			}
			an.Logf("created tmp dir %s", tmpDirName)
			defer func() {
				// clean up temp dir
				err = os.RemoveAll(tmpDirName)
				if err != nil {
					log.Printf("error cleaning up tmp dir %s: %s\n", tmpDirName, err)
				}
				an.Logf("cleaning up tmp dir %s", tmpDirName)
			}()

			originalFileGridFS, err := gridFS.Open(document.UploadGridFilename)
			if err != nil {
				log.Printf("error opening original file (%s) from GridFS: %s\n", document.UploadGridFilename, err)
				return
			}
			defer originalFileGridFS.Close()
			originalFileTmp, err := os.Create(path.Join(tmpDirName, "original.pdf"))
			if err != nil {
				log.Printf("error creating original file in tmpDir: %s\n", err)
				return
			}
			defer originalFileTmp.Close()

			// copy contents
			_, err = io.Copy(originalFileTmp, originalFileGridFS)
			if err != nil {
				log.Printf("error copying data from gridFS to tmp file: %s\n", err)
				return
			}
			originalFileGridFS.Close()

			// convert pdf to png's
			pdftoppmPath, err := exec.LookPath("pdftoppm")
			if err != nil {
				log.Printf("failed to find `pdftoppm` in PATH, is it even installed? err: %s\n", err)
				return
			}
			pdftoppm := exec.Cmd{
				Path: pdftoppmPath,
				Args: []string{
					"pdftoppm",
					"-r", "900",
					"-png",
					"original.pdf",
					"output",
				},
				Dir:    tmpDirName,
				Stdout: os.Stdout,
				Stderr: os.Stderr,
			}

			err = pdftoppm.Run()
			if err != nil {
				log.Printf("error running pdftoppm: %s\n", err)
				return
			}
			an.Logf("ran pdftoppm for document %s", documentID.Hex())

			tess, err := tesseract.NewTess("/usr/share/tesseract-ocr/tessdata/", tessLanguage)
			if err != nil {
				log.Printf("error creating new tesseract instance: %s\n", err)
				return
			}
			defer tess.Close()

			tmpDir, err := os.Open(tmpDirName)
			if err != nil {
				log.Printf("error opening tmpDir(%s): %s\n", tmpDirName, err)
				return
			}

			fileInfos, err := tmpDir.Readdir(0)
			if err != nil {
				log.Printf("error reading tmpDir(%s): %s\n", tmpDirName, err)
				return
			}
			var fileNames []string
			var fileInfosByName = make(map[string]os.FileInfo)
			for _, fileInfo := range fileInfos {
				fileName := fileInfo.Name()
				if regexpOutputFileName.MatchString(fileName) {
					fileInfosByName[fileName] = fileInfo
					fileNames = append(fileNames, fileName)
				}
			}
			sort.Strings(fileNames)
			for _, fileName := range fileNames {
				success := an.analyseFile(documentID, tess, tmpDirName, fileInfosByName[fileName])
				runtime.GC()
				if !success {
					return
				}
			}
			document.PageCount = len(fileNames)
			document.AnalyseState = "completed"
			document.Title = document.UploadFilename
			err = colDocuments.UpdateId(documentID, bson.M{"$set": document})
			if err != nil {
				log.Printf("error inserting document: %s\n", err)
				return
			}
			an.Logf("inserted document %s", documentID.Hex())
		}()
	}
}

func (an *analyser) analyseFile(documentID bson.ObjectId, tess *tesseract.Tess, tmpDirName string, fileInfo os.FileInfo) bool {
	pageNumberSubmatch := regexpOutputFileName.FindStringSubmatch(fileInfo.Name())
	pageNumberString := pageNumberSubmatch[1]
	pageNumberUint64, _ := strconv.ParseUint(pageNumberString, 10, 32)
	pageNumber := uint(pageNumberUint64)
	an.Logf("found page %d", pageNumber)

	outputTmpFile, err := os.Open(path.Join(tmpDirName, fileInfo.Name()))
	if err != nil {
		log.Printf("error opening output file(%s): %s\n", fileInfo.Name(), err)
		return false
	}
	defer outputTmpFile.Close()

	outputGridFileHighresName := fmt.Sprintf("highres/%s-%s.png", documentID.Hex(), pageNumberString)
	outputGridFileHighres, err := gridFS.Create(outputGridFileHighresName)
	if err != nil {
		log.Printf("error creating GridFS file(%s): %s\n", outputGridFileHighresName, err)
		return false
	}
	defer outputGridFileHighres.Close()

	// create buffer to be filled with image data
	imageBuf := bytes.NewBuffer(make([]byte, 0, fileInfo.Size()))

	// copy image data to gridFile while tee-reading to imageBuf
	_, err = io.Copy(outputGridFileHighres, io.TeeReader(outputTmpFile, imageBuf))
	if err != nil {
		log.Printf("error copying data from tempFile to gridFile: %s\n", err)
		return false
	}
	outputGridFileHighres.Close()
	an.Logf("read output png, saved highres. page %d", pageNumber)

	// get bytes from imageBuf and create leptonica pix
	imageBytes := imageBuf.Bytes()
	pix, err := leptonica.NewPixReadMem(&imageBytes)
	if err != nil {
		log.Printf("error creating new pix from imageBuf: %s\n", err)
		return false
	}
	defer pix.Close()

	// resize for thumbnail
	if pageNumber == 1 {
		imageBufReader := bytes.NewReader(imageBuf.Bytes())
		outputGridFileThumbnailName := fmt.Sprintf("document-thumbnails/%s.png", documentID.Hex())
		outputGridFileThumbnail, err := gridFS.Create(outputGridFileThumbnailName)
		if err != nil {
			log.Printf("error creating GridFS file(%s): %s\n", outputGridFileThumbnailName, err)
			return false
		}
		defer outputGridFileThumbnail.Close()
		err, _ = readResizeWrite(imageBufReader, outputGridFileThumbnail, thumbnailWidth)
		if err != nil {
			log.Printf("error performing readResizeWrite for gridFile(%s): %s\n", outputGridFileThumbnailName, err)
			return false
		}
		outputGridFileThumbnail.Close()
		an.Logf("resized page for thumbnail %d", pageNumber)
	}

	// resize for docviewer
	imageBufReader := bytes.NewReader(imageBuf.Bytes())
	outputGridFileDocviewerName := fmt.Sprintf("docviewer-pages/%s-%s.png", documentID.Hex(), pageNumberString)
	outputGridFileDocviewer, err := gridFS.Create(outputGridFileDocviewerName)
	if err != nil {
		log.Printf("error creating GridFS file(%s): %s\n", outputGridFileDocviewerName, err)
		return false
	}
	defer outputGridFileDocviewer.Close()
	err, sizes := readResizeWrite(imageBufReader, outputGridFileDocviewer, docviewerWidth)
	if err != nil {
		log.Printf("error performing readResizeWrite for gridFile(%s): %s\n", outputGridFileDocviewerName, err)
		return false
	}
	outputGridFileDocviewer.Close()
	an.Logf("resized page for docviewer %d", pageNumber)

	// hand leptonica pix to tess
	tess.SetImagePix(pix)

	// create page object
	page := &pageData{
		ID:            bson.NewObjectId(),
		DocumentID:    documentID,
		PageNumber:    pageNumber,
		Text:          tess.Text(),
		Lines:         make([]*[]*charData, 0),
		HighresWidth:  sizes.Dx(),
		HighresHeight: sizes.Dy(),
	}

	// get boxed text
	boxText, err := tess.BoxText(0)
	if err != nil {
		log.Printf("error retrieving boxText: %s\n", err)
		return false
	}
	// cleanup tess and pix for this page
	tess.Clear()
	pix.Close()
	an.Logf("retrieved boxText for page %d", pageNumber)

	// loop over box text and create lines
	var line []*charData
	for _, tessChar := range boxText.Characters {
		char := &charData{
			X1: (float32(tessChar.StartX) / float32(page.HighresWidth) * float32(100)),
			Y1: (float32(tessChar.StartY) / float32(page.HighresHeight) * float32(100)),
			X2: (float32(tessChar.EndX) / float32(page.HighresWidth) * float32(100)),
			Y2: (float32(tessChar.EndY) / float32(page.HighresHeight) * float32(100)),
			C:  string(tessChar.Character),
		}
		//TODO: \n won't ever happen with BoxText()
		// ++ need to mix this information with .Text() information to have whitespace
		if line == nil || char.C == "\n" {
			line = make([]*charData, 0)
			page.Lines = append(page.Lines, &line)
		}

		line = append(line, char)
	}

	err = colPages.Insert(page)
	if err != nil {
		log.Printf("error inserting page into collection: %s\n", err)
		return false
	}
	an.Logf("inserted page %s", page.ID.Hex())
	return true
}

func (an *analyser) Logf(format string, stuff ...interface{}) {
	if flags.Verbose {
		log.Printf(fmt.Sprintf("%d-%d: %s\n", an.num, an.jobCount.Last(), format), stuff...)
	}
}

func readResizeWrite(imageBuf io.Reader, to io.Writer, width uint) (error, *image.Rectangle) {
	img, err := png.Decode(imageBuf)
	if err != nil {
		return err, nil
	}
	imgResized := resize.Resize(width, 0, img, resize.MitchellNetravali)
	err = png.Encode(to, imgResized)
	if err != nil {
		return err, nil
	}
	sizes := img.Bounds()
	return nil, &sizes
}
