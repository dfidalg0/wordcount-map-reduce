all: wordcount

wordcount: bin/wordcount

bin/wordcount: wordcount/*.go mapreduce/*.go | bin
	cd ./wordcount && go build -o ../bin/ && cd ..

bin:
	mkdir bin

clean:
	rm -rf bin/ map/ reduce/ result/

reset:
	rm -rf map/ reduce/ result/
