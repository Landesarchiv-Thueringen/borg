# Konfiguration

Das Verhalten des Borg-Servers wird mittels eine [Konfigurationsdatei](config/server_config.yml) eingestellt. Die Datei bestimmt, wie die Werkzeuge angesprochen werden, unter welchen Bedingungen Validatoren ausgeführt werden und wie einzelne Werkzeugergebnisse gewichtet werden.

Die Konfigurationsdatei wird beim Start des Servers gelesen.

## Voreinstellungen

Borg wird mit einer bereits funktionalen Konfiguration ausgeliefert. Diese stellt sich vereinfacht wie folgt dar:

### Bedingungen für die Ausführung

| Werkzeug                  | Bedingung                                                                                                   |
| ------------------------- | ----------------------------------------------------------------------------------------------------------- |
| DROID                     | wird immer ausgeführt                                                                                       |
| Tika                      | wird immer ausgeführt                                                                                       |
| Magika                    | wird immer ausgeführt                                                                                       |
| JHOVE (PDF-Modul)         | PUID entspricht PDF Version 1.0 bis 1.7                                                                     |
| JHOVE (HTML-Modul)        | PUID entspricht HTML Version 3.2, 4.0 oder 4.01 (HTML 5 wird nicht unterstützt) oder MIME-Type enthält html |
| JHOVE (TIFF-Modul)        | PUID entspricht TIFF oder MIME-Type enthält tiff                                                            |
| JHOVE (JPEG-Modul)        | PUID entspricht JPEG oder MIME-Type enthält jpeg                                                            |
| JHOVE (JPEG2000-Modul)    | PUID entspricht JP2 (JPEG 2000 part 1) oder MIME-Type enthält jp2                                           |
| veraPDF (PDF/A-1a-Profil) | PUID entspricht PDF/A-1a oder Formatversion entspricht PDF/A-1a                                             |
| veraPDF (PDF/A-1b-Profil) | PUID entspricht PDF/A-1b oder Formatversion entspricht PDF/A-1b                                             |
| veraPDF (PDF/A-2a-Profil) | PUID entspricht PDF/A-2a oder Formatversion entspricht PDF/A-2a                                             |
| veraPDF (PDF/A-2b-Profil) | PUID entspricht PDF/A-2b oder Formatversion entspricht PDF/A-2b                                             |
| veraPDF (PDF/A-2u-Profil) | PUID entspricht PDF/A-2u oder Formatversion entspricht PDF/A-2u                                             |
| veraPDF (PDF/UA-Profile)  | MIME-Type enthält pdf, nach aktuellen Stand keine PUID verfügbar                                            |
| ODF Validator             | MIME-Type beginnt mit application/vnd.oasis.opendocument.                                                   |
| OOXML Validator           | MIME-Type beginnt mit application/vnd.openxmlformats-officedocument.                                        |

### Gewichtung der extrahierten Eigenschaften

| Werkzeug                 | PUID | MIME-Type    | Formatversion         | Validierung            |
| ------------------------ | ---- | ------------ | --------------------- | ---------------------- |
| DROID                    | 90%  | 90%          | 90%                   |                        |
| Tika                     |      | 90%          | 90%                   |                        |
| Magika                   |      | Magika-Score |                       |                        |
| JHOVE                    |      |              | 80%                   | 100%                   |
| JHOVE (PDF-Modul)        |      |              | 80% bzw. 0% bei PDF/A | 100% bzw. 0% bei PDF/A |
| veraPDF                  |      |              |                       | 100%                   |
| veraPDF (PDF/UA-Profile) |      |              |                       | 30%                    |
| ODF Validator            |      |              |                       | 100%                   |
| OOXML Validator          |      |              |                       | 100%                   |
