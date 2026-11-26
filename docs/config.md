# Konfiguration

Das Verhalten des Borg-Servers wird mittels einer [Konfigurationsdatei](https://github.com/Landesarchiv-Thueringen/borg/blob/main/config/server_config.yml) eingestellt. Die Datei bestimmt, wie die Werkzeuge angesprochen werden, unter welchen Bedingungen Validatoren ausgeführt werden und wie einzelne Werkzeugergebnisse gewichtet werden.

Die Konfigurationsdatei wird beim Start des Servers gelesen.

## Voreinstellungen

Borg wird mit einer bereits funktionalen Konfiguration ausgeliefert. Diese stellt sich vereinfacht wie folgt dar:

### Bedingungen für die Ausführung

| Werkzeug                  | Bedingung                                                                                                   |
| ------------------------- | ----------------------------------------------------------------------------------------------------------- |
| DROID                     | aus Performancegründen standardmäßig deaktiviert                                                            |
| Siegfried                 | wird immer ausgeführt                                                                                       |
| Tika                      | wird immer ausgeführt                                                                                       |
| Magika                    | wird immer ausgeführt                                                                                       |
| MediaInfo                 | MIME-Type beginnt mit audio oder video                                                                      |
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
| veraPDF (PDF/UA-Profile)  | PUID entspricht PDF/UA-1                                                                                    |
| ODF Validator             | MIME-Type beginnt mit application/vnd.oasis.opendocument.                                                   |
| OOXML Validator           | MIME-Type beginnt mit application/vnd.openxmlformats-officedocument.                                        |

### Gewichtung der Werkzeugergebnisse

| Werkzeug        | Gewichtung   | Bedingte Gewichtung | Bedingung                                       |
| --------------- | ------------ | ------------------- | ----------------------------------------------- |
| DROID           | 75%          |                     |                                                 |
| Siegfried       | 75%          |                     |                                                 |
| Tika            | 75%          |                     |                                                 |
| Magika          | Magika-Score | 0%                  | MIME-Type entspricht _application/octet-stream_ |
| MediaInfo       | 100%         |                     |                                                 |
| JHOVE           | 0%           | 100%                | Datei ist valide                                |
| veraPDF         | 0%           | 100%                | Datei ist valide                                |
| ODF Validator   | 0%           | 100%                | Datei ist valide                                |
| OOXML Validator | 0%           | 100%                | Datei ist valide                                |

### Bedingte Nachbewertung von Ergebnismengen

| Bedingung                                                                                       | Bewertung     |
| ----------------------------------------------------------------------------------------------- | ------------- |
| Formatversion entspricht _PDF/A-1a_ und es gibt ein gültiges Validationsergebnis für das Format | 100%          |
| Formatversion entspricht _PDF/A-1b_ und es gibt ein gültiges Validationsergebnis für das Format | 100%          |
| Formatversion entspricht _PDF/A-2a_ und es gibt ein gültiges Validationsergebnis für das Format | 100%          |
| Formatversion entspricht _PDF/A-2b_ und es gibt ein gültiges Validationsergebnis für das Format | 100%          |
| Formatversion entspricht _PDF/A-3a_ und es gibt ein gültiges Validationsergebnis für das Format | 100%          |
| Formatversion entspricht _PDF/A-3b_ und es gibt ein gültiges Validationsergebnis für das Format | 100%          |
| Formatversion entspricht _PDF/A-3u_ und es gibt ein gültiges Validationsergebnis für das Format | 100%          |
| Formatversion entspricht _PDF/UA_ und es gibt ein gültiges Validationsergebnis für das Format   | 100%          |
| Formatversion entspricht _PDF/UA-1_ und es gibt ein gültiges Validationsergebnis für das Format | 100%          |
