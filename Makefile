daemon_bin_path=/usr/sbin/piremd
service_path=/lib/systemd/system/piremd.service

configure:
	ls $(daemon_bin_path)
	ls $(service_path)

piremd: configure
	go build -o piremd .

install: piremd piremd.service
	cp piremd $(daemon_bin_path)
	cp piremd.service $(service_path)

uninstall: install
	rm piremd $(daemon_bin_path)
	rm $(service_path)
