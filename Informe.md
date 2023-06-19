# Informe Lab 6: Sistemas Embebidos / Webservices - Desarrollo

## Script de instalación

Para buildear los archivos y luego configurar los archivos necesarios para el funcionamiento de los servicios y nginx es necesario usar los siguientes comandos situado en la carpeta root del proyecto, donde se encuentra el Makefile:

    make build
    sudo -E make install	   # -E para preservar la variable de entorno PWD

Luego, iniciar nginx via:

    sudo systemctl start nginx     #(si no estaba corriendo)
    sudo nginx -s reload           #(si ya se estaba ejecutando)

Finalmente, levantamos los servicios con:
    
    sudo systemctl start users_service.service
    sudo systemctl start processing_service.service
    
## Diseño de los endpoints

Ambos endpoints siguen el patron de diseño de Controller - Service - Repository, al igual que el [ejemplo](https://github.com/SOiI-UNC/go-example) presentado en clases.

## Initializers
El package initializers tiene funciones que se ejecutan en `init`, es decir, una vez antes de ejecutar el main.
Estas:
- Cargan como variables de entorno aquellas encontradas en el archivo de configuración conf/.env, usando el package [godotenv](https://pkg.go.dev/github.com/joho/godotenv)
- Crear, en caso de no existir, el grupo en el que se van a encontrar los usuarios con permisos de conexión por ssh. 

## Servicios - Systemd
Los servicios, tanto el de usuarios como el de procesamiento, son manejados por systemd. Para esto hay que crear, por servicio, un archivo en `/etc/systemd/system`. En el proyecto, se encuentran en `'service'/conf/'service'.service`.
Las directivas más importantes son:

    ExectStart='path_to_bin'                  -> Path al binario a ejecutar al hacer start al servicio
    WorkingDirectory='path_to_service_file'   -> Configura el directorio de trabajo. En este caso debe estar en el directorio
                                                 donde se encuentra el binario, ya que desde ahí el programa busca el archivo
                                                 de configuracion (variables de entorno)
    User=sistemas_operativos                  -> Usuario bajo el que corren los servicios

El último punto se utiliza porque los servicios van a ejecutar comandos con sudo sin necesidad de contraseña. Ergo, por cuestiones de seguridad, se los encapsula en un usuario creado específicamente para correr estos servicios, con permisos específicos para solo usar con estos requerimientos (sudo y sin contraseña) los comandos `useradd`, `passwd` y `groupadd`.
Para esto es necesario modificar el archivo `/etc/sudoers` o bien incluir uno en `/etc/sudoers.d/`. En este laboratorio, se optó por la segunda opción y el archivo se encuentra en `users_service/conf/sistemas_operativos`. Si bien no es necesario, se eligió hacer una línea por cada comando, que tienen la siguiente forma:

    USR_SISTEMAS_OPERATIVOS ALL = NOPASSWD:CMD_GROUPADD   

Que se lee: El usuario sistemas_operativos (aparece como USR_SISTEMAS_OPERATIVOS porque previamente se definió un alias) tiene permisos en todos los hosts, de ejecutar el comando groupadd sin necesidad de contraseña. 

## Servicio de usuario

### Endpoint de autenticación - POST /api/users/login
Para realizar las peticiones a los endpoints del sistema operativo es necesario estar autorizado. Esta autorización se logra a través de un jwt. Para obtenerlo, es necesario primero hacer un POST a dashboard.com/api/users/login con las credenciales `username` y `password` de un usuario autorizado.

    curl -s --request POST \
              --url http://dashboard.com/api/users/login \
                -u admin:admin \
                --header 'accept: application/json' \
                --header 'content-type: application/json' \
                --data '{"username": "nicoAdmin", "password": "nicoPass"}'

Las credenciales de dicho usuario se encuentran en una base de datos sql. La contraseña hasheada. Para gestionar la bdd, se usa sqlite3. Como bibliografía se usaron [[1]](https://github.com/SOiI-UNC/go-example) [[2]](https://zetcode.com/golang/sqlite3/).
En el jwt, se incluyen dos campos en el payload:
  - iss: Nombre de usuario autenticado
  - exp: Indica el tiempo de expiración de la validez del token. En este caso, el token dura una hora.

### Endpoint de creación de usuario - POST /api/users/createuser
Este endpoint tiene como finalidad crear un usuario en el sistema operativo para que luego pueda acceder al sistema via ssh. Antes de realizar cuaquier acción, primero verifica la validez del token. En caso de no ser correcto, se retorna `401`

#### Creación del usuario
Para esta parte se usó el package [os/exec](https://pkg.go.dev/os/exec) siguiendo los ejemplos encontrados en [[3]](https://zetcode.com/golang/exec-command/). Estas funciones tienen comportamientos similares tipo fork/exec de c en el sentido de que no invocan una shell. 
Los comandos realizados de esta forma en este endpoint son:

    sudo useradd -g 'group_id' -s '/bin/bash' username     -> Para crear el usuario de nombre username. En el grupo 'group_id' y  que inicie sesion en bash.  
    sudo passwd username                                   -> Para asignarle contraseña al usuario.  

En el caso del último comando, la contraseña se pasa 2 veces vía stdin. Es por esto que se crea un string con las contraseñas separadas por \n y se lo pasa via pipe. 

#### Acceso via ssh
Para este punto, se decidió crear un grupo `operativos_ssh_clients` cuyo group_id, (obtenido del `/etc/group`) es el pasado como argumento en el comando useradd del punto anterior.
Luego, se le da acceso a todos los usuarios del mismo en el archivo de configuración `/etc/ssh/sshd_config` o en un archivo terminado en '*.conf*' en `/etc/ssh/sshd_config.d/` a través de la sentencia

    AllowGroups operativos_ssh_clients

### Endpoint que lista usuarios - GET /api/users/listall
Este endpoint es bastante simple, consta de parcear el archivo `/etc/passwd`, armar un slice y retornarlo en formato json. Antes de realizar cuaquier acción, primero verifica la validez del token. En caso de no ser correcto, se retorna `401`

## Servicio de procesamiento

Los endpoints de este servicio son bastante triviales:

    - /api/processing/submit    -> Incrementa un contador global
    - /api/processing/summary   -> Retorna el valor del contador

## Ngnix

Para exponer los servicios, se usa nginx como proxy inverso. Funciona como un especie de pasamanos, donde a partir de las sentencias:
    
    server_name sensors.com
    
    location /api/processing/submit {
		proxy_pass http://localhost:8080/api/processing/submit;
	}

al ingresar a `sensors.com/api/processing/submit`, nginx redirecciona la request a `http://localhost:8080/api/processing/submit`.
Los archivos de configuración de nginx se encuentran en el proyecto en la carpeta `'service'/conf/` y se deben ubicar en `/etc/nginx/sites-available`. Luego, poner un link simbólico en `/etc/nginx/sites-enabled` apuntando al archivo creado anteriormente. Esto podría omitirse poniendo el archivo directamente en el último path mencionado, pero se hace así por buena práctica, para activar o desactivar los sitios, simplemente creando o borrando los symlinks.

### Autenticación básica
Para el requerimiento de autenticación básica, es necesario:
- Agregar dos líneas al archivo de configuración de nginx:

      auth_basic                    "nginx auth";
      auth_basic_user_file		/etc/nginx/.htpasswd;

- Crear un archivo `htpasswd` con usuarios y contraseñas creados via `htpasswd -c /etc/nginx/.htpass 'user'` e ingresar la contraseña. El flag `-c` se usa solo para crear el primer usuario (porque crea el archivo).

Luego, es posible acceder a los endpoints utilizando el flag `-u USER:SECRET` de curl, que a fines prácticos se traduce en `-u usuario:contraseña`

### Modificando el /etc/hosts
Si uno hiciese 'sensors.com/api/processing/submit', la request se redirigiría a la ip cuyo domain name coincida con sensors.com según el servidor de dns. Como en este caso, queremos acceder a la dirección de localhost, hay que modificar el archivo `/etc/hosts` agregando las lineas:

    127.0.0.1	dashboard.com
    127.0.0.1	sensors.com

De esta forma, como siempre se comprueba este archivo antes de consultar al servidor de dns, ahora sí nos redirigirá al servidor de nginx. Ahora, si la dirección es la misma, ¿Por qué dashboard.com/api/processing/submit no lleva a sensors.com/api/processing/submit? Por los virtual hosts.

Gracias a la línea: `server_name 'servicio'.com` es que nginx puede discriminar qué request corresponde a qué servicio.

También se crea otro archivo de nginx con un default_server que redirecciona todo a 404 (en caso de no coincidir el nombre de dominio con el resto de virtual hosts, se matchea con este servidor).

## Testing

Tanto el folder de users_service como el de processing_service cuentan con una carpeta tests, donde hay algunos tests de bash.

### Processing_service
Funciona igual que el del tp5 y tiene el mismo nombre: `test_n_post_dos_get.sh`, solo que se realiza todo con curl, y se pasan las credenciales de autenticación básica:
- Se hace un GET al /summary para saber el valor del contador al iniciar el test
- Se realizan N1 POST al /submit
- Se realiza un GET y se guarda en CONTADOR_1
- Se realizan N2 POST
- Se realiza un GET y se guarda en CONTADOR_2
- Si (CONTADOR_1 - CONTADOR_BASE == N1_POST) y (CONTADOR_2 - CONTADOR_BASE) == (N1_POST+N2_POST) se imprime `TEST_PASSED`, si no `TEST_FAILED`

### Users_service
En este caso hay varios tests:
- `test_login_fail.sh` : Hace un curl al endpoint de login con credenciales incorrectas. Se espera que no devuelva token.
- `test_login_pass.sh` : Hace un curl al endpoint de login con credenciales correctas. Se espera que devuelva token.
- `test_get_users.sh`  : Hace un login correcto y con el token hace un curl a listall. Se espera que la cantidad de usuarios retornados sea igual a la cantidad de usuarios encontrados en el archivo /etc/passwd
- `test_create_user.sh`: Hace un login correcto y con el token hace un curl a createuser. Se espera encontrar el usuario en /etc/passwd.

