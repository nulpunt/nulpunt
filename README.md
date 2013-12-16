0.
================================

This repository contains the server and the static web files (html,css,js) for Nulpunt

For more information about Nulpunt, please visit [nulpunt.nu](http://nulpunt.nu)

### Quickstart
1. [Install go](http://golang.org/doc/install/)
2. [Install MongoDB](http://www.mongodb.org/downloads)
3. Fork the repository on GitHub
4. Clone to local machine: `git clone git@github.com:YOUR-USERNAME/nulpunt.git`
5. Execute: `cd nulpunt`
6. Execute: `GOPATH=$(PWD)/gopath go build npserver`
7. Run npserver: `./npserver`

For changes to go code, you must recompile and restart the server (steps 5 and 6). Changes to html/css/js only need browser refresh.

### Closed repo
The development of this application is closed fow now. If you know someone that wants to join, please ask an Owner ([GeertJohan](mailto:gjr19912@gmail.com)) to add this person to the "Contributors" team.

### Development

Nulpunt consists of a seperate server and client web-application. The client web-application is a standalone SPA (Single Page Application).
The server exposes a set of [services](notes/server-api.md) to the client. The server uses MongoDB as a database, it's outline/structure is defined [here](notes/database.md).

Please view the [issues](https://github.com/nulpunt/nulpunt/issues?state=open) on this repo. If you have an idea or suggestion, please [create a new issue](https://github.com/nulpunt/nulpunt/issues/new).

#### Server development

##### Go
The server is written using the go programming language. For more information, visit [golang.org](http://golang.org).

##### Go dependencies (packages/libraries)
This project uses several third-party dependencies. Such as the `mgo` driver for MongoDB.
These dependencies (third-party packages) are to be imported by nulpunt code with their fully qualified import name (e.g. `labix.org/v2/mgo`).
We are keeping the source for imported packages within this repository for several reasons:
- A commit can always refer to the right version of a third-party package, because it is included in the commit.
- New third-party code must go through a PR, and can easily be checked.
- Project will still build when remote dependency is unreachable or removed.

You can permanently set the GOPATH for this project in your `.profile` file or `.bashrc` file.

#### Client development
The client, or "front-end", is written using HTML, CSS and Javascript. The client uses several existing projects/libraries to make things easier.
 - [jQuery](http://jquery.com)
 - [Bootstrap](http://getbootstrap.com)
 - [AngularJS](http://angularjs.org), and some angular modules
 - [Underscore](http://underscorejs.org)
 - [CrytoJS](https://crypto-js.googlecode.com)

##### AngularJS
It is important to understand how [AngularJS](http://angularjs.org) works because this is the foundation for the nulpunt client application. If you have not worked with AngularJS yet, please folow some [basic tutorials (scroll down)](http://egghead.io/lessons), it's very easy to pick up.

#### How to contribute
1. Fork this repository on GitHub and clone to local.
2. Create a new branch and start developing
3. Make sure that the code is formated according to `go fmt`.
4. Push your branch+changes to github and create a pull request.
5. Pull request is automatically built by Jenkins.
6. When PR is approved, it is merged into the master branch.
7. Repeat from step 2 for each bugfix/feature.

#### CI
We have [jenkins](https://ci.nulpunt.nu)!

Jenkins performs two tasks:
- Run build and tests for each new PR (and new commits in that PR), then report status back to Github.
- Run nightly build and restart nightly when repository has changed (PR merged in).

#### OCR process
The quickstart and server instructions above do not include the OCR process (`npocr`).
To get `npocr` up and running, perform the following:
1. Install go.leptonica dependencies as explained [here](https://github.com/GeertJohan/go.leptonica)
2. Install go.tesseract dependencies as explained [here](https://github.com/GeertJohan/go.tesseract)
3. Install and run `nsqlookupd` and `nsqd` with their defaults (localhost): [follow this quick start](http://bitly.github.io/nsq/overview/quick_start.html).
4. Change dir into the root of your nulpunt repository clone
5. Build npocr: `GOPATH=$(pwd)/gopath go build npocr`
6. Run npocr: `./npocr`