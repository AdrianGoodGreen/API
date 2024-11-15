# Usar una imagen de Go como base
FROM golang:1.20-alpine AS builder

# Configurar el directorio de trabajo
WORKDIR /app

# Copiar los archivos go.mod y go.sum
COPY go.mod ./
COPY go.sum ./

# Descargar las dependencias
RUN go mod download

# Copiar todos los archivos fuente de Go
COPY . .

# Ejecuta las pruebas
RUN go test -v ./...

# Construir la aplicación
RUN go build -o main .

# Imagen ligera para producción
FROM alpine:latest

# Copiar el binario desde el contenedor de construcción
WORKDIR /root/
COPY --from=builder /app/main .

# Exponer el puerto donde la aplicación escuchará
EXPOSE 8083

# Comando para ejecutar la aplicación
CMD ["./main"]
