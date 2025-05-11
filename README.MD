# Веб-сервис Арифметических Вычислений
## Описание
Этот проект реализует веб-сервис для вычисления арифметических выражений. Пользователи могут отправлять математические выражения через HTTP-запросы, а сервис возвращает результат вычисления. Сервис поддерживает базовые арифметические операции: сложение, вычитание, умножение, деление, а также использование скобок для определения приоритетов операций.
## Структура Проекта
```
.
│   .env                # Переменные окружения
│   .gitignore          # Игнорируемые файлы Git
│   go.mod              # Модули Go
│   go.sum              # Контрольные суммы зависимостей
├───cmd
│       main.go         # Точка входа в приложение
└───internal
    ├───agent
    │       agent.go         # Логика агента
    │       agent_test.go    # Тесты для агента
    ├───ast
    │       ast.go           # Построение AST
    ├───auth
    │       auth.go          # JWT и авторизация
    ├───database
    │       connect.go       # Подключение к базе данных
    ├───models
    │       expression.go    # Модель выражения
    │       user.go          # Модель пользователя
    ├───proto
    │   │   agent.proto      # gRPC-протокол
    │   └───gen
    │           agent.pb.go         # Сгенерированный код protobuf
    │           agent_grpc.pb.go    # Сгенерированный gRPC-код
    ├───server
    │       server.go         # HTTP-сервер на Echo
    │       server_test.go    # Интеграционные тесты сервера
    └───utils
            utils.go          # Вспомогательные функции
```
## Установка
1. **Клонируйте Репозиторий:**
```bash
   git clone https://github.com/Yorshik/Go-Final.git
```
2. **Перейдите в Каталог Проекта:**
```bash
   cd Go-Final
```
3. Установите компилятор C (GCC)
### Windows
1. Установите [Msys2](https://www.msys2.org)
2. Откройте MSYS2 MSYS
3. выполните:
```
pacman -Syu
```
Если в процессе выполнения закроется терминал, откройте еще раз и напишите:
```
pacman -Su
```
Далее выполните команду:
```
pacman -S mingw-w64-x86_64-gcc
```
Перезапустите терминал и проверьте что GCC установился:
``` 
gcc --version
```
4. **Добавьте `C:\msys64\mingw64\bin` в системную переменную `PATH`**
### Linux (Debian/Ubuntu)
```
sudo apt update
sudo apt install build-essential
```
### MacOS
``` 
xcode-select --install
```
5. **Установите Необходимые Зависимости:**
```bash
go mod download
```
6. Перед запуском проекта заполните .env
- DB_NAME - имя базы данных
- TIME_ADDITION_MS - время сложения
- TIME_SUBTRACTION_MS - время вычитания
- TIME_MULTIPLICATIONS_MS - время умножения
- TIME_DIVISIONS_MS - время деления
- COMPUTING_POWER - количество воркеров
- JWT_SECRET - секретный ключ для создания JWT-токенов
## Запуск Проекта
Для запуска веб-сервиса выполните следующую команду из корневого каталога проекта:
```bash
go run ./cmd/main.go/
```
Сервис будет доступен по адресу [localhost:8080/api/v1/calculate](http://localhost:9000/api/v1/calculate).

# API документация
## Регистрация:
### Метод - POST
### На входе - username и password
### На выходе - сообщение (успех/неудача)
Пример запроса:
#### Windows
```
curl --location --request POST "http://localhost:8080/api/v1/register" ^
--header "Content-Type: application/json" ^
--data "{\"username\": \"testuser\", \"password\": \"testpass\"}"
```
#### Linux/MacOS
``` 
curl --location --request POST 'http://localhost:8080/api/v1/register' \
--header 'Content-Type: application/json' \
--data '{
  "username": "testuser",
  "password": "testpass"
}'
```
Результат:
``` 
{
  "message": "Registration is successfull"
}
```
## Авторизация
### Метод - POST
### На входе - username и password
### На выходе - token (если авторизация прошла успешно)
Пример запроса:
#### Windows
```
curl --location --request POST "http://localhost:8080/api/v1/login" ^
--header "Content-Type: application/json" ^
--data "{\"username\": \"testuser\", \"password\": \"testpass\"}"
```
#### Linux/MacOS
``` 
curl --location --request POST 'http://localhost:8080/api/v1/login' \
--header 'Content-Type: application/json' \
--data '{
  "username": "testuser",
  "password": "testpass"
}
```
Результат:
``` 
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```
## Калькулятор
### Метод - POST
### На входе - expression (+JWT токен для авторизации)
### На выходе - result (если все прошло хорошо, иначе - ошибка)
Пример запроса:
#### Windows
``` 
curl --location --request POST "http://localhost:8080/api/v1/calculate" ^
--header "Content-Type: application/json" ^
--header "Authorization: YOUR_JWT_TOKEN_HERE" ^
--data "{\"expression\": \"2+2*2\"}"
```
#### Linux
``` 
curl --location --request POST "http://localhost:8080/api/v1/calculate" \
--header "Content-Type: application/json" \
--header "Authorization: Bearer YOUR_JWT_TOKEN_HERE" \
--data '{"expression": "2+2*2"}'
```
Результат:
``` 
{
    "message": "Expression calculated",
    "result": "6.000000"
}
```
### Если ввели что-то некорректное (пример - несуществующий оператор `**`), то результат будет таким:
``` 
{
    "message": "Failed to send expression"
}
```
## Просмотр всех своих запросов
### Метод - GET
### На входе - ничего (только JWT токен для авторизации)
### На выходе - список выражений и результатов своего аккаунта
Пример запроса:
#### Windows
``` 
curl --location --request GET "http://localhost:8080/api/v1/expressions" ^
--header "Authorization: YOUR_JWT_TOKEN_HERE"
```
#### Linux
``` 
curl --location --request GET "http://localhost:8080/api/v1/expressions" \
--header "Authorization: YOUR_JWT_TOKEN_HERE"
```
Результат:
``` 
[
    {
        "id": 1,
        "expression": "2*2/0",
        "result": "+Inf",
        "user_id": 1
    },
    {
        "id": 2,
        "expression": "2*2*2",
        "result": "8.000000",
        "user_id": 1
    }
]
```
## Просмотр одного своего запроса
### Метод - GET
### На входе - через URL указываете ID запроса (+JWT токен для авторизации)
### На выходе - информация об этом запросе (выражение и его результат)
Пример запроса:
#### Windows
```
curl --location --request GET "http://localhost:8080/api/v1/expression/1" ^
--header "Authorization: YOUR_JWT_TOKEN_HERE"
```
#### Linux
``` 
curl --location --request GET "http://localhost:8080/api/v1/expression/1" \
--header "Authorization: YOUR_JWT_TOKEN_HERE"
```
Результат:
```
{
    "id": 2,
    "expression": "2*2*2",
    "result": "8.000000",
    "user_id": 1
}
```
