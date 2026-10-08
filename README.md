# ВышКат — консольное приложение учёта парка кикшеринга

Домашнее задание №1 по КПО (ВШЭ, ФКН): доменная модель, SOLID, DI-контейнер, юнит-тесты.
Слоган: **«Доедем до дедлайна»**.

Приложение ведёт учёт парка «ВышКат»: приём техники после техосмотра, инвентарные номера для
техники и вещей, суммарное суточное потребление энергии, список техники, допустимой для новичков.

---

## 1. Запуск

Требуется Go ≥ 1.27 (см. `go.mod`), зависимостей внешних нет.

```bash
go run ./cmd/vysh-kat        # запуск приложения
go test ./...                 # тесты
go test ./... -coverprofile=cover.out && go tool cover -func=cover.out | tail -1   # покрытие
```

Меню приложения:

```
------- ВышКат -------
0. Выход
1. Добавить транспорт
2. Добавить вещь
3. Суммарное суточное потребление энергии
4. Транспорт для новичков
5. Вывести все единицы и вещи
```

### Демонстрация сценариев

Полный сценарий (весь ввод подаётся одной лентой, удобно для демонстрации):

```bash
printf '1\nсамокат\nКрасный самокат\nINV-1\n1.5\n8\n2\nшлем\nСтроительный шлем\nINV-9\n3\n4\n5\n0\n' \
  | go run ./cmd/vysh-kat
```

```
--- Добавление транспорта ---
Тип транспорта (самокат / велосипед / электровелосипед): Имя средства: Инвентарный номер: Суточный расход энергии (кВт·ч): Простота для новичка (1-10): Транспорт принят в парк

--- Добавление вещи ---
Тип вещи (шлем / док-станция): Имя вещи: Инвентарный номер: Вещь добавлена на баланс

--- Суммарное суточное потребление энергии ---
1.50 кВт·ч

--- Транспорт для новичков ---
- Красный самокат (INV: INV-1), простота: 8

--- Все единицы и вещи на балансе ---
[Транспорт] Красный самокат (INV: INV-1)
[Вещь]      Строительный шлем (INV: INV-9)
```

Обрабатываемые ошибочные сценарии (выводятся как `Ошибка: ...`, приложение продолжает работу):

| Сценарий | Что происходит |
|---|---|
| Неизвестный тип транспорта/вещи | `ErrInvalidTransportName` / `ErrInvalidThingName`, ничего не добавляется |
| Пустое имя или пустой инвентарный номер | `ErrEmptyName` / `ErrEmptyInventoryNumber` |
| Дубликат инвентарного номера (в т.ч. между техникой и вещами) | `ErrDuplicateInventory` |
| Расход энергии не число / простота не число | `ErrEnergyNotANumber` / `ErrSimplicityNotANumber` |
| Техосмотр отклонил технику | `ErrInspectionRejected` с причиной, техника не попадает в парк |
| Неизвестный пункт меню | `ErrInvalidCommand` |

---

## 2. Структура проекта

```
cmd/vysh-kat/main.go          — точка входа: app.NewApp(app.NewContainer()).Run()
internal/domain/              — доменная модель (без I/O и без зависимостей)
internal/inspection/          — техосмотр: ServiceCenter + RealServiceCenter
internal/novice_policy/       — правило «для новичков»: NovicePolicy + BasicNovicePolicy
internal/repository/          — хранение парка: Repository + MemoryRepository
internal/service/             — сценарии использования (валидация, приёмка, отчёты)
internal/app/                 — консольный UI (App) + DI-контейнер (Container)
internal/mocks/               — mock/fake-объекты для тестов
```

Зависимости по слоям (высокий → низкий):

```
cmd → app(App) → app(Container) → service.Service → repository.Repository
                                              ↘ inspection.ServiceCenter
                      repository.MemoryRepository → novice_policy.NovicePolicy
```

---

## 3. Логика решений (доменная модель)

**Абстракции.** В `internal/domain` три маленьких интерфейса:

* `InventoryItem` (`internal/domain/item.go`) — всё, что ставится на баланс: `Name()`, `InventoryNumber()`.
  Инвентарный номер вынесен в отдельный интерфейс, а не продублирован в каждом типе.
* `EnergyConsumer` (`internal/domain/energy.go`) — потребитель энергии: `DailyEnergyKWh() float64`.
  Отделён от инвентаризации: шлем инвентаризируется, но энергию не ест.
* `Transport` (`internal/domain/transport.go`) = `InventoryItem` + `EnergyConsumer` + `GetSimplicity() int` —
  композиция интерфейсов вместо наследования.
* `Thing` (`internal/domain/thing.go`) = `InventoryItem` — вещи на балансе.

**Виды техники.** `Scooter`, `Bicycle`, `EBike` (`internal/domain/scooter.go`, `bicycle.go`, `ebike.go`) —
независимые структуры, реализующие `Transport`; общей базовой структуры нет, избыточное
наследование не эмулируется. Вещи: `Helmet`, `DockingStation` (`helmet.go`, `docking_station.go`).

**Фабрика.** `domain.NewTransport(kind, ...)` / `domain.NewThing(kind, ...)` (`internal/domain/factory.go`)
— единственное место, где строковый тип из ввода превращается в конкретный тип (с нормализацией
регистра и пробелов). Весь остальной код работает с интерфейсами и о строках ничего не знает.

**Техосмотр.** `inspection.ServiceCenter` (`internal/inspection/servicecenter.go`) возвращает
`ResultOfInspection{Accepted, Reason}` — причина отказа важна: заказчик хочет видеть, *почему*
технику не приняли. `RealServiceCenter` (`real_servicecenter.go`) проверяет диапазон простоты 1–10
и неотрицательный расход энергии. Приёмка вызывается **до** записи в репозиторий
(`service.Service.AddTransport`, `internal/service/service.go:36`).

**Правило новичков.** Вынесено в отдельную политику `novice_policy.NovicePolicy`, базовая реализация —
`BasicNovicePolicy` с порогом `simplicity >= 6`. Политика передаётся в репозиторий и меняется
на лету через `Repository.ChangeNovicePolicy` — правило «можно усложнить» из ТЗ закрыто без правок кода.

**Уникальность инвентарных номеров.** Проверяется сервисом по всему парку сразу
(`MemoryRepository.HasInventoryNumber` ищет и среди техники, и среди вещей) — номер не может
«утонуть» между коллекциями.

**Энергия.** `MemoryRepository.TotalDailyEnergyKWh` суммирует технику, а для вещей делает
type-assertion на `domain.EnergyConsumer`: зарядный шкаф/док-станция как только начнут
потреблять энергию, попадут в отчёт без изменения кода отчёта (так проверяется тестом
`TestTotalDailyEnergyKWh` с локальным `energyBox`).

**Почему консоль и память.** ТЗ требует консольное приложение; `MemoryRepository` — самая простая
реализация `Repository`. БД, HTTP, JSON-конфиги не вводились сознательно: это был бы
overengineering, за который снимают баллы.

---

## 4. Принципы SOLID — обоснование

### S — Single Responsibility
У каждого пакета одна причина к изменению:
* `domain` — модель предметной области, ни I/O, ни логики приёма, ни хранения;
* `service` — сценарии: валидация, дедупликация номера, вызов техосмотра, отчёты
  (`internal/service/service.go`);
* `repository` — только хранение и выборки (`internal/repository/memory_repository.go`);
* `app` — только ввод/вывод (`internal/app/app.go`): разбор строк и `fmt.Sprintf`,
  бизнес-правил там нет, все решения принимает `service`.

### O — Open/Closed
Открытость для расширения — на уровне **потребителей**: `service`, `repository`, отчёты и
консоль зависят от интерфейсов `Transport`/`Thing`/`InventoryItem`/`EnergyConsumer`, поэтому
новый вид техники не меняет ни один из них (см. вопрос 3 ниже). Точка, где всё же нужна правка, —
фабрика `domain/factory.go`; она сознательно одна и закрыта для расширения извне (см. вопрос 2).

### L — Liskov Substitution
* Любой `Transport` (самокат/велосипед/электровелосипед) подставляется везде, где ожидают
  `Transport`: `Service.AddTransport`, `NoviceSuitableTransports`, вывод инвентарных номеров —
  поведение нигде не сужается и не паникует.
* `Transport` является `InventoryItem`, поэтому техника и вещи спокойно лежат в одном
  `[]InventoryItem` (`Repository.GetAllItems`) и обрабатываются type-switch'ем в `App.printAllItems`.
* Обратное утверждение не требуется: `Thing` не обязан уметь в `GetSimplicity`, и ни один код
  этого не требует.

### I — Interface Segregation
Интерфейсы маленькие и по роли: `EnergyConsumer` (1 метод), `InventoryItem` (2), `NovicePolicy` (1),
`ServiceCenter` (1). `Helmet` реализует только `InventoryItem` и не заставляется реализовывать
методы энергии или простоты — «толстых» интерфейсов, которые пришлось бы реализовывать
пустышками, в модели нет.

### D — Dependency Inversion (применён наиболее явно)
* `service.Service` зависит от абстракций `repository.Repository` и `inspection.ServiceCenter`
  и получает их через конструктор (`internal/service/service.go:19`):
  `NewService(rep repository.Repository, servCenter inspection.ServiceCenter)`.
* `repository.MemoryRepository` зависит от абстракции `novice_policy.NovicePolicy`.
* Высокий слой (`app`) собирается вручную в DI-контейнере и не знает о конкретиках ниже `Service`.

**DI-контейнер** — `internal/app/container.go:10`, `app.Container` + `app.NewContainer()`:
создаёт политику → репозиторий → сервисный центр → сервис и отдаёт готовый граф `App`.
Контейнер написан вручную (≈20 строк), без reflection-реестра и сторонних библиотек —
это ровно столько DI, сколько нужно проекту.

---

## 5. Тесты

```bash
go test ./... -cover
```

Ориентир ТЗ — покрытие ≥ 60%, фактически **62.7%** по проекту
(`go tool cover -func`: суммарно), при этом логика покрыта заметно плотнее:

| Пакет | Покрытие | Что проверяется |
|---|---|---|
| `internal/service` | 100% | приёмка/отклонение техосмотра, валидация до осмотра, дубликаты, делегирование отчётов |
| `internal/inspection` | 100% | таблица вердиктов `RealServiceCenter` (границы 1/10, отрицательная энергия) |
| `internal/novice_policy` | 100% | граница правила новичков: 5 → false, 6 → true |
| `internal/repository` | 94.3% | подмена политики, подсчёт энергии, дубликаты, `GetAllItems` |
| `internal/domain` | 61.8% | фабрика (регистр/пробелы/неизвестный тип) |
| `internal/app` | 60.3% | сценарии ввода/вывода через подменённый `stdin`/`stdout` |

Mock-объекты (`internal/mocks/mocks.go`):
* `MockServiceCenter` — фиксирует вердикт, считает вызовы `Inspect` и запоминает, что осматривали;
* `MockRepository` — записывает добавленное, отдаёт заданные отчёты;
* `MockNovicePolicy` — произвольное правило + счётчик вызовов.

Поведенческие тесты подменяют зависимости через конструкторы DI, а UI-тесты — через подмену
`reader`/`writer` в `App` (`internal/app/app_test.go`, `newTestApp`).

---

## 6. Обязательные вопросы

**• Какой принцип SOLID применён наиболее явно и где?**
Dependency Inversion. `service.Service` не знает ни о `MemoryRepository`, ни о `RealServiceCenter` —
он получает интерфейсы `repository.Repository` и `inspection.ServiceCenter` в конструкторе
(`internal/service/service.go:19`, метод `NewService`). Благодаря этому в тестах вместо реального
техосмотра подставляется `mocks.MockServiceCenter` (`internal/service/service_test.go`,
`newTestService`), а репозиторий — `mocks.MockRepository`, не меняя ни строки production-кода.
Второй явный случай — `repository.MemoryRepository`, зависящий от `novice_policy.NovicePolicy`
(конструктор + `ChangeNovicePolicy`).

**• Какой принцип сознательно ограничен / «почти нарушен» и почему это приемлемо?**
Open/Closed — в двух местах. (1) Фабрика `domain.NewTransport`/`NewThing`
(`internal/domain/factory.go:12`) — это `switch` по строковому типу, который правится при
добавлении нового вида. (2) Интерфейс `Repository` описывает конкретные коллекции
(`AddTransport`, `AddThing`, `GetThings`, `GetTransports`), поэтому новая сущность потребует
новых методов в интерфейсе. Приемлемо, потому что: это ровно **одна** точка правки на сущность,
а не разброс по всему коду; глобальный реестр/reflect-фабрика/generics были бы «мёртвым»
SOLID'ом и overengineering'ом, за который по критериям снимают баллы; все *потребители* модели
(`service`, `repository`, отчёты, UI) от этого не меняются вовсе. То есть нарушена «открытость
создания», но закрытость «потребления» сохранена — с точки зрения стоимости изменений это то,
что нужно.

**• Что сломается (или не сломается) при расширении?**
* *Новый вид транспорта* (например, монборд): добавляется структура в `internal/domain`,
  `case` в `factory.NewTransport` и строка-подсказка в меню (`internal/app/app.go`,
  `addTransport`) + текст в `ErrInvalidTransportName`. **Не ломается ничего**: репозиторий,
  сервис, техосмотр, все отчёты, DI-контейнер и существующие тесты опираются на интерфейс
  `Transport` и не меняются.
* *Сотрудник службы эксплуатации*: по ТЗ на баланс не ставится, поэтому в модель он не входит.
  Если понадобится — это новый тип в `domain` + новые методы в `Repository`/`MemoryRepository`
  и `MockRepository` + пункт меню в `App`. Существующий код **скомпилируется и продолжит
  работать**, но правки затронут все перечисленные слои — это прямое следствие сознательно
  ограниченного OCP (см. вопрос 2), а не авария.
* *Склад запчастей*: если это инвентаризируемая единица — уже сейчас попадает в `Thing`/`InventoryItem`,
  и отчёт «все единицы и инвентарные номера» заработает без правок; если она ещё и потребляет
  энергию (зарядный шкаф) — `TotalDailyEnergyKWh` подхватит её type-assertion'ом
  (`internal/repository/memory_repository.go`), отчёт тоже не меняется. Останется добавить
  `case` в `factory.NewThing` и, при желании, пункт меню.

**• Какой тест фиксирует решение техосмотра и правило «для новичков»; как подменяется зависимость?**
* Техосмотр: `TestAddTransport_InspectionRejected` (`internal/service/service_test.go`) —
  в конструктор `service.NewService(rep, sc)` вместо `RealServiceCenter` передаётся
  `mocks.MockServiceCenter{Result: {Accepted: false, Reason: "..."}}`; тест фиксирует, что
  возвращается `domain.ErrInspectionRejected` с причиной, **репозиторий остаётся пустым** и
  `Inspect` вызван ровно один раз. Обратный случай — `TestAddTransport_Success` (принятая техника
  попадает в парк), плюс `TestAddTransport_EmptyName`/`_DuplicateInventoryNumber` доказывают,
  что валидация идёт *до* осмотра (`InspectCalls == 0`). Реальные правила осмотра покрыты
  таблицей `TestRealServiceCenter_Inspect` (`internal/inspection/real_servicecenter_test.go`).
* Правило новичков: `TestBasicNovicePolicy_IsNoviceFriendly`
  (`internal/novice_policy/basic_novice_policy_test.go`) — таблица с границей 5/6/10;
  подмена зависимости — в `TestNoviceSuitableTransports_UsesProvidedPolicy`
  (`internal/repository/memory_repository_test.go`): в `repository.NewMemoryRepository(policy)`
  передаётся `mocks.MockNovicePolicy` со своим правилом (`simplicity >= 3`), и тест проверяет,
  что отбор делает именно она (в т.ч. `policy.Calls == 2`). Там же
  `TestChangeNovicePolicy` — смена политики на лету и игнорирование `nil`.

**• Что сделал ИИ, а что изменили вы?**
ИИ использовался как инструмент при работе над проектом: обсуждение и отбор идей
архитектуры (границы слоёв, вынос техосмотра и правила новичков в отдельные интерфейсы,
отказ от глобального реестра), подсказки при рефакторинге и формулировка текста этого README.
Код написан и размещён вручную; студент отвечает за понимание каждой строки, читал и правил
предложенные варианты, прогоняет тесты и может доработать проект под новое требование
на защите (п. политики ИИ в задании: ответственность за понимание — на студенте).

**Что переработано самостоятельно** (существенно изменено/дополнено после ИИ-подсказок):
уточнение состава `domain` под формулировки заказчика (инвентарные номера, вещи на балансе,
энергия), реализация `MemoryRepository` и проверка уникальности номеров между коллекциями,
вариативная `NovicePolicy` + `ChangeNovicePolicy`, вся консольная обработка ошибок и тексты
на русском, а также наборы тестов по границам (1/10, 5/6) и дубликатам номеров.

---

## 7. Соответствие ТЗ

| Требование | Статус |
|---|---|
| Приём техники с техосмотром (принять/отклонить) | `service.Service.AddTransport` + `inspection.ServiceCenter` |
| Суммарное суточное потребление энергии парка | `MemoryRepository.TotalDailyEnergyKWh`, п. 3 меню |
| Техника для новичков (простота ≥ 6) | `novice_policy.BasicNovicePolicy`, п. 4 меню |
| Учёт вещей и вывод инвентарных номеров вместе с техникой | `Repository.GetAllItems`, п. 5 меню |
| SOLID обоснован в README | §4 |
| DI-контейнер + подмена техосмотра в тестах | `app.NewContainer`, `mocks.MockServiceCenter` |
| Юнит-тесты, покрытие ≥ 60% | 62.7% (`go test ./... -cover`) |
| Демонстрация сценариев и понятный ввод/вывод | §1 |
