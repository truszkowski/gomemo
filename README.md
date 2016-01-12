# gomemo - prosty serwer http key-value.

Projekt służy wprowadzeniu w język golang, od strony pisania mikroserwisow z API http. Zaimplementowane zostały dwa proste endppointy http, służące dodawaniu danych pod wskazany klucz oraz pobieraniu danych z ponadanego klucza.

W celu ćwiczenia należy zaimplementować kolejne funkcjonalności z sekcji TODO, opisać błedy jakie program posiada (jeśli).

## Opis działania

Serwer przechowuje w pamięci dane klucz-wartość. 
Klucz obiektu może być dowolną wartością alfa-numeryczną (regexp: `[a-zA-Z0-9]+`) o maksymalnie 100 znakach.
Dane obiektu mogą być dowolnymi danymi o maksymalnym rozmiarze 1 MB.

## Instalacja

Ze strony: https://golang.org/dl/

```
$ tar xf go1.5.2.linux-amd64.tar.gz
$ mv go ~/goroot
$ mkdir ~/gopath
$ export GOPATH=${HOME}/gopath
$ export GOROOT=${HOME}/goroot
$ export PATH=${GOPATH}/bin:${GOROOT}/bin:${PATH}
$ go get github.com/truszkowski/gomemo
$ cd ~/gopath/src/github.com/truszkowski/gomemo
```

## Budowanie, instalowanie, testy

Pobranie zależności:
```
$ go get -v -d -t ./...
# lub
$ make deps
```

Zbudowanie i instalacja (instaluje do katalogu ${GOPATH}):

```
$ go build -v ./...
$ go install -v ./...
# lub 
$ make build install
# uruchomienie:
$ gomemo
...
```

Testy:

```
$ go test -v ./...
# lub
$ make test
```

# Serwer HTTP API

Poniżej opis aktualnie zaimplementowanych endpointow HTTP oraz ich krótki opis z przykładami.

## `PUT /v1/objects/<object_id>`

Dodawanie danych pod klucz <object_id>.

```
$ curl -si 127.0.0.1:1234/v1/objects/abc -XPUT -d 'przykladowe dane'
HTTP/1.1 201 Created
```

```
$ curl -si 127.0.0.1:1234/v1/objects/niepoprawny_klucz -XPUT -d 'przykladowe dane'
HTTP/1.1 400 Bad Request
```

```
$ curl -si 127.0.0.1:1234/v1/objects/abc -XPUT -d 'ponad 1 MB danych...'
HTTP/1.1 413 Request Entity Too Large
```

## `GET /v1/objects/<object_id>`

Pobieranie danych spod klucza `<object_id>`.

```
$ curl -si 127.0.0.1:1234/v1/objects/abc
HTTP/1.1 200 Ok
```

```
$ curl -si 127.0.0.1:1234/v1/objects/abc--
HTTP/1.1 400 Bad Request
```

```
$ curl -si 127.0.0.1:1234/v1/objects/niema
HTTP/1.1 404 Not Found
```

# TODO

Zadania do zrobienia.

## 1. Obsługa błędów

Zdefiniowanie odpowiednich limitów i ograniczeń kiedy serwer powinien zwracać odpowiednie błędy. 

Np, co gdy będziemy chcieli dodać plik o rozmiarze 100 GB?

## 2. `GET /v1/objects`

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

## 3. `DELETE /v1/objects/<object_id>`

Dodanie endpointu HTTP usuwającego wskazany obiekt.

```
$ curl -s 127.0.0.1:1234/v1/objects/a1 -XPUT -d 'v1'
$ curl -s 127.0.0.1:1234/v1/objects/a1
v1
$ curl -s 127.0.0.1:1234/v1/objects/a1 -XDELETE
$ curl -s 127.0.0.1:1234/v1/objects/a1 -i
HTTP/1.1 404 Not Found
```

## 4. Dopisać test zapisu danych o niepoprawnym kluczu.

W `main_test.go` dopisać test sprawdzający poprawne zachowanie serwera w przypadku podania niepoprawnych danych dla endpointu `PUT ...`.

## 5. Automatyczne czyszczenie obiektów

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
HTTP/1.1 404 Not Found
```

## 6. Dodanie obsługi sygnałów SIGINT i SIGTERM

Jeśli program otrzyma któryś z sygnałów, powinien zamknąć socket nasłuchujący oraz zakończyć pracę.

## 7. Dodanie obsługi typu danych (MIME)

Dodanie obslugi typu danych według MIME, czyli tego co jest podawane w nagłówku `"Content-Type: text/html"`. 
Klient dodając dane, np typu text/html, przy ich odczycie powinien dostać w odpowiedzi HTTP również taki nagłówek.

Np:
```
$ curl -s 127.0.0.1:1234/v1/objects/plik.html -XPUT -H'Content-Type: text/html' --data-binary @plik.html
HTTP/1.1 201 Created
$ curl -si 127.0.0.1:1234/v1/objects/plik.html
HTTP/1.1 200 Ok
Content-Type: text/html
...
```

## 8. Zapisywanie obiektów na dysk

Serwer powinien zapisywać dane nie w pamięci, lecz na dysku aby po zrestartowaniu serwera dane były nadal dostępne.

## 9. Dostęp tylko z wybranych IP, klas IP

Umożliwienie zapisów tylko z wybranych adresów/klas IP IP. Dla pozostałych adresów na próby `PUT` powinniśmy zwracać kod `403 Forbidden`.
