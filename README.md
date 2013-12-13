0.
================================

This repository contains the server and the static web files (html,css,js) for Nulpunt

For more information about Nulpunt, please visit [nulpunt.nu](http://nulpunt.nu)

### Quickstart
1. [Install go](http://golang.org/doc/install/)
2. Fork the repository on GitHub
3. Clone to local machine: `git clone git@github.com:YOUR-USERNAME/nulpunt.git`
4. Execute: `cd nulpunt`
5. Execute: `GOPATH=$(PWD)/gopath go build npserver`
6. Run npserver: `./npserver`

For changes to go code, you must recompile and restart the server (steps 5 and 6). Changes to html/css/js only need browser refresh.

### Closed repo
The development of this application is closed fow now. If you know someone that wants to join, please ask an Owner ([GeertJohan](mailto:gjr19912@gmail.com)) to add this person to the "Contributors" team.

### Development
Please view the [issues](https://github.com/nulpunt/nulpunt/issues?state=open) on this repo. If you have an idea or suggestion, please [create a new issue](https://github.com/nulpunt/nulpunt/issues/new).

If you want to create something new, make sure the idea was approved by a nulpunt maintainer. It would be sad to have development time spilled on duplicates or features that won't be merged in.

### Dependencies
This project uses several third-party dependencies. Such as the `mgo` driver for MongoDB.
These dependencies (third-party packages) are to be imported by nulpunt code with their fully qualified import name (e.g. `labix.org/v2/mgo`).
We are keeping the source for imported packages within this repository for several reasons:
- A commit can always refer to the right version of a third-party package, because it is included in the commit.
- New third-party code must go through a PR, and can easily be checked.
- Project will still build when remote dependency is unreachable or removed.

Cons:
- Need to set GOPATH environment variable to make this work properly.

You can permanently set the GOPATH for this project in your `.profile` file or `.bashrc` file.

### How to contribute
1. Fork this repository on GitHub and clone to local.
2. Create a new branch and start developing
3. Make sure that the code is formated according to `go fmt`.
4. Push your branch+changes to github and create a pull request.
5. Pull request is automatically built by Jenkins.
6. When PR is approved, it is merged into the master branch.
7. Repeat from step 2 for each bugfix/feature.

### CI
We have [jenkins](https://ci.nulpunt.nu)!

Jenkins performs two tasks:
- Run build and tests for each new PR (and new commits in that PR), then report status back to Github.
- Run nightly build and restart nightly when repository has changed (PR merged in).

