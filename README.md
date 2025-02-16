# **Merch shop system**

Микросервис для сотрудников, позволяющий покупать товары в магазине, обмениваться монетками между пользователями и получать информацию о пользователях.

### Используемые технологии:
* База данных - `PostgreSQL`
* Миграции - `golang-migrate/migrate`
* Управление БД - `pgxpool`
* Менеджер тразнакций - `avito-tech/go-transaction-manager`
* Router - `chi router`
* Развертывание - `Docker`
* Тестирование - `testify`, `mock`, `testcontainers`
* Архитектура - `Clean Architecture`

### Запуск
Запустить сервис можно с помощью команды `make run-build`

Запустить тестовую среду для e2e тестирования - `make e2e-stand`

Запустить тесты - `make test`

Запустить тесты с покрытием - `make test-cover`

Тестовая среда для интеграционных тестов создается с помощью `go-testcontainers`