# Group membership

В этой задаче вам надо реализовать распределенный сервис для отслеживания состава группы (group membership).

Пусть имеется распределенная система, состоящая из _N_ узлов (node). С каждым узлом свяжем процесс _Node_. У каждого узла есть свой уникальный идентификатор. Будем говорить, что все работающие в некоторый момент узлы системы образуют группу, то есть являются ее участниками. В ходе выполнения системы узлы могут включаться и выключаться, тем самым вступая в группу или покидая её. Таким образом, состав участников группы может изменяться со временем. Сервис, который вам надо реализовать, должен позволять узлам узнавать текущий состав группы, а также выполнять вход и выход из группы.

Сервис должен поддерживать следующие локальные операции на каждом узле системы:
- `JOIN(seed)` - выполнить подключение локального узла к группе. В результате этой операции узел должен стать участником группы. Операции передается адрес одного из доступных в настоящий момент участников группы (seed). Если передается адрес самого узла, то он должен создать новую пустую группу и добавить в неё себя.
- `LEAVE()` - выполнить выход локального узла из группы. В результате этой операции узел должен перестать быть участником группы.
- `GET_MEMBERS()` - выдать список текущих участников группы. В результате этой операции возвращается список идентификаторов узлов.

Как обычно, данные операции должны быть реализованы с помощью локальных сообщений. Запросы `JOIN` и `LEAVE` не предполагают какого-либо ответа. Ответ на запрос `GET_MEMBERS` надо вернуть в виде сообщения `MEMBERS`.

Помимо отключения в штатном режиме путем отправки `LEAVE`, узлы могут покидать группу из-за отказов. В нашей системе возможны следующие виды отказов:
- Узел внезапно прекращает свою работу (падает). Позднее узел может быть перезапущен и заново подключен к группе.
- Отказы сети, приводящие к потере сообщений между некоторыми парами узлов, возможно только в одном направлении. В результате нормально работающие узлы могут отказаться полностью или частично отрезаны от других узлов. Из-за разделения сети (network partition) система может распасться на несколько изолированных групп. Отказы сети обычно носят временный характер, после чего связь восстанавливается.

Ваша реализация должна обнаруживать описанные выше отказы и обрабатывать их. Для этого вам потребуется реализовать _детектор отказов_. Ваша реализация детектора отказов должна гарантировать полноту (если узел стал недоступным из-за отказа, то в конечном счете все живые узлы признают его отказавшим). Также постарайтесь найти хороший баланс между скоростью обнаружения отказов, точностью детектора (долей false positives) и сопутствующими накладными расходами (объем служебного трафика, нагрузка на узел).

Отказавшие или недоступные ни с одного из участников группы узлы должны автоматически исключаться из состава группы. Таким образом, в состав группы должны входить только живые узлы, способные общаться по сети (в обе стороны) с хотя бы одним другим участником группы. Вам не требуется обеспечивать строгую согласованность списков участников между узлами. Иными словами, допустимо, если в некоторый момент времени списки участников группы на разных узлах отличаются (например, отказавший узел кое-где еще выдается как участник группы). Главное, чтобы через некоторое время после отказа или входа/выхода участника состав группы стал одинаковым и корректным на всех живых узлах.

Наконец, ваша реализация должна масштабироваться на большое число узлов. А именно, с ростом размера системы накладные расходы (служебный трафик, нагрузка на узел) должны расти не быстрее, чем линейно от числа узлов _N_.

В `solution` находится заготовка для решения. Обратите внимание, что используется [событийная модель программирования](../../dslib/readme.md#обработка-событий-callbacks) из _dslib_ на основе обратных вызовов. Эта модель более удобна для данной задачи, так как поддерживает _таймеры_ для реализации периодических выполняемых действий (не забывайте продлевать таймеры, так как они срабатывают только один раз).

## Тестирование

Необходимые зависимости должны быть уже установлены (см. предыдущие задания). Добавьте корень репозитория в переменную окружения PYTHONPATH и запустите тесты:

```console
$ PYTHONPATH=$ROOT_DIRECTORY_PATH python3 test.py solution
```

Если Вам нужна бОльшая debug информация, добавьте флаг `-d` к тестированию.

Поднимаются 10 процессов узлов (число можно изменить с помощью флага `-n`), которые коммуницируют между собой через `test_server`. Все аналогично предыдущим задачам.

Мы уже написали за Вас тесты, которы должны проходить при правильном решении. Вы можете добавлять новые тесты, по умолчанию они не будут оцениваться. Однако если вы придумаете и реализуете тест, который увеличит покрытие наших тестов, то за него можно получить бонус 2 балла. За деталями о том, как мы тестируем, обращайтесь к тестам, осознать, что там происходит -- Ваша задача.

## Оценивание

Баллы за успешное прохождение тестов (максимум 9):
- BasicTestCase: 1
- RandomSeedTestCase: 1
- NodeJoinTestCase: 0.5
- NodeLeaveTestCase: 0.5
- NodeCrashTestCase: 1
- NodeCrashRecoverTestCase: 1
- NodeOfflineTestCase: 0.5
- NodeOfflineRecoverTestCase: 0.5
- NetworkPartitionTestCase: 0.5
- NetworkPartitionRecoverTestCase: 0.5
- NodeCannotReceiveTestCase: 0.5
- FlakyNetworkTestCase: 0.5
- FlakyNetworkStartTestCase: 0.5
- FlakyNetworkCrashTestCase: 0.5

Еще 1 балл добавляется, если ваша реализация масштабируется на большое число узлов.

## Сдача

Сдача и проверка решений будет вестись через сервис Gradescope.
См. инструкцию в первой задаче.

Далее зайдите в Assignments, там вы увидите задачу Membership, в которую можно
сдавать только zip архивы. Поместите тесты и решение в архив следующей командой:

```console
$ zip -r solution.zip test.py solution/
```

Добавьте  в папку `solution` файл `readme.md` с кратким описанием вашего решения и дополнительными комментариями, который также попадет в архив. Также, если вы добавляете свои тесты, то напишите к ним docstring - что проверяется в тесте и каким образом. **В случае отсутствия описаний решения и тестов оценка может быть снижена на 1 балл**.

Далее сдайте Ваш `solution.zip`. Дождитесь пока отработают Ваши тесты на Вашем решении (учтите, что это может занять несколько минут из-за нагрузки серверов). Если все тесты прошли успешно, то решение проставляется 2 балла и решение принимается на ручную проверку. Если какие-то тесты не прошли, то выставляется 0 баллов, и решение не принимается на ручную проверку.

Если Вы не смогли что-то реализовать, можете закомментировать некоторые тесты в [test.py](./test.py), но обязательно отразите почему Вы это сделали в readme-файле. Это сильно упростит нам проверку.

Дедлайн задачи -- __2 недели__, в Gradescope корректно проставлена дата окончания.