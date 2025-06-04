# BorgFormat

BorgFormat (kurz Borg) ist ein Programm für die Formaterkennung und -validierung. Die Anwendung integriert verschiedene Werkzeuge um eine möglichst umfassende Abdeckung bei der Identifizierung und Validierung von Dateiformaten zu erreichen.

## Lizenz

Dieses Projekt wird unter der [GNU General Public License Version 3 (GPLv3)](https://www.gnu.org/licenses/gpl-3.0.de.html) veröffentlicht. Weitere Informationen finden Sie in der Datei [LICENSE](LICENSE).

Diese Lizenz gilt nicht für eingebundene Werkzeuge, die jeweils unter eigenen Lizenzen stehen (siehe unten).

## Integrierte Werkzeuge

| Name            | Version | Funktion        | Resourcen                                                        | Lizenz                                                                                                        |
| --------------- | ------- | --------------- | ---------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------- |
| Droid           | 6.8.1   | Formaterkennung | [Homepage ](https://digital-preservation.github.io/droid/)       | [BSD License](https://github.com/digital-preservation/droid/blob/master/license.md)                           |
| Tika            | 2.9.2   | Formaterkennung | [Homepage](https://tika.apache.org/)                             | [Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)                                    |
| Magika          | 0.6.0   | Formaterkennung | [GitHub](https://github.com/google/magika)                       | [Apache License, Version 2.0](https://github.com/google/magika?tab=Apache-2.0-1-ov-file)                      |
| JHOVE           | 1.30.1  | Validierung     | [Homepage](https://jhove.openpreservation.org)                   | [GNU Lesser General Public License](https://www.gnu.org/licenses/lgpl-3.0.html)                               |
| veraPDF         | 1.26.2  | Validierung     | [Homepage](https://verapdf.org/)                                 | [GNU General Public License v3.0](https://github.com/veraPDF/veraPDF-validation/blob/integration/LICENSE.GPL) |
| ODF Validator   | 0.12.0  | Validierung     | [Homepage](https://odftoolkit.org/conformance/ODFValidator.html) | [Apache License, Version 2.0](https://github.com/tdf/odftoolkit/blob/master/validator/LICENSE.txt)            |
| OOXML Validator | 2.1.5   | Validierung     | [GitHub](https://github.com/mikeebowen/OOXML-Validator)          | [MIT License](https://github.com/mikeebowen/OOXML-Validator/blob/main/LICENSE)                                |

## Roadmap

Die Weiterentwicklung von Borg wird sich hauptsächlich um die Integration neuer Werkzeuge und die Extraktion von weitereren Metadaten aus den Werkzeugergebnissen drehen. Folgende Weiterentwicklungen sind für die nächsten Veröffentlichungen vorgesehen:

- ~~Integration des Werkzeugs [Google Magika](https://github.com/google/magika) für die Formaterkennung, besonders für textbasierte Formate~~ ✓
- ~~Integration des Werkzeugs [MediaInfo](https://mediaarea.net/de/MediaInfo) für die Extraktion von AV-Metadaten~~ ✓
- ~~Verfeinerung der Bedingungen für das Ansprechen der Validatoren durch Kombination der verschiedenen extrahierten Eigenschaften (Erweiterung der Konfiguration)~~ ✓
- Integration des Werkzeugs [jpylyzer](https://github.com/openpreserve/jpylyzer) für die Validierung von JP2-Dateien (JPEG 2000 Part 1)
- PDF-Export von Gesamt- und Teilergebnissen
