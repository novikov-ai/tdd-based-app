# Минимизируем шаги изменений

### [Задача по умножению двух строк](https://github.com/novikov-ai/leetcode/blob/master/medium/multiply_strings/multiply_strings.go)

[Тесты](https://github.com/novikov-ai/leetcode/blob/master/medium/multiply_strings/multiply_strings_test.go)

История коммитов:

1. init
2. make test pass
3. dummy implementation
4. add case with zero
5. add zero case reverse
6. add new case and implement
7. add new case
8. clean up

Реализацию задачи начал с объявления сигнатуры функции. Затем написал "глупый" тест, который бы проходил. 

Далее шаг за шагом расширял тестовые кейсы и добавлял недостающую функциональность в реализацию основной функции. 

Когда работа над функцией была завершена, удалил лишнее из кода. 

### Задача по добавлению нового ендпоинта в существующий сервис

[Тесты](../../internal/server/server_test.go)


История коммитов:

1. init stub for new feature
2. fix tests
3. register handler with players league
4. implement players of league fetching
5. add test for new endpoint
6. refactor

Текущее приложение работает с in-memory БД, поэтому реализацию задачи начал с добавления соответствующего функционала на уровне интерфейса (с обозначением пред и постусловий), но без непосредственной реализации. 

Очевидно, существующие тесты потребовали также правок. 

После этого, добавил регистрацию нового эндпоинта (параллельно писал тесты, но не коммитил). 

Добавил реализацию на уровне БД. 

Когда тесты стали проходить успешно, сделал их коммит. 

Подумал как можно еще упростить код и сделать его лучше поддерживающим - поправил. 

### Выводы

1. Удобно делать много понятных коммитов, при этом важно, чтобы они составляли единый инкремент. 
2. Каждая из задач потребовала 6-8 коммитов. 
3. В ежедневной работе обычно количество коммитов варьируется от 3-15 (в зависимости от сложности задачи).

Желательно любое атомарное изменение заворачивать в соответствующий коммит с понятным комментарием, но также и важно придерживаться разумному балансу, чтобы не создавать излишнего шума из большого количества изменений. 