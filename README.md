# api-gen-doc
server witch generates docs (pdf, word) by template

Создание шаблона
## POST `http://localhost:8080/api/gendoc`
Отправляем

```json
{
  "recordID": 1,
  "urlTemplate": "https://sycret.ru/service/apigendoc/forma_025u.xml",
  "text": "1",
  "use": "Иван,Иванов,Иванович"
}
```

Получаем

```json
{
	"resultdata": "1, Иван,Иванов,Иванович",
	"resultdescription": "Ok"
}

```
Поиск документов
## POST `http://localhost:8080/api/find`
Отправляем

```json
{
  "recordID": 1,
  "urlTemplate": "https://sycret.ru/service/apigendoc/forma_025u.xml"
}
```

Получаем

```json
{
	"pdf": [
		"http://localhost:8080/api/download/forma_025u.xml/1/pdf/2022-06-09%2004:16:55.pdf"
	],
	"word": [
		"http://localhost:8080/api/download/forma_025u.xml/1/word/2022-06-09%2004:16:55.doc"
	]
}
```

Получаем word
## GET `http://localhost:8080/api/download/forma_025u.xml/1/word/2022-06-09%2004:16:55.doc`

Получаем pdf
## GET `http://localhost:8080/api/download/forma_025u.xml/1/pdf/2022-06-09%2004:16:55.pdf`


