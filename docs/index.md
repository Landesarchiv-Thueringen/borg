# BorgFormat

BorgFormat (kurz Borg) ist ein Programm für die Formaterkennung und -validierung. Die Anwendung integriert verschiedene Werkzeuge um eine möglichst umfassende Abdeckung bei der Identifizierung und Validierung von Dateiformaten zu erreichen.

Borg lässt sich auf zwei Arten nutzen:

Zum einen stehen die Funktionen zur Formatverifikation über eine REST-API zur Verfügung, die sich in beliebige Systeme integrieren lässt – beispielsweise in ein digitales Archiv.

Zum anderen gibt es eine Webanwendung zur manuellen Analyse, die direkt im Browser am Arbeitsplatz verwendet werden kann.

## Lizenz

Dieses Projekt wird unter der [GNU General Public License Version 3 (GPLv3)](https://www.gnu.org/licenses/gpl-3.0.de.html) veröffentlicht.

Diese Lizenz gilt nicht für eingebundene Werkzeuge, die jeweils unter eigenen Lizenzen stehen (siehe unten).

## Integrierte Werkzeuge

| Name            | Version       | Funktion            | Ressourcen                                                       | Lizenz                                                                                                        |
| --------------- | ------------- | ------------------- | ---------------------------------------------------------------- | ------------------------------------------------------------------------------------------------------------- |
| Droid           | 6.8.1         | Formaterkennung     | [Homepage ](https://digital-preservation.github.io/droid/)       | [Apache License, Version 2.0](https://github.com/richardlehane/siegfried?tab=Apache-2.0-1-ov-file)            |
| Siegfried       | 1.11.2        | Formaterkennung     | [Homepage ](https://www.itforarchivists.com/siegfried)           | [BSD License](https://github.com/digital-preservation/droid/blob/master/license.md)                           |
| Tika            | 2.9.2         | Formaterkennung     | [Homepage](https://tika.apache.org/)                             | [Apache License, Version 2.0](https://www.apache.org/licenses/LICENSE-2.0)                                    |
| Magika          | standard_v3_3 | Formaterkennung     | [GitHub](https://github.com/google/magika)                       | [Apache License, Version 2.0](https://github.com/google/magika?tab=Apache-2.0-1-ov-file)                      |
| MediaInfo       | 24.11         | Metadatenextraktion | [Homepage](https://mediaarea.net/en/MediaInfo)                   | [BSD-style license](https://mediaarea.net/en/MediaInfo/License)                                               |
| JHOVE           | 1.30.1        | Formatvalidierung   | [Homepage](https://jhove.openpreservation.org)                   | [GNU Lesser General Public License](https://www.gnu.org/licenses/lgpl-3.0.html)                               |
| veraPDF         | 1.26.2        | Formatvalidierung   | [Homepage](https://verapdf.org/)                                 | [GNU General Public License v3.0](https://github.com/veraPDF/veraPDF-validation/blob/integration/LICENSE.GPL) |
| ODF Validator   | 0.12.0        | Formatvalidierung   | [Homepage](https://odftoolkit.org/conformance/ODFValidator.html) | [Apache License, Version 2.0](https://github.com/tdf/odftoolkit/blob/master/validator/LICENSE.txt)            |
| OOXML Validator | 2.1.5         | Formatvalidierung   | [GitHub](https://github.com/mikeebowen/OOXML-Validator)          | [MIT License](https://github.com/mikeebowen/OOXML-Validator/blob/main/LICENSE)                                |
