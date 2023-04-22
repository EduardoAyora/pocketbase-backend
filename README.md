### Despliegue

```
docker build -t pb .
docker run -d --mount type=bind,source="$(pwd)"/pb_data,target=/pb/pb_data -p 8080:8080 pb
```

### Quitar contenedor una vez se detiene

```
docker run -d --rm -p 8080:8080 pb
docker run -it --rm -p 8080:8080 pb

-it is short for --interactive
-d is short for --detach
```
