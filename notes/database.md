accounts
 - `_id` (bson.ObjectId)
 - `handle` (string e.g. "GeertJohan" from "@GeertJohan")
 - `email` (string, optional)
 - `avatar` (to be decided)

documents
technical parameters (for system)
 - `_id` (bson.ObjectId)
 - `original` (string, refers to location of the orginal document in GridFS)
 - `published` boolean; true: document is visible for users; false: new or not yet processed document
 - `uploaded_date` (time.Time); date of *publication* on Nulpunt.
content parameters (for people)
 - `uploader` (bson.ObjectId, refers to `accounts._id`)
 - `title` (string)
 - `summary` (string)
 - `source` (string)
 - `categories` ([]string)  // These come from the Tags-table
 - `originalDate` (time.Time)  // Time of publishing by the gov-ment agency or date of FOI-response.

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
 - `text` string; the text in the same order as the lines-attribute, use for search/sharing. Contains ocr-errors

char-object (inside page):
 - `x1` (int, left) in pixels
 - `y1` (int, top) in pixels
 - `x2` (int, bottom) in pixels
 - `y2` (int, right) in pixels
 - `c` (string, character)

images
 - `_id` (bson.ObjectId)
 - `documentId` (bson.ObjectId, refers to `documents._id`)
 - `pageNr` (int, page number)
 - `image-location` (string) address in gridfs of the image data

annotations
 - `_id` (bson.ObjectId)
 - `annotatorId` (bson.ObjectId, refers to `accounts._id`)
 - `createDate` (time.Time)
 - `annotation` (string)
 - `location` ([]object) // In future, there could be multiple sections in a single annotation.
    - `page` (int)
    - `x1` (int))
    - `y1` (int)
    - `x2` (int)
    - `y2` (int)

comments
 - `_id` (bson.ObjectId)
 - `documentID` (bson.ObjectId) refers to document
 - `annotationID` (bson.ObjectId) refers to annotation
 - `commenterId` (bson.ObjectId, refers to `accounts._id`)
 - `createDate` (time.Time)
 - `comment` (string)
 - `parentID` (bson.ObjectId) refers to comment

uploads
 - `_id` (bson.ObjectId)
 - `uploaderId` (bson.ObjectId)
 - `original` 
 - `filename` (string)
 - `uploadDate` (time.Time)
 - `language` (string); language of the document to help the OCR (default 'nld')
