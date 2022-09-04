# url-shortener-swoyo

## Task

Требуется создать микросервис сокращения url на хосту localhost и порту 8080<br>
Эндпоинты:<br>
* POST

Request: (body): 
```json
{
  "url": "http://cjdr17afeihmk.biz/kdni9/z9d112423421"
}
```
Response: `http://localhost:8080/88d2d0f8fe07c98da23165c7a8a7acae`
* GET

Request (url query): `http://localhost:8080/88d2d0f8fe07c98da23165c7a8a7acae`
Response (body): 
```json
{
  "url": "http://cjdr17afeihmk.biz/kdni9/z9d112423421"
}
```

Микросервис должен уметь хранить информацию в памяти и в postgres в зависимости от флага  
запуска -d