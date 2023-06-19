# Informe Lab 6: Sistemas Embebidos / Webservices - Desarrollo

## Diseño de los endpoints

Ambos endpoints siguen el patron de diseño de Controller - Service - Repository, al igual que el ejemplo presentado en clases[[1]](https://github.com/SOiI-UNC/go-example).

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


## TODO temas a tocar en el informe:
### Servicio (systemd)              .
### Conexion base de datos sqlite   .
### Autenticacion con JWT           .
### Nginx 
### Nginx - basic auth  
### Sudoers                         .
### Sshd                            .
### Initializers                    .
### Variables de entorno            .
### Testing
### Script de instalacion - uso
