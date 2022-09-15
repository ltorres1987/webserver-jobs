# Project Name: webserver-jobs

### Descripción

Crea un Web server que genera jobs y son procesados por workerpools.

### Construcción 🛠️
* **Language:** Golang


## Instalación

Pasos:

1. Clone el proyecto.
2. Ejecute el archivo ```main.go``` para subir el servidor web.

## Consumo de la Api

Pasos:

1. Consumir el endpoint ```http://localhost:8081/fib``` con el body
   ```
   curl --location --request POST 'http://localhost:8081/fib' \
   --header 'Content-Type: application/x-www-form-urlencoded' \
   --data-urlencode 'name=Fib25' \
   --data-urlencode 'value=12' \
   --data-urlencode 'delay=3s'
   ```