Server API for the Nulpunt server.
=========================

This document specifies the API between the nulpunt server and the
html front end.  In here you'll find every call that can be made from
the front end, the data structures, the parameters and the results.

    This document is leading. Any deviation between this document and
    the code is considered a bug. Either one (code or this document)
    needs to be changed.

This document does not specify data storage strucutures, See
Datastore-design.md for that.

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

### GET /admin/process/doc?docId

Parameter: 
- docId, the internal id of the document;

Returns the selected document, together with any existing metadata, see POST parameters;

User can edit all parameters.

TODO: invent something to correct OCR-errors.

### POST 

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
value in the document classificiation, not by reference.


# Document viewing

## GET /document/$docId/#commentid

This shows the document (or the first part), the selected page and the
selected comment.  It is designed to be the full, static URL of the
document with the comment on the page.

It's for static deep-linking. People can post this URL everywhere and
be sure other readers can read their comment on the document.

The #commentid is optional. Without it, it shows the first page/all pages.

Parameters:
- none;  it's in the URL

Result:
- The page of the document: a html page with the
          image(gif/png) of the specified page and centered on the
          specified comment;
- The Javascript has to fetch the details of the selection and
          the contents of the comments in a separate request;
- Details include: comment, details of commenter: name, color.

Side effects (on the server): none.

## POST /document/$docid

Add a quote(selection) and comment for the world to see. IE, people
can add a selection of a document and their comments.

Parameters:
	- one or more ranges of start-end coordinates;
	- commentary text;

Result:
- Ok, added, gives bookmarkable, static URL to the document with the comment;
- Error

Side effects: 
- when valid: add the comment to the database;

Requirements:
- be logged in. (we need to know who you are).


## POST /document/$docid/$commentid

Add a reponse to a selection/comment
The idea is to provide a way to add a comment to an existing quotation.
It allows people to discuss a certain quotation.

Parameters:
- Response-text

Returns:
- OK/Error

Side Effects:
- Add a comment to an existing quotation. Sorted by submission date.

# Document selection and ordering

This part deals with document selection and ordering.

## GET /trending

This retrieve a list of docuement that are sorted to the  'trending' criterium.

Trending can be defined simply as 'ordered by timestamp of latest annotation'. It will be a 'jumpy' list.

Or more complex as weighted number of anntations and comments in the last X minutes.

Parameters: 
- limit; max number of documents to return

Returns:
- a list of documents. (max 50)
  For each document expect: 
  - docId;
  - title;
  - summary (if not too long);
  - latest annotation-id

Front end code needs to fetch page contents (image, overlay) and the annotation and comments to display it all.

Perhaps we should consider making this call a special case that precomputes the pages and makes it mostly static, as it is used on the front page. 