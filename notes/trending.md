The main page shows a list of trending documents. This is likely to be the most-loaded resource available by npserver.
Because this resource changes only once in 5-20 (?) minutes, we can easily optimizte this.

npserver provides a serivce /service/trending, which doesn't take any parameters and simply returns a json structure holding an array with 'trending annotation objects', each containing the following fields:
 - cropped image (selected text from document) in base64 encoding
 - annotation
 - annotator (user handle)
 - ++ more fields (?)

This json structure can be used by the the frond-end without the server performing any modifications.

A seperate process, `nptrending`, calculates trending annotations and generates the json structure described above.
This structure is saved on [etcd](https://github.com/coreos/etcd) with key "/trending".
