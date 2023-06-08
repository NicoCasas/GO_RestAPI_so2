##
# You should look at the following URL's in order to grasp a solid understanding
# of Nginx configuration files in order to fully unleash the power of Nginx.
# https://www.nginx.com/resources/wiki/start/
# https://www.nginx.com/resources/wiki/start/topics/tutorials/config_pitfalls/
# https://wiki.debian.org/Nginx/DirectoryStructure
#
# In most cases, administrators will remove this file from sites-enabled/ and
# leave it as reference inside of sites-available where it will continue to be
# updated by the nginx packaging team.
#
# This file will automatically load configuration files provided by other
# applications, such as Drupal or Wordpress. These applications will be made
# available underneath a path with that package name, such as /drupal8.
#
# Please see /usr/share/doc/nginx-doc/examples/ for more detailed examples.
##

# Default server configuration
#
server {
	listen [::]:80;
	listen 80;
	# Add index.php to the list if you are using PHP
	#index index.html index.htm index.nginx-debian.html;

	server_name dashboard.com;

	location /{
		return 404;
	}
	
	location /api/users/login{
		proxy_pass http://localhost:8030/users/login;
	}

	location /api/users/createuser{
		proxy_pass http://localhost:8030/api/users/createuser;
	}

	location /api/users/listall{
		proxy_pass http://localhost:8030/api/users/listall;
	}

	location /api/users/ping{
		proxy_pass http://localhost:8030/api/users/ping;
	}

}


# Virtual Host configuration for example.com
#
# You can move that to a different file under sites-available/ and symlink that
# to sites-enabled/ to enable it.
#
#server {
#	listen 80;
#	listen [::]:80;
#
#	server_name example.com;
#
#	root /var/www/example.com;
#	index index.html;
#
#	location / {
#		try_files $uri $uri/ =404;
#	}
#}