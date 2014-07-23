# Makefile for gomon

# By default, just build the gomon binary into bin/
all: gomon

# Builds and moves the gomon binary
gomon:
	go build cmd/gomon.go
	@mkdir -p bin
	@mv gomon bin/

# Install builds and moves to /usr/local/bin
install: gomon
	cp bin/gomon /usr/local/bin/gomon

# Clean removes everything in the local bin/ directory
clean:
	rm -rf bin/*

# Uninstall removes from /usr/local/bin
uninstall: clean
	rm /usr/local/bin/gomon

config:
	@if [ ! -e /usr/local/etc/gomon.gcfg ]; then \
		cp cmd/gomon.gcfg /usr/local/etc/gomon.gcfg; \
	fi;

# Installs the systemd service, enables it and starts it
install_systemd: install config
	cp cmd/systemd/gomon.service /etc/systemd/system/
	systemctl enable /etc/systemd/system/gomon.service
	systemctl start gomon.service

# Uninstalls the service
uninstall_systemd: uninstall
	systemctl stop gomon.service
	systemctl disable gomon.service
	rm /etc/systemd/system/gomon.service

install_upstart: install config
	cp cmd/upstart/gomon.conf /etc/init/gomon.conf
	start gomon

uninstall_upstart: uninstall
	stop gomon
	rm /etc/init/gomon.conf

# Uninstalls everything and removes the config file
implode: uninstall uninstall_systemd
	rm /usr/local/etc/gomon.gcfg
