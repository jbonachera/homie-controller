test:
	docker run --rm -v $$(pwd):/usr/local/src/go/src/github.com/jbonachera/homie-controller $$(docker build -qf Dockerfile.build .) 
