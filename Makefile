daemon_bin_path=/usr/sbin/piremd
service_path=/lib/systemd/system/piremd.service

configure:
	@ls /usr/sbin > /dev/null
	@ls /lib/systemd/system > /dev/null

piremd: configure
	go build -o piremd .
	chmod 711 piremd

install: piremd daemon_install_files/piremd.service
	cp piremd $(daemon_bin_path)
	cp daemon_install_files/piremd.service $(service_path)
	sudo systemctl daemon-reload
	rm piremd

uninstall:
	rm $(daemon_bin_path)
	rm $(service_path)
	sudo systemctl daemon-reload
