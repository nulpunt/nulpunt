## Collections

### accounts
 - `_id` (bson.ObjectId)
 - `username` (string,) primary key. ( e.g. "GeertJohan" in "@GeertJohan", **indexed**)
 - `email` (string, optional)
 - `avatar` (to be decided, link to GridFS file?)
 - `admin` (bool) whether user is administrator or ordinary user 

### profile
 - `_id` (bson.ObjectId)
 - `username` primary key, refers to account.username
 - `tags` ([]string) list of tags this user is interested in

### document
technical parameters (for system)
 - `_id` (bson.ObjectId)

 - `uploaderUsername` (string, refers to `accounts.username`)
 - `uploadFilename` (string) // Filename of the original upload
 - `uploadGridFilename` (string, refers to location of the orginal document in GridFS)
 - `uploadDate` (time.Time); date of *publication* on Nulpunt.
content parameters (for people)
 - `language` (string) // language of the document to help the OCR (default 'nl_NL')
 - `pageCount` (int) // number of pages
 - `analyseState` (string) options("uploaded", "started", "completed", "error")

 - `title` (string)
 - `summary` (string)
 - `category (string) // "Kamerbrief", "Rapport", ...
 - `tags` ([]string)  // These come from the Tags-table
 
 - `FOIRequester` (string) // Wobber
 - `FOIARequest` (string) // Wob-verzoek
 - `originalDate` (time.Time)  // Time of publishing by the gov-ment agency or date of FOI-response.
 - `source` (string) // "NL - Binnenlandse zaken", "EN - Foreign affairs", "US - Foreign affairs"
 - `country` (string) // "NL", "EN"
 
 - `published` boolean; true: document is visible for users; false: new or not yet processed document
 
 - `hits` (int)

### tags
 - `_id` (bson.ObjectId)
 - `tag` (string)

Note: tags have an ObjectId, these are not for referencing in other collections.
Just insert the tag-string into other collections where needed.

### pages
 - `_id` (bson.ObjectId)
 - `documentId` (bson.ObjectId, refers to `documents._id`)
 - `pageNumber` (int, page number)
 - `lines` ([][]char-object)
 - `text` (string); the text in the same order as the lines-attribute, use for search/sharing. Contains ocr-errors
 - `highresWidth` (int) The width (in pixels) for the highres(900dpi) render.
 - `highresHeight` (int) The height (in pixels) for the highres(900dpi) render.

#### char-object
 - `x1` (float32, left) in percentage relative to the image
 - `y1` (float32, top) in percentage relative to the image
 - `x2` (float32, bottom) in percentage relative to the image
 - `y2` (float32, right) in percentage relative to the image
 - `c` (string, character)

### annotations
 - `_id` (bson.ObjectId)
 - `documentId` (bson.ObjectId, refers to Documents)
 - `annotator` (string)
 - `createDate` (time.Time)
 - `annotation` (string)
 - `comments` ([]comment)
 - `location` ([]object) // In future, there could be multiple sections in a single annotation.
    - `pageNumber` (int) index
    - `y1` (float32, left) in percentage relative to the image
    - `x1` (float32, top) in percentage relative to the image
    - `x2` (float32, bottom) in percentage relative to the image
    - `y2` (float32, right) in percentage relative to the image

#### comment
 - `_id` (bson.ObjectId) // needed to do treewalking to get new comments in the right place
 - `commenterUsername` (string, refers to `accounts.username`)
 - `createDate` (time.Time)
 - `commentText` (string)
 - `comments` ([]comment) *recursion, disabled for first version??* (JANUARI/FEBRUARI)

## GridFS
We're using GridFS to store files.

### uploads
Filename must be formatted as: `uploads/<uploader-username>-<unix-timestamp>-<random-string-10-chars>-<original-filename>`
Holds original uploaded file.

### highres
Filename must be formatted as: `highres/<docIdHex>-<pageNumber>.png`
Holds png for each page for given document rendered at 600dpi from pdf file

### docviewer-pages
Filename must be formated as: `docviewer-pages/<docIdHex>-<pageNumber>.png`
Holds png for each page for any given document resized to a width of 1000 px
