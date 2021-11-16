all:
	docker build -t system-monitor .
run: all
	docker run system-monitor
clean:
	docker rm `docker ps -q -a -f ancestor=system-monitor` && docker rmi system-monitor