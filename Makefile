daemon_bin_path=/usr/sbin/piremd
service_path=/lib/systemd/system/piremd.service
config_path=/etc/piremd

configure:
	@ls /dev/null
	@ls /usr/sbin > /dev/null
	@ls /lib/systemd/system > /dev/null
	@ls /etc > /dev/null

piremd: configure
	go build -o piremd .
	chmod 711 piremd

install: piremd daemon_install_files/piremd.service
	sudo cp piremd $(daemon_bin_path)
	rm piremd
	sudo cp daemon_install_files/piremd.service $(service_path)
	-mkdir $(config_path)
	sudo cp daemon_install_files/config.json $(config_path)
	-systemctl daemon-reload

uninstall: install
	sudo rm $(daemon_bin_path)
	sudo rm $(service_path)
	sudo rm -rf $(config_path)
	-systemctl daemon-reload
