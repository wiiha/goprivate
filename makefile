
all: webfront
	@echo "DONE"

.PHONY: webfront
webfront:
	./compile-webfront.sh

clean:
	rm -rf webfront/dist/
	rm -rf server/webfront