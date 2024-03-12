# BorgFormat

BorgFormat (kurz Borg) ist ein Programm für die Formaterkennung und -validierung. Es integriert mehrere externe Programme

## Roadmap

Die Weiterentwicklung von Borg wird sich hauptsächlich um die Integration neuer Werkzeuge und die Extraktion von weitereren Metadaten aus den Werkzeugergebnissen. Folgende Werkzeuge sind für die Integration in den nächsten Veröffentlichungen vorgesehen:

- [Google Magika](https://github.com/google/magika) für die Formaterkennung, besonders für textbasierte Formate
- [jpylyzer](https://github.com/openpreserve/jpylyzer) für die Validierung von JP2-Dateien (JPEG 2000 Part 1)

## Motivation

Die Formaterkennung und -validierung erfordern eine Vielzahl unterschiedlicher Programme. Aufgrund der Komplexität des Problems kann jedoch kein einzelnes Programm dieses vollständig lösen. In der Regel sind Anwendungen darauf spezialisiert, entweder Dateien mit unbekannten Formaten zu identifizieren oder eine Auswahl an Dateiformaten zu validieren.

Um eine möglichst umfassende Abdeckung bei der Identifizierung und Validierung von Dateiformaten zu erreichen, ist es daher notwendig, mehrere Programme miteinander zu kombinieren. Es gibt bereits einige Anwendungen, die verschiedene Programme für die Formaterkennung und -validierung einbinden. Diese eingebundenen Programme werden in der Regel direkt integriert oder lokal ausgeführt. Für Borg wurde jedoch ein anderer Ansatz gewählt. Die Programme werden nicht direkt integriert, sondern werden in eigenen Containern ausgeführt und über eine Web-API angesprochen. Das hat den Vorteil, dass die Abhängigkeit

## Technische Umsetzung

Jedes integrierte Programm wird mittels Docker in einem eigenen Container gestartet. Die Werkzeug-Container teilen ein gemeinsamen Speicher (Docker Volume) um die übermittelte Datei zu teilen.

## Installation

Für den Betrieb von Borg wird [Docker](https://docs.docker.com/) inklusive [Docker Compose](https://docs.docker.com/compose/). Wir empfehlen für den Betrieb die Installation auf einem Linux-Server. Es ist für Testzwecke aber durchaus denkbar Borg lokal

Um alle für den Betrieb von Borg benötigten Container zu starten, genügt der folgende Befehl:

```sh
docker compose up --build -d
```

## Integrierte Programme

| Name            | Funktion        | Resourcen                                            |
| --------------- | --------------- | ---------------------------------------------------- |
| Droid           | Formaterkennung | https://digital-preservation.github.io/droid/        |
| Tika            | Formaterkennung | https://tika.apache.org/                             |
| JHOVE           | Validierung     | https://jhove.openpreservation.org                   |
| verapdf         | Validierung     | https://verapdf.org/                                 |
| ODF Validator   | Validierung     | https://odftoolkit.org/conformance/ODFValidator.html |
| OOXML-Validator | Validierung     | https://github.com/mikeebowen/OOXML-Validator        |

## Konfiguration

Das Verhalten des Borg-Servers wird mittels eine [Konfigurationsdatei](server/config/server_config.yml) eingestellt. Die Datei bestimmt, wie die Werkzeuge angesprochen werden, unter welchen Bedingungen Validatoren ausgeführt werden und wie einzelne extrahierte Eigenschaften gewichtet werden.

### Voreinstellungen

Borg wird mit einer funktionalen Konfiguration ausgeliefert.

#### Droid

**Bedingung für die Ausführung**

- wird immer ausgeführt

##### Extrahierte Eigenschaften

| Name | Standard Zuversichtswert |
| ---- | ------------------------ |
| PUID | 90%                      |

#### Tika

##### Bedingung für die Ausführung

- wird immer ausgeführt

##### Extrahierte Eigenschaften

| Name               | Standard Zuversichtswert |
| ------------------ | ------------------------ |
| MIME-Type          | 90%                      |
| Dateiformatversion | 90%                      |
| Textkodierung      | 90%                      |

#### JHOVE

##### Bedingung für die Ausführung

| Modulname   | Bedingung                                                                       |
| ----------- | ------------------------------------------------------------------------------- |
| PDF-Modul   | PUID entspricht PDF Version 1.0 bis 1.7                                         |
| HTML-Module | PUID entspricht HTML Version 3.2, 4.0 oder 4.01 (HTML 5 wird nicht unterstützt) |
| TIFF-Module | PUID entspricht TIFF (fmt/153)                                                  |
| JPEG-Module | PUID entspricht TIFF (fmt/153)                                                  |

##### Extrahierte Eigenschaften

Die extrahierten Eigenschaften und Zuversichtswerte sind für die meisten JHOVE-Module identisch. Falls die Werte abweichen, sind diese in einer gesonderten Übersicht aufgeführt.

| Name               | Standard Zuversichtswert |
| ------------------ | ------------------------ |
| Validität          | 100%                     |
| Wohlgeformtheit    | 100%                     |
| Dateiformatversion | 80%                      |

###### JHOVE (PDF-Modul)

| Name               | Standard Zuversichtswert | Bedingter Zuversichtswert                |
| ------------------ | ------------------------ | ---------------------------------------- |
| Validität          | 100%                     | 0% falls die Formatversion PDF/A enthält |
| Wohlgeformtheit    | 100%                     | 0% falls die Formatversion PDF/A enthält |
| Dateiformatversion | 80%                      | 0% falls die Formatversion PDF/A enthält |

#### veraPDF

##### Extrahierte Eigenschaften

| Name      | Standard Zuversichtswert |
| --------- | ------------------------ |
| Validität | 100%                     |

#### ODF Validator

##### Extrahierte Eigenschaften

| Name      | Standard Zuversichtswert |
| --------- | ------------------------ |
| Validität | 100%                     |

#### OOXML-Validator

##### Extrahierte Eigenschaften

| Name      | Standard Zuversichtswert |
| --------- | ------------------------ |
| Validität | 100%                     |
