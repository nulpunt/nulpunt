Server API for the Nulpunt server.
=========================

This document specifies the API between the nulpunt server and the
html front end.  In here you'll find every call that can be made from
the front end, the data structures, the parameters and the results.

    This document is leading. Any deviation between this document and
    the code is considered a bug. Either one (code or this document)
    needs to be changed.

This document does not specify data storage strucutures, See
database.md for that.

# Account management

## /sign-in

Parameters:
- username
- password

Returns: 
 - success message;
 - or error message;

Result depends on database records with a valid account an entry of
the corresponding parameters.

Side effects on the server:
- none.

## /register

Parameters:
- Username;
- Email address;
- Chosen Password;
- Chosen Color (RGB);

Returns:
- Success (account created);
- Error, (no account created);

Result depends on database records with a valid account an entry of
the corresponding parameters.

Side effects: 
- When the username has not been used, the account, with specified
password, email address and color is created in the database.

## /profile

To be defined.
Expected functions:
 - change password, old password, new password;
 - change tag-subscriptions, add tags, remove tags;

# Admin

The admin interface deals with uploading, processing and publishing documents.

## /admin/upload

Upload a document to be published.

Parameters:
- Document (the PDF contents, with a file name);

Returns:
- Ok, upload succeded.

Side effects:
- Saves the document in the storage. 
- It will be OCR-ed into:
  - images for each page;
  - coordinates of each character;
  - asynchronously. (it can take a while).

## /admin/process

Add metadata to an uploaded document

### GET /admin/process/list

Returns a list of unpublished documents.

For each document expect these fields:
- docId; internal id of the document. needed for reference;
- original file name; File name as it was when it was uploaded;
- timestamp of upload;

### GET /admin/process

Parameter: 
- docId, the internal id of the document;

Returns the selected document, together with any existing metadata, see POST parameters;

User can edit all parameters.

TODO: invent something to correct OCR-errors.

### POST /admin/updateDocument

Updates the metadata of a document.

POST Parameters:
- To be defined. Examples: 
  - title;
  - deparment, author, subject; 
  - tags;
  - dates;
  - whatever;
  - Publish Yes/No;
  - Delete Yes/No;	  
Result:
	Updated document.

Side effects:
- The metadata of the document gets updated with the specified values;
- if Published == yes, document will become visible on the site. No: remove from site;
- if Delete == yes, document, all metadata and any comments will be deleted from the database.

## admin/analyse 

This gets removed. 

Rationale: Documents get added to a queue for OCR'ing after uploading. OCR'ing happens automatically. 
When OCR'ed succesfully, documents get visible in the /process list.

## admin/tags

A page devoted to managing the list of tags to assign to documents.

GET retrieve the list of tags,
It returns a list of JSON-encoded 

    [{ID: 123abc, Tag: 'example'}, ..]

With status code 200.

POST adds a tags
It takes a JSON encoded object in the request-body:

    { Tag: 'example' }

On succes, it returns a 200-status code and the new list in the same
way as the GET request. Clients can use this to update their view.

On error, it returns a 400/405/500 status code with a plain text string in the page body.

## admin/tags/delete 
POST deletes a tag from the list. It cannot be selected anymore for new taggings.

It takes the same parameter as POST admin/tags call.
It returns the same types of results. 

NOTE: Documents tagged with it stay as they are. IE, Tags are used by
value in the document classification, not by reference.


# Document viewing

## GET /service/getDocumentList

ListDocuments returns a list of documents that match the specified criteria.

Criteria are specified in bson.

## GET /service/getDocument

Parameters:
- docID
- annotationID
- commentID

This shows the document with, the selected page and the
selected annotation and the comments.

It is designed to be the full, static URL of the
document with the annotation and comment on the page.

It's for static deep-linking. People can post this URL everywhere and
be sure other readers can read their annotation and comment on the document.

The $commentid is optional. Without it, it shows the document on the page with requested annotations.

The $annotation is optional. Without it, it returns the document with the first annotation.

Returns:
- document record;
- annotation-record;
- comment-records;
- pages-record

Clients need to fetch the page image data in a separate call to getImage

Side effects (on the server): none.


## GET /service/getImage/$docId/$pagenumber.png

Fetches the document page image data. It's static data, so ideal for caching at clients.

Parameters: none, it's in the URL.

Result: the image data in a http-body.

TODO: cache this stuff at apache level.

## POST /sevice/session/add-annotation

Add a quote(selection) and comment for the world to see. IE, people
can add a selection of a document and their comments.

Parameters:
- documentId
- one or more ranges of start-end coordinates;
- commentary text;

Server adds the annotatorId from the session and stores it in the database.

Result:
- Ok, added, gives bookmarkable, static URL to the document with the comment;
- Error

Side effects: 
- when valid: add the annotation to the database;

Requirements:
- be logged in. (we need to know who you are).


## POST /service/session/add-comment

Add a reponse to an annotation/comment
The goal is to provide a way to add a comment to an existing annotation.
It allows people to discuss that annotation

Parameters:
- documentid
- annotationid;
- parentid to which this is a reply; nil for first comment or for no threading;
- Response-text

Server adds the commenterID from the session.

Returns:
- OK/Error with the full URL to the comment.

Side Effects:
- Add a comment to an existing quotation. Sorted by submission date.

## GET /service/getPage

Gets a page-record of a document. Needed for lazy loading.

Parameters:
- documentid
- page number. Starts at 1. (there is no page 0).

Returns: 
- page-record
- []annotations on the page

## GET /service/getComments

Gets the comment for the page. Used in lazy loading

Parameters: 
- documentId
- annotationId

Returns: comments on the page (to be defined further).

# Document selection and ordering

This part deals with document selection and ordering.

## GET /trending

This retrieves a list of documents that are sorted to the  'trending' criterium.

Trending can be defined simply as 'ordered by timestamp of latest annotation'. It will be a 'jumpy' list.

Or more complex as weighted number of anntations and comments in the last X minutes.

The server keeps a mostly static list of trending documents. This to keep the amount of work at request time to a minimum. This is really important for this call, as it will be used for every request at the home page. IE, everyone.

Clients can request parts of the list that they are interested in. (to make lazy loading possible)

Parameters: 
- startAt; position in the list where to start; 0 is the most trending one.
- limit; max number of documents to return
- nrOfAnnotations; number of recent annotations per trending document

Returns:
- a list of documents. (max 50)
  For each document expect: 
  - document-record (as in the database.md)
  - []annotation-records (as in db)
  - []pages-records (as in db) that are specified in the annotations
  - []crops of the images, one for each annotation.

This response contains all that's needed to show the trending page.

Expected use:
* Front end code starts with a call to GET /trending with startAt 0, limit to what is wanted (1, 3, ...) and NrOfAnnotations to what's wanted.
* Front end displays page;
* user presses buttons on the pages to see more (older) trending data. Can load lazily, just change the start-at parameter.
