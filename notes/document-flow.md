Document Flow
============

This describes the document-flow through the nulpunt server.

### Uploading & Analyzing
- Admins (users with sufficient rights) can upload documents in a 'bulk' way on the 'upload' admin page.
- The server stores the pdf file into GridFS (see [database.md](database.md) for more info) and creates an entry in the `uploads` collection.
- Uploaded documents are visible for uploader only on his/her 'analyse' admin page.
- Admin inputs meta-data for each uploaded document on the analyse page.
 - text's language (for ocr purposes).
 - color or grayscale document.
 - summary? extra information about 'why' it was uploaded?
 - original source? etc?
- The server updates the mongo-document in the `uploads` collection with the upload-metadata and marks it ready for ocr.

### OCR
The `npanalyse` application analyses uploaded data. To let `npanalyse` do it's job, `npserver` must invoke it somehow. To not re-invent the wheel, we use NSQ, a Message Queue written in Go at bitly. Read up on NSQ [here](http://bitly.github.io/nsq/).

The NSQ topic for the upload-job is: `uploads`
The NSQ channel for the `npanalyse` consumer is: `analyse`

During development, a single `nsqlookupd` and a single `nsqd` instance are ran. We can scale when required.

The message sent on the uploads topic should be a json encoded data-structure with only one field (for starters):
- field:`uploadId` type:`bson.ObjectId (hex string)`

`npanalyse` reads messages from the Queue and performs the following:
- extracts/converts image data from the pdf, one image per page
- ocr's each page
- create 1 document entry in the document collection (to get the documentID for referencing to in the other collections) documents have their published-flag set to false
- create 1 page-record for each page

[More info on database collections](/notes/database.md).

When `npanalyse` succesfully analysed a document, it deletes the entry in the `uploads` collection.

### Metadata attachment and publishing
The admininstrator browses the list of documents that have their published-flag set to false
- The admin selects one
- It comes up on screeen
- Admin adds the metadata;
- Admin sets published-flag to true;
- Admin saves it. Now it is published.

