0.
================================

This repository contains the server and the static web files (html,css,js) for Nulpunt

For more information about Nulpunt, please visit [nulpunt.nu](http://nulpunt.nu)

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
- Need to set GOPATH to get this to work.

If anyone can point to or provide a better way to do this, please open an issue.

### Building
Clone this repository (or, preferably, a fork):
`git clone git@github.com:nulpunt/nulpunt.git`

Changedir into the nulpunt directory
`cd nulpunt`

Invoke go build with a specific GOPATH:
`GOPATH=$(PWD)/gopath go build`

You can permanently set the GOPATH for this project in your `.profile` file or `.bashrc` file.

### How to contribute
1. Fork this repository on GitHub  
2. Edit your fork (preferably use a new branch for each feature/bugfix)
3. Send pull request
4. ????
5. Profit!

### CI
We have jenkins!
https://ci.nulpunt.nu
Jenkins performs two tasks:
- Run build and tests for each new PR (and new commits in that PR), then report status back to Github.
- Run nightly build and restart nightly when repository has changed (PR merged in).

