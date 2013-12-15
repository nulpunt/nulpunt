accounts
 - `_id` (bson.ObjectId)
 - `handle` (string e.g. "GeertJohan" from "@GeertJohan")
 - `email` (string, optional)
 - 

documents
technical parameters (for system)
 - `_id` (bson.ObjectId)
 - `original` (string, refers to location in GridFS)
 - `published` boolean; true: document is visible for users; false: new or not yet processed document
content parameters (for people)
 - `accountId` (bson.ObjectId, refers to `accounts._id`)
 - `title` (string)
 - `summary` (string)
 - `source` (string)
 - `categories` ([]string)  // These come from the Tags-table
 - `publicationDate` (time.Time)  // Time that it gets published on Nulpunt.


tags
 - `_id` (bson.ObjectId)
 - `tag` (string)

Note: tags have an ObjectId, these are not for referencing in other tables.
Just insert the tag-string into other tables where needed.

pages
 - `_id` (bson.ObjectId)
 - `documentId` (bson.ObjectId, refers to `documents._id`)
 - `pageNr` (int, page number)
 - `lines` ([][]char-object)
 - `image` []byte; the png image data of the page
 - `text' []string; the text in the same order as the lines-attribute

char-object (inside page):
 - `x` (int, left)
 - `y` (int, top)
 - `s` (int, size in pixels)
 - `c` (string, character)

annotations
 - `_id` (bson.ObjectId)
 - `accountId` (bson.ObjectId, refers to `accounts._id`)
 - `createDate` (time.Time)
 - `location` (object)
  - `page` (int)
  - `start` (object)
    - `x` (int)
    - `y` (int)
  - `end` (object)
    - `x` (int)
    - `y` (int)

uploads
 - `_id` (bson.ObjectId)
 - `uploaderId` (bson.ObjectId)
 - `original` 
 - `filename` (string)
 - `uploadDate` (time.Time)
 - `language` (string); language of the document to help the OCR (default 'nld')