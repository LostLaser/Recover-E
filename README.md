# election
[![PkgGoDev](https://pkg.go.dev/badge/github.com/LostLaser/election)](https://pkg.go.dev/github.com/LostLaser/election)
[![Go Report Card](https://goreportcard.com/badge/github.com/LostLaser/election)](https://goreportcard.com/report/github.com/LostLaser/election)
[![Release](https://img.shields.io/github/release/LostLaser/election.svg?style=flat-square)](https://github.com/LostLaser/election/releases/latest)
[![codecov](https://codecov.io/gh/LostLaser/election/branch/master/graph/badge.svg?token=Y0ZFB9MLTZ)](https://codecov.io/gh/LostLaser/election)

## Overview
election is a project that simulates various leader election algorithms in distributed systems. This project's primary use is for educational purposes.

Visual implementation located here: https://recover-e.herokuapp.com

### To install
```
go get github.com/LostLaser/election
```

### Algorithms currently implemented
* [Bully election](https://en.wikipedia.org/wiki/Bully_algorithm)
* [Ring election](https://en.wikipedia.org/wiki/Leader_election#Leader_election_in_rings)


## Contributing
Want to improve the project? Open an issue and list what you would like to change/add so we can discuss it. 

### Retrieving the source
```
git clone https://github.com/LostLaser/election
cd election
go test -v
```
