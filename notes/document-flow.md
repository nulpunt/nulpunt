Document Flow
============

This describes the document-flow through the nulpunt server.

# uploading

Admins (users with sufficient rights) can upload documents.
Uploading requires the user to specify the language.
It stores the pdf into gridfs
It creates an Upload-record with the upload-metadata (pretend it to be a queue)

# ocr-ing

An independent proces:
- takes the top of the queue, (upload-table)
- extracts/converts image data from the pdf, one image per page;
- ocr's each page;

The process creates:
- 1 document entry in the document table (to get the documentID for referencing to in the other tables)
    documents have their published-flag set to false;
- 1 page-record for each page

If successful, the process deletes the entry in the upload-table (remove from queue)

# metadata attachment and publishing

The admininstrator browses the list of documents that have their published-flag set to false
- The admin selects one
- It comes up on screeen
- Admin adds the metadata;
- Admin sets published-flag to true;
- Admin saves it. Now it is published.

