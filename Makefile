run:	
	go run *.go -file profile.jpg | tee out.log
	
dump:
	hex-dump -f profile.jpg  | tee dump.log

build:
	go build -o image-analyze.bin *.go

display:
	python3 display.py

test:
	make run && make display