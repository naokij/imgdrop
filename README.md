## ImageDrop
A little one-click image sharing web app with user management. Written in Golang.

### Build
glide is used to manage package. Please [install](https://github.com/Masterminds/glide) it first

	cd $GOPATH
	go get github.com/naokij/imgdrop
	go get glide
	cd src/github.com/naokij/imgdrop
	glide install
	go build
	cp src/github.com/naokij/imgdrop/data.db.sample src/github.com/naokij/imgdrop/data.db
	cp src/github.com/naokij/imgdrop/conf/app.conf.sample src/github.com/naokij/imgdrop/conf/app.conf

### Run

	cd src/github.com/naokij/imgdrop
	mkdir log
	./imgdrop

* open http://localhost:8080 in your browser
* username: admin
* password: p@sswd