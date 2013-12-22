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
 - `original` (string, refers to location of the orginal document in GridFS)
 - `published` boolean; true: document is visible for users; false: new or not yet processed document
 - `upload_date` (time.Time); date of *publication* on Nulpunt.
content parameters (for people)
 - `uploader` (string, refers to `accounts.username`)
 - `title` (string)
 - `summary` (string)
 - `source` (string)
 - `country` (string)
 - `FOIRequester` (string) // Wobber
 - `FOIARequest` (string) // Wob-verzoek
 - `category (string) // "Kamerbrief", "Rapport", ...
 - `language` (string) // same value as in upload-table.
 - `tags` ([]string)  // These come from the Tags-table
 - `originalDate` (time.Time)  // Time of publishing by the gov-ment agency or date of FOI-response.

### tags
 - `_id` (bson.ObjectId)
 - `tag` (string)

Note: tags have an ObjectId, these are not for referencing in other collections.
Just insert the tag-string into other collections where needed.

### pages
 - `_id` (bson.ObjectId)
 - `documentId` (bson.ObjectId, refers to `documents._id`)
 - `pageNr` (int, page number)
 - `lines` ([][]char-object)
 - `text` (string); the text in the same order as the lines-attribute, use for search/sharing. Contains ocr-errors

#### char-object
 - `x1` (int, left) in pixels
 - `y1` (int, top) in pixels
 - `x2` (int, bottom) in pixels
 - `y2` (int, right) in pixels
 - `c` (string, character)

### annotations
 - `_id` (bson.ObjectId)
 - `documentId` (bson.ObjectId, refers to Documents)
 - `annotator` (string)
 - `createDate` (time.Time)
 - `annotation` (string)
 - `comments` (comment)
 - `location` ([]object) // In future, there could be multiple sections in a single annotation.
    - `page` (int)
    - `x1` (int))
    - `y1` (int)
    - `x2` (int)
    - `y2` (int)

#### comment
 - `_id` (bson.ObjectId) // needed to do treewalking to get new comments in the right place
 - `commenter` (string, refers to `accounts.handle`)
 - `createDate` (time.Time)
 - `comment` (string)
 - `comments` ([]comment) *recursion, disabled for first version??*

### uploads
 - `_id` (bson.ObjectId)
 - `uploader` (string, refers to `accounts.username`)
 - `original` (string); reference to the original pdf file.
 - `filename` (string)
 - `uploadDate` (time.Time)
 - `language` (string); language of the document to help the OCR (default 'nl_NL')

## GridFS
We're using GridFS to store files.

### uploads
Filename must be formatted as: `upload/<uploader-handle>/<unix-timestamp>-<random-string-10-chars>-<original-filename>`
Holds original uploaded file.

### images
Filename must be formated as: `pages/<documentId>-<pageNumber>.png`
