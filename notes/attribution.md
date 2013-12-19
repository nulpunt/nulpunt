## Attribution
Nulpunt uses several third-party libraries and packages. We wish to credit those responsible for creating and maintaining these libraries and packages.
This page lists the third-party projects we incorporate. If you believe information is missing or obsolete, please create a pull request or mail contact@nulpunt.nu.

### List of included software
#### Go
##### Go and the Go standard library
 - Source: [http://code.google.com/p/go](http://code.google.com/p/go)
 - Authors: The [Go authors](http://code.google.com/p/go/source/browse/AUTHORS)
 - License: [BSD](http://code.google.com/p/go/source/browse/LICENSE) and a [patent grant](http://code.google.com/p/go/source/browse/PATENTS)

##### go.crypto
 - Source: [http://code.google.com/p/go.crypto](http://code.google.com/p/go.crypto)
 - Authors: The [Go authors](http://code.google.com/p/go/source/browse/AUTHORS)
 - License: [BSD](http://code.google.com/p/go/source/browse/LICENSE) and a [patent grant](http://code.google.com/p/go/source/browse/PATENTS)
 - Location: [/gopath/src/code.google.com/p/go.crypto](/gopath/src/code.google.com/p/go.crypto)

##### mgo, mgo/bson, mgo/txn
 - Project: [http://labix.org/mgo](http://labix.org/mgo)
 - Source: TODO
 - Author: Gustavo Niemeyer <gustavo@niemeyer.net>
 - License: Simplified BSD License (TODO: link)
 - Location: [/gopath/src/labix.org/v2/mgo](/gopath/src/labix.org/v2/mgo)

##### go.leptonica
 - Source: [https://github.com/geertjohan/go.leptonica](https://github.com/geertjohan/go.leptonica)
 - Author: Geert-Johan Riemer
 - License: [BSD-style License](https://github.com/GeertJohan/go.leptonica/blob/master/LICENSE)
 - Location: [/gopath/src/github.com/GeertJohan/go.leptonica](/gopath/src/github.com/GeertJohan/go.leptonica)

##### go.tesseract
 - Source: [https://github.com/geertjohan/go.tesseract](https://github.com/geertjohan/go.tesseract)
 - Author: Geert-Johan Riemer
 - License: [BSD-style License](https://github.com/GeertJohan/go.tesseract/blob/master/LICENSE)
 - Location: [/gopath/src/github.com/GeertJohan/go.tesseract](/gopath/src/github.com/GeertJohan/go.tesseract)

##### go-spew
 - Source: [https://github.com/davecgh/go-spew](https://github.com/davecgh/go-spew)
 - Author: Dave Collins
 - License: [License](https://github.com/davecgh/go-spew/blob/master/LICENSE)
 - Location: [/gopath/src/github.com/davecgh/go-spew](/gopath/src/github.com/davecgh/go-spew)

##### gorilla/context
 - Project: [http://www.gorillatoolkit.org/](http://www.gorillatoolkit.org/)
 - Source: [https://github.com/gorilla/context](https://github.com/gorilla/context)
 - Author: Rodrigo Moraes
 - License: [BSD-style License](https://github.com/gorilla/context/blob/master/LICENSE)
 - Location: [/gopath/src/github.com/gorilla/context](/gopath/src/github.com/gorilla/context)

##### gorilla/mux
 - Project: [http://www.gorillatoolkit.org/](http://www.gorillatoolkit.org/)
 - Source: [https://github.com/gorilla/mux](https://github.com/gorilla/mux)
 - Author: Rodrigo Moraes
 - License: [BSD-style License](https://github.com/gorilla/mux/blob/master/LICENSE)
 - Location: [/gopath/src/github.com/gorilla/mux](/gopath/src/github.com/gorilla/mux)

##### go-flags
 - Source: [https://github.com/jessevdk/go-flags](https://github.com/jessevdk/go-flags)
 - Author: Jesse van den Kieboom
 - License: [BSD-style License](https://github.com/jessevdk/go-flags/blob/master/LICENSE)
 - Location: [/gopath/src/github.com/jessevdk/go-flags](/gopath/src/github.com/jessevdk/go-flags)

#### HTML/CSS/Javascript
##### AngularJS
 - Project: [http://angularjs.org/](http://angularjs.org/)
 - Source: [https://github.com/angular/angular.js](https://github.com/angular/angular.js)
 - Copyright: Google, Inc.
 - License: [MIT License](https://github.com/angular/angular.js/blob/master/LICENSE)
 - Location: Several files in [/http-files/js](/http-files/js)

##### angular.ui bootstrap
 - Project: [http://angular-ui.github.io/bootstrap/](http://angular-ui.github.io/bootstrap/)
 - Source: [https://github.com/angular-ui/bootstrap](https://github.com/angular-ui/bootstrap)
 - Authors: [AngularUI Team](https://github.com/organizations/angular-ui/teams/291112)
 - License: [MIT License](https://github.com/angular-ui/bootstrap/blob/master/LICENSE)
 - Location: Several files in [/http-files/js](/http-files/js) and [/http-files/css](/http-files/css)

##### bootstrap

##### underscore.js

##### jQuery

##### ngStorage

##### checklist-model

##### bootstrap-datepicker

##### angular-file-upload

#### Other
Other types package/libraries that are being used

##### pdftoppm
While not incorporated in the nulpunt source code, this `pdftoppm` is used to convert pdf files into images.
 - Project: TODO
 - Source: TODO
 - Authors: TODO
 - License: TODO
 - Location: Must be installed on systems running `npanalyse`.

##### libleptonica
Used by `go.leptonica` for `npanalyse`.

##### libtesseract3
Used by `go.tesseract` for `npanalyse`.

### Internal notes
When including third party work, please check wether it's license is compatible with our license (AGPL v3).
If work includes source files, please check that the file states copyright or author or project name / web url. If this is absent, please add the following:
```
/*
This file was retrieved from: <url>
The original author is: <person and/or company and/or project>
*/
```