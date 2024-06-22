# Домашнее задание по курсу HLF 2024

## Условие
- (+1 балл) написать чейнкод типа saac с методами set / get, разграничением доступа по MSP ID и пользователям: 
   - только создатель записи может менять свою запись
   - только участники одного MSP могут получать данные
- ( +2 балла) реализовать клиент на базе fabric-chaincode-api с методами:
   - GET /{username}/assets/get
   - POST /{username}/assets/set
- (+3 балла) применил современные паттерны по layer архитектуре, IoC / DI (uber-go.fx)

## Реализация
- Клиент разворачивается на `::8082`
- Модель Asset находится в директории `saac-v2/saac-v2-core`
- Chaincode представляет из себя усовершенствованный saac - saac-v2
- Доступны три метода с соответствующими ручками на клиенте
  - CreateAsset(asset Asset) -- `/assets/create` (asset передается в body)
  - ReadAsset(id int) -- `/assets/get` (id передается через заголовок `id`)
  - UpdateAsset(asset Asset) -- `/assets/update` (asset передается в body)
- При обращении на любую ручку клиента также необходимо передавать заголовок `user`, в котором указаывается имя клиента из каталога organizations в формате `ClientName@OrgMspName`(пр. `User1@Org1MSP`)
- Папка organizations для удобства является полной копией аналогичной директории из fabric-samples/test-network после поднятия сети
- Имя канала регулируется через env-переменную `CHANNEL_NAME`, а chaincode - `CHAINCODE_NAME`

## Особенности национального деплоя
- Для деплоя chaincode **ОБЯЗАТЕЛЬНО** использовать политику `OR('Org2MSP.peer','Org1MSP.peer')`
- При поднятии сети и создании канала необходимо использовать флаг `-ca`, иначе будут проблемы с сертификатами

## Пример команд для развертывания
- `./network.sh createChannel -ca -s couchdb -c testchannel` - поднятие сети с помощью скрипта `./network.sh` из `fabric-samples`
- `./network.sh deployCC -c testchannel -ccn saacv2 -ccp ../saac-v2/saac-v2-chaincode -ccl go -ccv 1 -cccg ../saac-v2/saac-v2-chaincode/collections-config.json -ccep "OR('Org2MSP.peer','Org1MSP.peer')"` - deploy chaincode
