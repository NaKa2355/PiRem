bin=bin/piremd
daemon_bin_path=/usr/sbin/piremd
config_path=/etc/piremd
service_path=/lib/systemd/system/piremd.service
plugins_path=/opt/piremd

build:
	go build -o $(bin) .
	@echo "build completed!"

install: bin/piremd
	cp $^ $(daemon_bin_path) #copy daemon binary to right place
	cp daemon_install_files/piremd.service $(service_path) #copy systemd service file to right place

	mkdir $(config_path) #make a directory that contains a config file
	cp daemon_install_files/config.json $(config_path) #copy the default config file to the directory
	
	mkdir $(plugins_path) #make a directory that contains plugin binaries

	systemctl daemon-reload #reload daemon to recognize the service file

	rm $(bin)
	@echo "install completed!"

uninstall:
	-rm $(daemon_bin_path)

	-rm $(service_path)

	-rm -rf $(config_path)

	-rm -rf $(plugins_path)

	systemctl daemon-reload

	@echo "uninstall completed..."

update: bin/piremd
	cp $^ $(daemon_bin_path) #copy daemon binary to right place
	rm $(bin)
	@echo "update completed!"

all: build