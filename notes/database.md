accounts
 - `_id` (bson.ObjectId)
 - `handle` (string e.g. "GeertJohan" from "@GeertJohan")
 - `email` (string, optional)
 - 

documents
 - `_id` (bson.ObjectId)
 - `accountId` (bson.ObjectId, refers to `accounts._id`)
 - `title` (string)
 - `summary` (string)
 - `source` (string)
 - `categories` ([]string)
 - `publicationDate` (time.Time)
 - `original` (string, refers to location in GridFS)

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