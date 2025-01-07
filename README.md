# Сервис отслеживания посылок

## Задание
Реализовать сервис отслеживания посылок со следующими функциями:
- регистрация посылки,
- получение списка посылок клиента,
- изменение статуса посылки,
- изменение адреса доставки,
- удаление посылки.
Информация о посылках хранится в БД. Посылка может быть зарегистрирована, отправлена или доставлена. При регистрации посылки создаётся новая запись в БД. У только что зарегистрированной должен быть статус «зарегистрирована». Трек-номер посылки равен её идентификатору в таблице. Если посылка в статусе «зарегистрирована», можно изменить адрес доставки или удалить посылку.
