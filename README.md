# gomemo - prosty serwer http key-value.

Projekt służy wprowadzeniu w język golang, od strony pisania mikroserwisow z API http. Zaimplementowane zostały dwa proste endppointy http, służące dodawaniu danych pod wskazany klucz oraz pobieraniu danych z ponadanego klucza.

W celu ćwiczenia należy zaimplementować kolejne funkcjonalności z sekcji TODO, opisać błedy jakie program posiada (jeśli).

## Opis działania

Serwer przechowuje w pamięci dane klucz-wartość. 
Klucz obiektu może być dowolną wartością alfa-numeryczną (regexp: `[a-zA-Z0-9]+`) o maksymalnie 100 znakach.
Dane obiektu mogą być dowolnymi danymi o maksymalnym rozmiarze 1 MB.

# Serwer HTTP API

Poniżej opis aktualnie zaimplementowanych endpointow HTTP oraz ich krótki opis z przykładami.

## `PUT /v1/objects/<object_id>`

Dodawanie danych pod klucz <object_id>.

```
$ curl -si 127.0.0.1:1234/v1/objects/abc -XPUT -d 'przykladowe dane'
200 Ok
```

```
$ curl -si 127.0.0.1:1234/v1/objects/niepoprawny_klucz -XPUT -d 'przykladowe dane'
400 Bad Request
```

```
$ curl -si 127.0.0.1:1234/v1/objects/abc -XPUT -d 'ponad 1 MB danych...'
413 Request Entity Too Large
```

## `GET /v1/objects/<object_id>`

Pobieranie danych spod klucza `<object_id>`.

```
$ curl -si 127.0.0.1:1234/v1/objects/abc
200 Ok
```

```
$ curl -si 127.0.0.1:1234/v1/objects/abc__
400 Bad Request
```

```
$ curl -si 127.0.0.1:1234/v1/objects/niema
404 Not Found
```

# TODO

Zadania do zrobienia.

## Obsługa błędów

Zdefiniowanie odpowiednich limitów i ograniczeń kiedy serwer powinien zwracać odpowiednie błędy. 

Np, co gdy będziemy chcieli dodać plik o rozmiarze 100 GB?

## `GET /v1/objects`

Dodanie endpointu HTTP listującego dostępne w pamięci obiekty

```
$ curl -s 127.0.0.1:1234/v1/objects/a1 -XPUT -d 'v1'
$ curl -s 127.0.0.1:1234/v1/objects/a2 -XPUT -d 'v2'
...
$ curl -s 127.0.0.1:1234/v1/objects/a9 -XPUT -d 'v9'
$ curl -s 127.0.0.1:1235/v1/objects
a1
a2
...
a9
```

## `DELETE /v1/objects/<object_id>`

Dodanie endpointu HTTP usuwającego wskazany obiekt.

```
$ curl -s 127.0.0.1:1234/v1/objects/a1 -XPUT -d 'v1'
$ curl -s 127.0.0.1:1234/v1/objects/a1
v1
$ curl -s 127.0.0.1:1234/v1/objects/a1 -XDELETE
$ curl -s 127.0.0.1:1234/v1/objects/a1 -i
404 Not Found
```

## Dopisać test zapisu danych o niepoprawnym kluczu.

W `main_test.go` dopisać test sprawdzający poprawne zachowanie serwera w przypadku podania niepoprawnych danych dla endpointu `PUT ...`.

## Automatyczne czyszczenie obiektów

Serwer powinien czyścić z pamięci obiekty dodane lub pobrane po minucie.

```
$ curl -s 127.0.0.1:1234/v1/objects/a1 -XPUT -d 'v1'
(po 40 sekundach)
$ curl -s 127.0.0.1:1234/v1/objects/a1
v1
(po 40 sekundach)
$ curl -s 127.0.0.1:1234/v1/objects/a1
v1
(po 61 sekundach)
$ curl -s 127.0.0.1:1234/v1/objcets/a1 -i
404 Not Found
```

## Dodanie obsługi sygnałów SIGINT i SIGTERM

Jeśli program otrzyma któryś z sygnałów, powinien zamknąć socket nasłuchujący oraz zakończyć pracę.

## Zapisywanie obiektów na dysk

Serwer powinien zapisywać dane nie w pamięci, lecz na dysku aby po zrestartowaniu serwera dane były nadal dostępne.

