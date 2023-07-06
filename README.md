# EasyBalance

Микросервис для работы с балансом пользователей

Это простой микросервис на Go для управления балансом пользователей. Он предоставляет HTTP API для выполнения следующих операций:

Создание нового пользователя с указанием их идентификатора и начального баланса.
Получение текущего баланса пользователя по идентификатору.
Зачисление средств на баланс пользователя.
Списание средств с баланса пользователя.
Перевод средств от одного пользователя к другому.

Использование API

Микросервис предоставляет следующие эндпоинты (Все примеры запросов для PowerShell):

POST '/users': Создание нового пользователя. Пример запроса: 
Invoke-WebRequest -Method POST -Uri http://localhost:8080/users -Headers @{"Content-Type" = "application/json"} -Body '{"ID":"user123","Balance":100.0}'

GET '/balance?user_id={user_id}': Получение баланса пользователя по идентификатору. Пример запроса: 
Invoke-WebRequest -Method GET -Uri "http://localhost:8080/balance?user_id=user123"

POST '/balance/update': Обновление баланса пользователя (зачисление). Пример запроса: 
Invoke-WebRequest -Method POST -Uri http://localhost:8080/balance/update -Headers @{"Content-Type" = "application/json"} -Body '{"user_id":"user123","amount":50.0,"is_debit":false}'

POST '/balance/update': Обновление баланса пользователя (списание средств). Пример запроса: 
Invoke-WebRequest -Method POST -Uri http://localhost:8080/balance/update -Headers @{"Content-Type" = "application/json"} -Body '{"user_id":"user123","amount":30.0,"is_debit":true}'

Зачисление и списание зависит от значения флага 'is_debit';
