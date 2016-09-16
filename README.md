## ImageDrop
A little one-click image sharing web app with user management

### Build


	cd $GOPATH
	go get github.com/naokij/imgdrop
	go get glide
	cd src/github.com/naokij/imgdrop
	glide update
	go build

### Run

	cd src/github.com/naokij/imgdrop
	./imgdrop

* open http://localhost:8080 in your browser
* username: admin
* password: p@sswd