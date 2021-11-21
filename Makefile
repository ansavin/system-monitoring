all:
	docker build -t system-monitor .
run: all
	docker run \
	--net="host" \
	--pid="host" \
	-v "/:/host:ro,rslave" \
	--name system-monitor \
	system-monitor
clean:
	docker rm system-monitor && docker rmi system-monitor