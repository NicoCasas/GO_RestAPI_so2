# Informe Lab 6: Sistemas Embebidos / Webservices - Desarrollo

## Diseño de los endpoints

Ambos endpoints siguen el patron de diseño de Controller - Service - Repository, al igual que el ejemplo presentado en clases.

## Servicio de usuario

### Endpoint de autenticación - POST /api/users/login
Para realizar las peticiones a los endpoints del sistema operativo es necesario estar autorizado. Esta autorización se logra a través de un jwt. Para obtenerlo, es necesario primero hacer un POST a dashboard.com/api/users/login con las credenciales `username` y `password` de un usuario autorizado.

    curl -s --request POST \
             --url http://dashboard.com/api/users/login \
                -u admin:admin \
                --header 'accept: application/json' \
                --header 'content-type: application/json' \
                --data '{"username": "nicoAdmin", "password": "nicoPass"}'

Las credenciales de dicho usuario se encuentran en una base de datos sql. Para gestionarla, se usa sqlite3. Como bibliografía se usaron [[1]](https://github.com/SOiI-UNC/go-example) [[2]](https://zetcode.com/golang/sqlite3/).
En el jwt, se incluyen dos campos en el payload:
  - iss: Nombre de usuario autenticado
  - exp: Indica el tiempo de expiración de la validez del token. En este caso, el token dura una hora.

### Endpoint de creación de usuario - POST /api/users/createuser
Este endpoint tiene como finalidad crear un usuario en el sistema operativo para que luego pueda acceder al sistema via ssh.

#### Creación del usuario
Para esta parte se usó el package [os/exec](https://pkg.go.dev/os/exec) siguiendo los ejemplos encontrados en [3](https://zetcode.com/golang/exec-command/). Estas funciones tienen comportamientos similares tipo fork/exec de c en el sentido de que no invocan una shell. 
Los comandos realizados de esta forma en este endpoint son:

    sudo useradd -g 'group_id' -s '/bin/bash' username     -> Para crear el usuario de nombre username. En el grupo 'group_id' y  que inicie sesion en bash.  
    sudo passwd username                                   -> Para asignarle contraseña al usuario.  

En el caso del último comando, la contraseña se pasa 2 veces vía stdin. Es por esto que se crea un string con las contraseñas separadas por \n y se lo pasa via pipe. 

#### Acceso via ssh
Para este punto, se decidió crear un grupo `operativos_ssh_clients` cuyo group_id, (obtenido del `/etc/group`) es el pasado como argumento en el comando useradd del punto anterior.
Luego, se le da acceso a todos los usuarios del mismo a través de la sentencia

    AllowGroups operativos_ssh_clients


## TODO temas a tocar en el informe:
### Servicio (systemd)
### Conexion base de datos sqlite
### Autenticacion con JWT
### Nginx - basic auth
### Sudoers
### Sshd
### Initializers
### Variables de entorno