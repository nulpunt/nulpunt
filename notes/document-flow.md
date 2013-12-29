Document Flow
============

This describes the document-flow through the nulpunt server.

### Uploading & Analyzing
- Admins (users with sufficient rights) can upload documents in a 'bulk' way on the 'upload' admin page.
- Admin select a text's language (for ocr purposes) before uploading documents.
- The server stores the pdf file into GridFS (see [database.md](database.md) for more info) and creates an entry in the `documents` collection with given language and upload information (gridFS location, etc.).

### OCR
The `npanalyse` application analyses uploaded data. To let `npanalyse` do it's job, `npserver` must invoke it somehow. To not re-invent the wheel, we use NSQ, a Message Queue written in Go at bitly. Read up on NSQ [here](http://bitly.github.io/nsq/).

`npanalyse` reads messages from the Queue and performs the following:
- extracts/converts image data from the pdf, one image per page
- ocr's each page
- update the document entry in the `documents` collection
- create one page-record for each page

[More info on database collections](/notes/database.md).

When `npanalyse` succesfully analysed a document, it deletes the entry in the `uploads` collection.

### Metadata attachment and publishing
The admininstrator browses the list of documents that have their published-flag set to false
- The admin selects one
- It comes up on screeen
- Admin adds the metadata;
- Admin can always 'save' the document.
- Admin can only 'save+publish' the document when it is analysed (`analyseState` = 'completed')
