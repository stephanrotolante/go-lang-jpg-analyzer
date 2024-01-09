run:	
	go run *.go -file cat.jpg -output cat.color.out | tee cat.out.log
	
dump:
	hex-dump -f cat.jpg  | tee cat.dump.log

build:
	go build -o image-analyze.bin *.go

display:
	python3 display.py -f cat.color.out

test:
	make run && make display