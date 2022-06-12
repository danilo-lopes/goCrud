# goCrud

Simples API em Golang que Faz CRUD em Banco MySQL

## Uso

```
docker run -d -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=password -v /$HOME/.docker-volumes/dev-book:/var/lib/mysql mysql/mysql-server:8.0
```

Crie usuario de conexão externo

```
CREATE USER 'admin'@'%' IDENTIFIED BY 'password';
```

```
GRANT ALL ON *.* TO 'admin'@'%';
```

```
go run main.go
```

## API

### Endpoints

#### Get

`/usuarios`

`/usuario/{id}`

#### Delete

`/usuario/{id}`

#### Update

`/usuario/{id}`

#### Create

`/usuario/`