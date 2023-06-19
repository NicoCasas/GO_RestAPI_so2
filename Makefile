.ONESHELL:

default: build

build:
	cd users_service/
	go build

	cd ..

	cd processing_service/
	go build

	cd ..
	

install:
	# Referido a user_service // dashboard.com
	
	cd users_service/
	
	## Servicio
	sed -in 's@^WorkingDirectory=.*@WorkingDirectory=${PWD}/@ ; s@^ExecStart=.*@ExecStart=${PWD}/users_service@' ./conf/users_service.service
	rm ./conf/users_service.servicen 
	sudo cp ./conf/users_service.service /etc/systemd/system/users_service.service
	
	## Nginx
	sudo cp ./conf/dashboard.com /etc/nginx/sites-available/dashboard.com
	sudo ln -s /etc/nginx/sites-available/dashboard.com /etc/nginx/sites-enabled/
	sudo cp ./conf/.htpasswd /etc/nginx/.htpasswd

	## Sudoers
	sudo useradd sistemas_operativos -s /usr/sbin/nologin
	sudo cp ./conf/sistemas_operativos	/etc/sudoers.d/ 

	## Sshd
	#sudo cp ./conf/lab6_operativos.conf /etc/ssh/sshd_config.d/
	sudo echo "AllowGroups operativos_ssh_clients" >> /etc/ssh/sshd_config.d/operativos.conf

	## Hosts
	sudo echo -e "127.0.0.1\tdashboard.com" >> /etc/hosts

	#
	cd ..

	# Referido a processing_service // sensors.com
	cd processing_service

	## Servicio
	sed -in 's@^WorkingDirectory=.*@WorkingDirectory=${PWD}/@ ; s@^ExecStart=.*@ExecStart=${PWD}/processing_service@' ./conf/processing_service.service 
	rm ./conf/processing_service.servicen
	sudo cp ./conf/processing_service.service /etc/systemd/system/processing_service.service

	## Nginx
	sudo cp ./conf/sensors.com /etc/nginx/sites-available/sensors.com
	sudo ln -s /etc/nginx/sites-available/sensors.com /etc/nginx/sites-enabled/

	## Hosts
	sudo echo -e "127.0.0.1\tsensors.com" >> /etc/hosts

	#
	cd ..

	# Nginx - default_server
	sudo cp ./default_server /etc/nginx/sites-available/
	sudo ln -s /etc/nginx/sites-available/default_server /etc/nginx/sites-enabled/default_server