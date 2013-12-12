Server API for the Nulpunt server.
=========================

This document specifies the API between the nulpunt server and the
html front end.  In here you'll find every call that can be made from
the front end, the data structures, the parameters and the results.

This document is leading. Any deviation between this document and the
current code is considered a bug. Either one (code or this document)
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
     - When the username has not been used, the account, with
       specified password, email address and color is created in the
       database.

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
       - coordinates of each characater;
       - asynchronously. (it can take a while).

## /admin/process

Add metadata to an uploaded document
GET gives a list of document to be processed.
POST updates a document.

POST Parameters:
	- To be defined. Examples: 
	  - title;
	  - deparment, author, subject; 
	  - tags;
	  - dates;
	  - whatever;
	  - Publish Yes/No;
	  
Result:
	Updated document.

Side effects:
	- The metadata of the document gets updated with the specified values.
	- if Published == yes, document will become visible on the site. No: remove from site.

## admin/analyse 

This gets removed. 

Rationale: Documents get added to a queue for OCR'ing after uploading. OCR'ing happens automatically. 
When OCR'ed succesfully, documents get visible in the /process list.

# Document viewing, 

## GET /document/$docid/#commentid

This shows the document (or the first part), the selected page and the
selected comment.  It is designed to be the full, static URL of the
document with the comment on the page.

It's for static deep-linking. People can post this URL everywhere and
be sure readers get their comment on the document.

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

Add a quote(selection) and comment for the world to see. IE, people can add a selection of a document and their comments.

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
The idea is to provide a way to add a comment to an existing quotaion.
It allows people to discuss a certain quotation.

Parameters:
	- Response-text

Returns:
	- OK/Error

Side Effects:
     - Add a comment to an existing quotation. Sorted by submission date.

