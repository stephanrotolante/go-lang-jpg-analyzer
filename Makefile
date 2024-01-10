run:	
	go run *.go -file cat.jpg -output color.bmp | tee out.log
	
dump:
	hex-dump -f cat.jpg  | tee dump.log

build:
	go build -o image-analyze.bin *.go

display:
	python3 display.py -f color.out

test:
	make run && make display