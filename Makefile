daemon_bin_path=/usr/sbin/piremd
service_path=/lib/systemd/system/piremd.service

configure:
	@ls /usr/sbin > /dev/null
	@ls /lib/systemd/system > /dev/null

piremd: configure
	go build -o piremd .

install: piremd piremd.service
	cp piremd $(daemon_bin_path)
	cp piremd.service $(service_path)

uninstall: install
	rm piremd $(daemon_bin_path)
	rm $(service_path)
